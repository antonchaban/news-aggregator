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

package v1

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"slices"
	"time"
)

// SourceClientWrapper allows substituting the Kubernetes client with a mock for testing.
type SourceClientWrapper struct {
	Client client.Client
}

// SourceClient is a global instance of SourceClientWrapper used to interact with the Kubernetes client.
var SourceClient SourceClientWrapper

// SetupWebhookWithManager sets up the webhook for the Source resource with the provided manager.
// It assigns the Kubernetes client from the manager to the SourceClient global variable and
// configures the webhook for the Source resource using the manager.
//
// Parameters:
// - mgr (ctrl.Manager): The manager that handles the Source resource and its webhooks.
//
// Returns:
// - error: Returns an error if the webhook setup fails.
func (r *Source) SetupWebhookWithManager(mgr ctrl.Manager) error {
	SourceClient.Client = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-source,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=sources,verbs=create;update,versions=v1,name=vsource.kb.io,admissionReviewVersions=v1

// ValidateCreate is a webhook method that validates the Source resource upon creation.
// It checks that the required fields (Name, ShortName, and Link) are non-empty and that their lengths
// do not exceed the specified limits. It also ensures that these fields are unique within the namespace.
//
// Returns:
// - admission.Warnings: Returns a list of warnings if applicable
// - error: Returns an error if validation fails.
func (r *Source) ValidateCreate() (admission.Warnings, error) {
	logrus.Println("validate create", "name", r.Name)
	return r.validateSource()
}

// ValidateUpdate is a webhook method that validates the Source resource upon update.
// It checks that the updated fields (Name, ShortName, and Link) are still non-empty, their lengths are within limits,
// and their values are unique within the namespace.
//
// Returns:
// - admission.Warnings: Returns a list of warnings if applicable (currently returns nil).
// - error: Returns an error if validation fails.
func (r *Source) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	logrus.Println("validate update", "name", r.Name)
	return r.validateSource()
}

// ValidateDelete is a webhook method that validates the Source resource upon deletion.
//
// Returns:
// - admission.Warnings: Returns a list of warnings if applicable
// - error: Returns an error if validation fails.
func (r *Source) ValidateDelete() (admission.Warnings, error) {
	logrus.Println("validate delete", "name", r.Name)

	// Check if the source is referenced by any HotNews resources
	var hotNewsList HotNewsList
	err := SourceClient.Client.List(context.Background(), &hotNewsList, &client.ListOptions{Namespace: r.Namespace})
	if err != nil {
		return nil, err
	}

	for _, hotNews := range hotNewsList.Items {
		if slices.Contains(hotNews.Spec.Sources, r.Spec.ShortName) {
			return nil, fmt.Errorf("cannot delete Source %s as it is referenced by HotNews %s", r.Spec.ShortName, hotNews.Name)
		}
	}

	return nil, nil
}

// validateSource performs the actual validation logic for the Source resource.
// It ensures that the Name, ShortName, and Link fields are non-empty, their lengths do not exceed 20 characters,
// and that the Link field contains a valid URL. It also checks that these fields are unique within the namespace.
//
// Returns:
// - admission.Warnings: Returns a list of warnings if applicable (currently returns nil).
// - error: Returns an error if validation fails.
func (r *Source) validateSource() (admission.Warnings, error) {
	var allErrs field.ErrorList
	var warnings admission.Warnings

	// Check if Name, ShortName, and Link are non-empty
	if len(r.Spec.Name) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("name"), "name must be present"))
	}
	if len(r.Spec.ShortName) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("shortName"), "short_name must be present"))
	}
	if len(r.Spec.Link) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("link"), "link must be present"))
	}

	// Check if Name and ShortName exceed the length limit
	if len(r.Spec.Name) > 20 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("name"), r.Spec.Name, "name cannot be more than 20 characters"))
	}
	if len(r.Spec.ShortName) > 20 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("shortName"), r.Spec.ShortName, "short_name cannot be more than 20 characters"))
	}

	// Check if the Link is a valid URL
	if !isValidURL(r.Spec.Link) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("link"), r.Spec.Link, "link must be a valid URL"))
	}

	// Check for uniqueness within the namespace
	if err := r.checkUniqueFields(SourceClient.Client, &allErrs); err != nil {
		allErrs = append(allErrs, field.InternalError(field.NewPath("spec"), fmt.Errorf("failed to check uniqueness: %v", err)))
	}

	// If there are accumulated errors, return them as an error
	if len(allErrs) > 0 {
		return warnings, errors.NewInvalid(GroupVersion.WithKind("Source").GroupKind(), r.Name, allErrs)
	}

	return warnings, nil
}

// isValidURL checks if the provided link is a valid URL.
// It parses the URL and ensures it has a valid scheme and host.
//
// Parameters:
// - link (string): The URL to be validated.
//
// Returns:
// - bool: Returns true if the URL is valid, false otherwise.
func isValidURL(link string) bool {
	u, err := url.Parse(link)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// checkUniqueFields checks if the Name, ShortName, and Link fields are unique within the namespace.
// It lists all Source resources in the namespace and compares their fields with the current Source resource.
//
// Parameters:
// - cl (client.Client): The Kubernetes client used to interact with the cluster.
// - allErrs (*field.ErrorList): A list of errors to append to if any uniqueness violations are found.
//
// Returns:
// - admission.Warnings: Returns a list of warnings if applicable (currently returns nil).
// - error: Returns an error if the Name, ShortName, or Link fields are not unique within the namespace.
func (r *Source) checkUniqueFields(cl client.Client, allErrs *field.ErrorList) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var sources SourceList
	logrus.Info("Listing sources in namespace", "namespace", r.Namespace)
	if err := cl.List(ctx, &sources, client.InNamespace(r.Namespace)); err != nil {
		return err
	}

	logrus.Info("Sources retrieved", "count", len(sources.Items))
	for _, source := range sources.Items {
		logrus.Info("Source found", "name", source.Name, "shortName", source.Spec.ShortName, "link", source.Spec.Link)
		if source.Name != r.Name {
			if source.Spec.Name == r.Spec.Name {
				*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("name"), r.Spec.Name, "name must be unique in the namespace"))
			}
			if source.Spec.ShortName == r.Spec.ShortName {
				*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("shortName"), r.Spec.ShortName, "short_name must be unique in the namespace"))
			}
			if source.Spec.Link == r.Spec.Link {
				*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("link"), r.Spec.Link, "link must be unique in the namespace"))
			}
		}
	}

	return nil
}
