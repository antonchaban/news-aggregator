package controller

import (
	"bytes"
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller/predicates"
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
	"slices"
	"time"
)

// SourceReconciler reconciles a Source object.
// It contains the logic for handling Source resources, including creating, updating, and deleting them.
type SourceReconciler struct {
	Client                      client.Client   // Client for interacting with the Kubernetes API.
	Scheme                      *runtime.Scheme // Scheme for the reconciler.
	HTTPClient                  *http.Client    // HTTP client for making external requests.
	NewsAggregatorSrcServiceURL string          // URL of the news aggregator source service.
}

const (
	SrcFinalizer     = "source.finalizers.teamdev.com" // Finalizer string for Source resources.
	ReasonFailedUpd  = "FailedUpdate"                  // Reason for a failed update.
	ReasonFailedCr   = "FailedCreation"                // Reason for a failed creation.
	ReasonSuccessCr  = "SuccessfulCreation"            // Reason for a successful creation.
	ReasonSuccessUpd = "SuccessfulUpdate"              // Reason for a successful update.
)

// Reconcile is the main logic for reconciling a Source resource.
// It handles creating, updating, and deleting sources in the news aggregator service.
// - Retrieves the Source resource specified in the reconcile request.
// - If the resource is marked for deletion, it calls the deleteSource method and removes the finalizer.
// - If the resource does not exist in the news aggregator, it creates it.
// - If the resource exists, it updates the resource in the news aggregator.
func (r *SourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var source aggregatorv1.Source
	err := r.Client.Get(ctx, req.NamespacedName, &source)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info("Source resource not found, possibly deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle resource finalization logic
	if !source.ObjectMeta.DeletionTimestamp.IsZero() {
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

	if !slices.Contains(source.Finalizers, SrcFinalizer) {
		source.Finalizers = append(source.Finalizers, SrcFinalizer)
		if err := r.Client.Update(ctx, &source); err != nil {
			return ctrl.Result{}, err
		}
	}

	logrus.Info("Reconciling Source", "ID", source.Status.ID, "Name", source.Spec.Name, "Link", source.Spec.Link)

	// If source ID is not set, create a new source in the news aggregator
	if source.Status.ID == 0 {
		return r.createSource(ctx, &source)
	}
	// Update the existing source in the news aggregator
	return r.updateSource(ctx, source.Status.ID, &source)
}

// createSource creates a new source in the news aggregator service due to a new Source resource being created.
// - Marshals the SourceSpec into JSON and sends a POST request to the news aggregator service.
// - Updates the status of the Source resource based on the response.
// - Sets the Source ID in the Status field to keep track of the resource in the external service.
func (r *SourceReconciler) createSource(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Creating source:", source.Spec.Name)

	// Marshal source spec into JSON
	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}
	logrus.Println("Source data:", string(sourceData))

	// Create an HTTP POST request to the news aggregator service
	req, err := http.NewRequest(http.MethodPost, r.NewsAggregatorSrcServiceURL, bytes.NewBuffer(sourceData))
	if err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Handle non-OK status codes
	if resp.StatusCode != http.StatusOK {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, resp.Status)
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, fmt.Errorf("failed to create source in news aggregator: %s", resp.Status)
	}

	// Parse the response body into a new SourceSpec
	var createdSource aggregatorv1.SourceSpec
	if err := json.NewDecoder(resp.Body).Decode(&createdSource); err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionFalse, ReasonFailedCr, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}

	// Update the Source resource status with the new ID
	source.Status.ID = createdSource.Id
	if err := r.Client.Status().Update(ctx, source); err != nil {
		return ctrl.Result{}, err
	}

	// Update the Source status with a successful creation condition
	err = r.updateSourceStatus(ctx, source, aggregatorv1.SourceAdded, metav1.ConditionTrue, ReasonSuccessCr, "Source successfully created in news aggregator")
	if err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully created source in news aggregator")
	return ctrl.Result{}, nil
}

// updateSource updates an existing source in the news aggregator service due to a Source resource being updated.
// - Marshals the updated SourceSpec into JSON and sends a PUT request to the news aggregator service.
// - Updates the status of the Source resource based on the response.
func (r *SourceReconciler) updateSource(ctx context.Context, sourceID int, source *aggregatorv1.Source) (ctrl.Result, error) {
	logrus.Println("Updating source:", source.Spec.Name)

	// Marshal source spec into JSON
	sourceData, err := json.Marshal(source.Spec)
	if err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}

	// Create an HTTP PUT request to update the source in the news aggregator service
	url := fmt.Sprintf("%s/%d", r.NewsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(sourceData))
	if err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, err.Error())
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Handle non-OK status codes
	if resp.StatusCode != http.StatusOK {
		errUp := r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionFalse, ReasonFailedUpd, resp.Status)
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, fmt.Errorf("failed to update source in news aggregator: %s", resp.Status)
	}

	// Update the Source status with a successful update condition
	err = r.updateSourceStatus(ctx, source, aggregatorv1.SourceUpdated, metav1.ConditionTrue, ReasonSuccessUpd, "Source successfully updated in news aggregator")
	if err != nil {
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully updated source in news aggregator")
	return ctrl.Result{}, nil
}

// deleteSource deletes a source from the news aggregator service.
// - Sends a DELETE request to the news aggregator service using the provided source ID.
func (r *SourceReconciler) deleteSource(sourceID int) (ctrl.Result, error) {
	logrus.Println("Deleting source with ID:", sourceID)

	// Create an HTTP DELETE request to remove the source from the news aggregator service
	url := fmt.Sprintf("%s/%d", r.NewsAggregatorSrcServiceURL, sourceID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return ctrl.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Handle non-OK status codes
	if resp.StatusCode != http.StatusOK {
		return ctrl.Result{}, fmt.Errorf("failed to delete source from news aggregator: %s", resp.Status)
	}

	logrus.Info("Successfully deleted source from news aggregator")
	return ctrl.Result{}, nil
}

// updateSourceStatus updates the SourceStatus of a source resource with the given condition.
// - Adds a new SourceCondition to the SourceStatus.Conditions list.
// - Updates the status of the Source resource in the Kubernetes API.
func (r *SourceReconciler) updateSourceStatus(ctx context.Context, source *aggregatorv1.Source, conditionType aggregatorv1.SourceConditionType, status metav1.ConditionStatus, reason, message string) error {
	newCondition := aggregatorv1.SourceCondition{
		Type:           conditionType,
		Status:         status,
		LastUpdateTime: metav1.Time{Time: time.Now()},
		Reason:         reason,
		Message:        message,
	}

	source.Status.SetCondition(newCondition)
	return r.Client.Status().Update(ctx, source)
}

// SetupWithManager sets up the controller with the Manager and uses predicates to filter events.
// - Registers the controller with the Manager for the Source resource.
// - Adds event filters to optimize reconciliation using the Source predicates.
func (r *SourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Source{}).
		WithEventFilter(predicates.Source()).
		Complete(r)
}
