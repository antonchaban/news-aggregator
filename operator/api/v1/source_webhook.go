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
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var sourcelog = logf.Log.WithName("source-resource")

// SourceClientWrapper allows substituting the Kubernetes client with a mock for testing.
type SourceClientWrapper struct {
	Client client.Client
}

var SourceClient SourceClientWrapper

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Source) SetupWebhookWithManager(mgr ctrl.Manager) error {
	SourceClient.Client = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-source,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=sources,verbs=create;update,versions=v1,name=msource.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Source{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
// This function able to set default values for the Source resource.
func (r *Source) Default() {
	sourcelog.Info("default", "name", r.Name)
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-source,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=sources,verbs=create;update,versions=v1,name=vsource.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Source{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
// This function validates the Source resource upon creation.
func (r *Source) ValidateCreate() (admission.Warnings, error) {
	sourcelog.Info("validate create", "name", r.Name)
	return r.validateSource()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
// This function validates the Source resource upon update.
func (r *Source) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	sourcelog.Info("validate update", "name", r.Name)
	return r.validateSource()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
// This function validates the Source resource upon deletion.
func (r *Source) ValidateDelete() (admission.Warnings, error) {
	sourcelog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateSource validates the fields of a Source resource.
// It ensures that the Name, ShortName, and Link fields are non-empty,
// the Name and ShortName fields are no longer than 20 characters,
// and the Link field is a valid URL.
// It also ensures that these fields are unique within the namespace.
func (r *Source) validateSource() (admission.Warnings, error) {
	if len(r.Spec.Name) == 0 || len(r.Spec.ShortName) == 0 || len(r.Spec.Link) == 0 {
		return nil, fmt.Errorf("name, short_name, and link fields cannot be empty")
	}

	if len(r.Spec.Name) > 20 || len(r.Spec.ShortName) > 20 {
		return nil, fmt.Errorf("name and short_name cannot be more than 20 characters")
	}

	if !isValidURL(r.Spec.Link) {
		return nil, fmt.Errorf("link field must be a valid URL")
	}

	// Check for uniqueness within the namespace
	return r.checkUniqueFields(SourceClient.Client)
}

// isValidURL checks if the provided link is a valid URL.
func isValidURL(link string) bool {
	u, err := url.Parse(link)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// checkUniqueFields ensures that the Name, ShortName, and Link fields are unique within the namespace.
func (r *Source) checkUniqueFields(cl client.Client) (admission.Warnings, error) {
	ctx := context.Background()

	var sources SourceList
	sourcelog.Info("Listing sources in namespace", "namespace", r.Namespace)
	if err := cl.List(ctx, &sources, client.InNamespace(r.Namespace)); err != nil {
		sourcelog.Error(err, "failed to list sources")
		return nil, err
	}

	sourcelog.Info("Sources retrieved", "count", len(sources.Items))
	for _, source := range sources.Items {
		sourcelog.Info("Source found", "name", source.Name, "shortName", source.Spec.ShortName, "link", source.Spec.Link)
		if source.Name != r.Name {
			if source.Spec.Name == r.Spec.Name {
				return nil, fmt.Errorf("name must be unique in the namespace")
			}
			if source.Spec.ShortName == r.Spec.ShortName {
				return admission.Warnings{}, fmt.Errorf("short_name must be unique in the namespace")
			}
			if source.Spec.Link == r.Spec.Link {
				return nil, fmt.Errorf("link must be unique in the namespace")
			}
		}
	}

	return nil, nil
}
