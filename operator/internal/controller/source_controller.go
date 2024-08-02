/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// SourceReconciler reconciles a Source object
type SourceReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

const (
	newsAggregatorSrcServiceURL = "https://news-alligator-service.news-alligator.svc.cluster.local:8443/sources"
)

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Source object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
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
		if !containsString(source.Finalizers, "source.finalizers.teamdev.com") {
			source.Finalizers = append(source.Finalizers, "source.finalizers.teamdev.com")
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if containsString(source.Finalizers, "source.finalizers.teamdev.com") {
			if _, err := r.deleteSourceFromAggregator(ctx, source.Status.ID); err != nil {
				return ctrl.Result{}, err
			}
			source.Finalizers = removeString(source.Finalizers, "source.finalizers.teamdev.com")
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	logrus.Info("Reconciling Source", "ID", source.Status.ID, "Name", source.Spec.Name, "Link", source.Spec.Link)

	if source.Status.ID == 0 {
		return r.createSourceInAggregator(ctx, &source)
	} else {
		return r.updateSourceInAggregator(ctx, source.Status.ID, &source)
	}
}

func (r *SourceReconciler) createSourceInAggregator(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Creating source:", source.Name)

	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	req, err := http.NewRequest(http.MethodPost, newsAggregatorSrcServiceURL, bytes.NewBuffer(sourceData))
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := c.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error(fmt.Errorf("failed to create source in news aggregator"), "Status", resp.Status)
		return ctrl.Result{}, fmt.Errorf("failed to create source in news aggregator: %s", resp.Status)
	}

	var createdSource aggregatorv1.SourceSpec
	if err := json.NewDecoder(resp.Body).Decode(&createdSource); err != nil {
		return ctrl.Result{}, err
	}

	// Update the Source status with the created Source ID
	source.Status.ID = createdSource.Id
	if err := r.Client.Status().Update(ctx, source); err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully created source in news aggregator")
	return ctrl.Result{}, nil
}

func (r *SourceReconciler) updateSourceInAggregator(ctx context.Context, sourceID int, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Updating source:", source.Name)

	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	url := fmt.Sprintf("%s/%d", newsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(sourceData))
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := c.Do(req)
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

func (r *SourceReconciler) deleteSourceFromAggregator(ctx context.Context, sourceID int) (ctrl.Result, error) {
	logrus.Println("Deleting source with ID:", sourceID)

	url := fmt.Sprintf("%s/%d", newsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := c.Do(req)
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

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// SetupWithManager sets up the controller with the Manager.
func (r *SourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Source{}).
		Complete(r)
}
