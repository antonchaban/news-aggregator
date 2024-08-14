package v1

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"time"
)

// log is for logging in this package.
var hotnewslog = logf.Log.WithName("hotnews-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient() // Set the global k8sClient variable
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-hotnews,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HotNews{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HotNews) Default() {
	hotnewslog.Info("default", "name", r.Name)
	if r.Spec.SummaryConfig.TitlesCount == 0 {
		r.Spec.SummaryConfig.TitlesCount = 10
	}
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update,versions=v1,name=vhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	hotnewslog.Info("validate create", "name", r.Name)

	if warnings, err := r.validateHotNews(); err != nil {
		return warnings, err
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	hotnewslog.Info("validate update", "name", r.Name)

	if warnings, err := r.validateHotNews(); err != nil {
		return warnings, err
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	hotnewslog.Info("validate delete", "name", r.Name)

	// No specific validation needed on delete for now.
	return nil, nil
}

func (r *HotNews) validateHotNews() (admission.Warnings, error) {
	var allErrs field.ErrorList
	var warnings admission.Warnings

	// Ensure keywords are present
	if len(r.Spec.Keywords) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("keywords"), "keywords must be present"))
	}

	// Validate dates
	if r.Spec.DateStart != "" && r.Spec.DateEnd != "" {
		startTime, err := time.Parse("2006-01-02", r.Spec.DateStart)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("dateStart"), r.Spec.DateStart, "invalid date format, should be YYYY-MM-DD"))
		}
		endTime, err := time.Parse("2006-01-02", r.Spec.DateEnd)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("dateEnd"), r.Spec.DateEnd, "invalid date format, should be YYYY-MM-DD"))
		}
		if err == nil && startTime.After(endTime) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("dateStart"), r.Spec.DateStart, "dateStart must be before dateEnd"))
		}
	} else {
		if r.Spec.DateStart != "" {
			if _, err := time.Parse("2006-01-02", r.Spec.DateStart); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("dateStart"), r.Spec.DateStart, "invalid date format, should be YYYY-MM-DD"))
			}
		}
		if r.Spec.DateEnd != "" {
			if _, err := time.Parse("2006-01-02", r.Spec.DateEnd); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("dateEnd"), r.Spec.DateEnd, "invalid date format, should be YYYY-MM-DD"))
			}
		}
	}

	if len(r.Spec.Sources) > 0 {
		// Check that FeedGroups is empty if Sources is not empty
		if len(r.Spec.FeedGroups) > 0 {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec").Child("feedGroups"), "feedGroups cannot be used when sources are specified"))
		}

		sourceList := &SourceList{}
		err := k8sClient.List(context.Background(), sourceList, &client.ListOptions{Namespace: r.Namespace})
		logrus.Println("SourceList: ", sourceList.Items)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("sources"), r.Spec.Sources, "unable to fetch SourceList"))
		} else {
			validSources := make(map[string]bool)
			for _, source := range sourceList.Items {
				validSources[source.Spec.ShortName] = true
			}
			for i, source := range r.Spec.Sources {
				if !validSources[source] {
					allErrs = append(allErrs, field.NotFound(field.NewPath("spec").Child("sources").Index(i), source))
				}
			}
		}
	}

	if len(allErrs) > 0 {
		return warnings, errors.NewInvalid(GroupVersion.WithKind("HotNews").GroupKind(), r.Name, allErrs)
	}

	return nil, nil
}
