package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"slices"
	"time"
)

// SourceReconciler reconciles a Source object
// It contains the logic for handling Source resources, including creating, updating, and deleting them.
type SourceReconciler struct {
	Client                      client.Client   // Client for interacting with the Kubernetes API
	Scheme                      *runtime.Scheme // Scheme for the reconciler
	HTTPClient                  *http.Client    // HTTP client for making external requests
	NewsAggregatorSrcServiceURL string          // URL of the news aggregator source service
}

const (
	SrcFinalizer     = "source.finalizers.teamdev.com" // Finalizer string for Source resources
	ReasonFailedUpd  = "FailedUpdate"                  // Reason for a failed update
	ReasonFailedCr   = "FailedCreation"                // Reason for a failed creation
	ReasonSuccessCr  = "SuccessfulCreation"            // Reason for a successful creation
	ReasonSuccessUpd = "SuccessfulUpdate"              // Reason for a successful update
)

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/finalizers,verbs=update

// Reconcile is the main logic for reconciling a Source resource.
// It handles creating, updating, and deleting sources in the news aggregator service.
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
		if !slices.Contains(source.Finalizers, SrcFinalizer) {
			source.Finalizers = append(source.Finalizers, SrcFinalizer)
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if slices.Contains(source.Finalizers, SrcFinalizer) {
			if _, err := r.deleteSource(source.Status.ID); err != nil {
				return ctrl.Result{}, err
			}
			source.Finalizers = slices.Delete(source.Finalizers,
				slices.Index(source.Finalizers, SrcFinalizer), slices.Index(source.Finalizers, SrcFinalizer)+1)
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
		return r.updateSource(ctx, source.Status.ID, &source)
	}
}

// createSource creates a new source in the news aggregator service due to a new Source resource being created.
func (r *SourceReconciler) createSource(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Creating source:", source.Spec.Name)

	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	logrus.Println("Source data:", string(sourceData))

	req, err := http.NewRequest(http.MethodPost, r.NewsAggregatorSrcServiceURL, bytes.NewBuffer(sourceData))
	if err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, resp.Status)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, fmt.Errorf("failed to create source in news aggregator: %s", resp.Status)
	}

	var createdSource aggregatorv1.SourceSpec
	if err := json.NewDecoder(resp.Body).Decode(&createdSource); err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	source.Status.ID = createdSource.Id
	if err := r.Client.Update(ctx, source); err != nil { // todo check is okay?
		return ctrl.Result{}, err
	}

	err = r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionTrue, ReasonSuccessCr, "Source successfully created in news aggregator")
	if err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully created source in news aggregator")
	return ctrl.Result{}, nil
}

// updateSource updates an existing source in the news aggregator service due to a Source resource being updated.
func (r *SourceReconciler) updateSource(ctx context.Context, sourceID int, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Updating source:", source.Spec.Name)

	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	url := fmt.Sprintf("%s/%d", r.NewsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(sourceData))
	if err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		err := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, resp.Status)
		if err != nil {
			return ctrl.Result{}, err
		}
		logrus.Error(fmt.Errorf("failed to update source in news aggregator"), "Status", resp.Status)
		return ctrl.Result{}, fmt.Errorf("failed to update source in news aggregator: %s", resp.Status)
	}

	err = r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionTrue, ReasonSuccessUpd, "Source successfully updated in news aggregator")
	if err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully updated source in news aggregator")
	return ctrl.Result{}, nil
}

// deleteSource deletes a source from the news aggregator service.
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

// updateSourceStatus updates the SourceStatus of a source resource with the given condition.
func (r *SourceReconciler) updateSourceStatus(ctx context.Context, source *aggregatorv1.Source, conditionType aggregatorv1.SourceConditionType, status metav1.ConditionStatus, reason, message string) error {
	newCondition := aggregatorv1.SourceCondition{
		Type:           conditionType,
		Status:         status,
		LastUpdateTime: metav1.Time{Time: time.Now()},
		Reason:         reason,
		Message:        message,
	}

	source.Status.Conditions = append(source.Status.Conditions, newCondition)
	return r.Client.Update(ctx, source) // todo check is okay
}

// SetupWithManager sets up the controller with the Manager, uses predicates to filter events.
func (r *SourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Source{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return !e.DeleteStateUnknown
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
			},
		}).
		Complete(r)
}
