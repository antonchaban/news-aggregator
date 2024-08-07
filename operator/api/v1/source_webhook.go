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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
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

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-source,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=sources,verbs=create;update,versions=v1,name=vsource.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Source{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Source) ValidateCreate() (admission.Warnings, error) {
	sourcelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Source) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	sourcelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Source) ValidateDelete() (admission.Warnings, error) {
	sourcelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
