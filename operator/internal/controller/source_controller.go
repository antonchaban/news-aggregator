package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
)

// SourceReconciler reconciles a Source object
type SourceReconciler struct {
	Client                      client.Client
	Scheme                      *runtime.Scheme
	HTTPClient                  *http.Client
	NewsAggregatorSrcServiceURL string
}

const (
	srcFinalizer = "source.finalizers.teamdev.com"
)

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/finalizers,verbs=update

func (r *SourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var source aggregatorv1.Source
	err := r.Client.Get(ctx, req.NamespacedName, &source)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info("Source resource not found, possibly deleted. Removing source from news-aggregator.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if source.ObjectMeta.DeletionTimestamp.IsZero() {
		if !slices.Contains(source.Finalizers, srcFinalizer) {
			source.Finalizers = append(source.Finalizers, srcFinalizer)
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if slices.Contains(source.Finalizers, srcFinalizer) {
			if _, err := r.deleteSource(source.Status.ID); err != nil {
				return ctrl.Result{}, err
			}
			source.Finalizers = slices.Delete(source.Finalizers,
				slices.Index(source.Finalizers, srcFinalizer), slices.Index(source.Finalizers, srcFinalizer)+1)
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	logrus.Info("Reconciling Source", "ID", source.Status.ID, "Name", source.Spec.Name, "Link", source.Spec.Link)

	if source.Status.ID == 0 {
		return r.createSource(ctx, &source)
	} else {
		return r.updateSource(source.Status.ID, &source)
	}
}

func (r *SourceReconciler) createSource(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Creating source:", source.Name)

	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	req, err := http.NewRequest(http.MethodPost, r.NewsAggregatorSrcServiceURL, bytes.NewBuffer(sourceData))
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ctrl.Result{}, fmt.Errorf("failed to create source in news aggregator: %s", resp.Status)
	}

	var createdSource aggregatorv1.SourceSpec
	if err := json.NewDecoder(resp.Body).Decode(&createdSource); err != nil {
		return ctrl.Result{}, err
	}

	source.Status.ID = createdSource.Id
	if err := r.Client.Status().Update(ctx, source); err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully created source in news aggregator")
	return ctrl.Result{}, nil
}

func (r *SourceReconciler) updateSource(sourceID int, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Updating source:", source.Name)

	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	url := fmt.Sprintf("%s/%d", r.NewsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(sourceData))
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error(fmt.Errorf("failed to update source in news aggregator"), "Status", resp.Status)
		return ctrl.Result{}, fmt.Errorf("failed to update source in news aggregator: %s", resp.Status)
	}

	logrus.Info("Successfully updated source in news aggregator")
	return ctrl.Result{}, nil
}

func (r *SourceReconciler) deleteSource(sourceID int) (ctrl.Result, error) {
	logrus.Println("Deleting source with ID:", sourceID)

	url := fmt.Sprintf("%s/%d", r.NewsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error(fmt.Errorf("failed to delete source from news aggregator"), "Status", resp.Status)
		return ctrl.Result{}, fmt.Errorf("failed to delete source from news aggregator: %s", resp.Status)
	}

	logrus.Info("Successfully deleted source from news aggregator")
	return ctrl.Result{}, nil
}

func (r *SourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Source{}).
		Complete(r)
}
