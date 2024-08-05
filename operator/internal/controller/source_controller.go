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
	"golang.org/x/exp/slices"
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
	newsAggregatorSrcServiceURL        = "https://news-alligator-service.news-alligator.svc.cluster.local:8443/sources"
	sourceFinalizer                    = "source.finalizers.teamdev.com"
	contentTypeHeader                  = "Content-Type"
	contentTypeJSON                    = "application/json"
	httpClientTimeout                  = 10 * time.Second
	logSourceNotFound                  = "Source resource not found, possibly deleted. Removing source from news-aggregator."
	logSourceReconcile                 = "Reconciling Source"
	logSourceCreate                    = "Creating source:"
	logSourceUpdate                    = "Updating source:"
	logSourceDelete                    = "Deleting source with ID:"
	logSourceCreateSuccess             = "Successfully created source in news aggregator"
	logSourceUpdateSuccess             = "Successfully updated source in news aggregator"
	logSourceDeleteSuccess             = "Successfully deleted source from news aggregator"
	logErrorCreateSourceInAggregator   = "failed to create source in news aggregator"
	logErrorUpdateSourceInAggregator   = "failed to update source in news aggregator"
	logErrorDeleteSourceFromAggregator = "failed to delete source from news aggregator"
)

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var source aggregatorv1.Source
	err := r.Client.Get(ctx, req.NamespacedName, &source)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info(logSourceNotFound)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if source.ObjectMeta.DeletionTimestamp.IsZero() {
		if !slices.Contains(source.Finalizers, sourceFinalizer) {
			source.Finalizers = append(source.Finalizers, sourceFinalizer)
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if slices.Contains(source.Finalizers, sourceFinalizer) {
			if _, err := r.deleteSource(source.Status.ID); err != nil {
				return ctrl.Result{}, err
			}
			source.Finalizers = slices.Delete(source.Finalizers, slices.Index(source.Finalizers, sourceFinalizer),
				slices.Index(source.Finalizers, sourceFinalizer)+1)
			if err := r.Client.Update(ctx, &source); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	logrus.Info(logSourceReconcile, "ID", source.Status.ID, "Name", source.Spec.Name, "Link", source.Spec.Link)

	if source.Status.ID == 0 {
		return r.createSource(ctx, &source)
	} else {
		return r.updateSource(ctx, source.Status.ID, &source)
	}
}

func (r *SourceReconciler) createSource(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println(logSourceCreate, source.Name)

	return r.sendSrcReq(ctx, http.MethodPost, newsAggregatorSrcServiceURL,
		source, logSourceCreateSuccess, logErrorCreateSourceInAggregator)
}

func (r *SourceReconciler) updateSource(ctx context.Context, sourceID int, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println(logSourceUpdate, source.Name)

	url := fmt.Sprintf("%s/%d", newsAggregatorSrcServiceURL, sourceID)
	return r.sendSrcReq(ctx, http.MethodPut, url, source,
		logSourceUpdateSuccess, logErrorUpdateSourceInAggregator)
}

func (r *SourceReconciler) deleteSource(sourceID int) (ctrl.Result, error) {
	logrus.Println(logSourceDelete, sourceID)

	url := fmt.Sprintf("%s/%d", newsAggregatorSrcServiceURL, sourceID)
	return r.sendRequest(http.MethodDelete, url, nil, logSourceDeleteSuccess,
		logErrorDeleteSourceFromAggregator)
}

func (r *SourceReconciler) sendSrcReq(ctx context.Context, method, url string,
	source *aggregatorv1.Source, successMsg, errorMsg string) (ctrl.Result, error) {
	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	result, err := r.sendRequest(method, url, bytes.NewBuffer(sourceData), successMsg, errorMsg)
	if err == nil && method == http.MethodPost {
		var createdSource aggregatorv1.SourceSpec
		if err := json.NewDecoder(bytes.NewBuffer(sourceData)).Decode(&createdSource); err != nil {
			return ctrl.Result{}, err
		}
		source.Status.ID = createdSource.Id
		if err := r.Client.Status().Update(ctx, source); err != nil {
			return ctrl.Result{}, err
		}
	}
	return result, err
}

func (r *SourceReconciler) sendRequest(method, url string,
	body *bytes.Buffer, successMsg, errorMsg string) (ctrl.Result, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set(contentTypeHeader, contentTypeJSON)

	c := r.getHTTPClient()
	resp, err := c.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error(fmt.Errorf(errorMsg), "Status", resp.Status)
		return ctrl.Result{}, fmt.Errorf("%s: %s", errorMsg, resp.Status)
	}

	logrus.Info(successMsg)
	return ctrl.Result{}, nil
}

func (r *SourceReconciler) getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: httpClientTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Source{}).
		Complete(r)
}
