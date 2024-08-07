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

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Source) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-source,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=sources,verbs=create;update,versions=v1,name=msource.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Source{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Source) Default() {
	sourcelog.Info("default", "name", r.Name)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-source,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=sources,verbs=create;update,versions=v1,name=vsource.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Source{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Source) ValidateCreate() (admission.Warnings, error) {
	sourcelog.Info("validate create", "name", r.Name)
	return r.validateSource()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Source) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	sourcelog.Info("validate update", "name", r.Name)
	return r.validateSource()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Source) ValidateDelete() (admission.Warnings, error) {
	sourcelog.Info("validate delete", "name", r.Name)
	return nil, nil
}

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
	return r.checkUniqueFields()
}

func isValidURL(link string) bool {
	u, err := url.Parse(link)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (r *Source) checkUniqueFields() (admission.Warnings, error) {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()
	if err := AddToScheme(scheme); err != nil {
		sourcelog.Error(err, "failed to add scheme")
		return nil, err
	}
	cl, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		sourcelog.Error(err, "failed to create client")
		return nil, err
	}

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
				return nil, fmt.Errorf("short_name must be unique in the namespace")
			}
			if source.Spec.Link == r.Spec.Link {
				return nil, fmt.Errorf("link must be unique in the namespace")
			}
		}
	}

	return nil, nil
}
