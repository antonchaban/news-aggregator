package v1

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// HotNewsClientWrapper wraps a Kubernetes client for use in validation functions.
type HotNewsClientWrapper struct {
	Client client.Client
}

// HotNewsClient is a global client wrapper for accessing Kubernetes resources.
var HotNewsClient HotNewsClientWrapper

const (
	dateFormat             = "2006-01-02"
	feedGroupConfigMapName = "feed-group-source"
)

// SetupWebhookWithManager sets up the webhook with the controller manager.
// It initializes the HotNewsClient and registers the webhook for the HotNews resource.
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	HotNewsClient.Client = mgr.GetClient() // Set the global client
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-aggregator-com-teamdev-v1-hotnews,mutating=true,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update;delete,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HotNews{}

// Default sets default values for the HotNews resource.
// It ensures that TitlesCount is set to a default value if not provided.
func (r *HotNews) Default() {
	logrus.Info("default", "name", r.Name)
	if r.Spec.SummaryConfig.TitlesCount == 0 {
		r.Spec.SummaryConfig.TitlesCount = 10
	}
}

// +kubebuilder:webhook:path=/validate-aggregator-com-teamdev-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=aggregator.com.teamdev,resources=hotnews,verbs=create;update;delete,versions=v1,name=vhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate validates the creation of a HotNews resource.
// It ensures that all required fields are present and valid.
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	logrus.Info("validate create", "name", r.Name)
	return r.validateHotNews()
}

// ValidateUpdate validates the update of a HotNews resource.
// It ensures that all required fields remain valid after an update.
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	logrus.Info("validate update", "name", r.Name)
	return r.validateHotNews()
}

// ValidateDelete validates the deletion of a HotNews resource.
// Currently, it does not perform any validation on deletion.
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	logrus.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateHotNews performs validation on the HotNews resource.
// It checks that required fields are present and valid.
func (r *HotNews) validateHotNews() (admission.Warnings, error) {
	var allErrs field.ErrorList

	// Ensure keywords are present
	if len(r.Spec.Keywords) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("keywords"), "keywords must be present"))
	}

	// Validate dates
	r.validateDate(r.Spec.DateStart, r.Spec.DateEnd, &allErrs)

	// Validate sources
	r.validateSrc(r.Spec.Sources, &allErrs)

	// Validate feed groups
	r.validateFeedGroups(r.Spec.FeedGroups, &allErrs)

	if len(allErrs) > 0 {
		return nil, errors.NewInvalid(GroupVersion.WithKind("HotNews").GroupKind(), r.Name, allErrs)
	}
	return nil, nil
}

// validateDate validates the dateStart and dateEnd fields.
// It ensures that the dates are in the correct format and that dateStart is before dateEnd.
func (r *HotNews) validateDate(dateStart, dateEnd string, allErrs *field.ErrorList) {
	var startTime, endTime time.Time
	var err error

	if dateStart != "" {
		startTime, err = time.Parse(dateFormat, dateStart)
		if err != nil {
			*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("dateStart"), dateStart, "invalid date format, should be YYYY-MM-DD"))
		}
	}

	if dateEnd != "" {
		endTime, err = time.Parse(dateFormat, dateEnd)
		if err != nil {
			*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("dateEnd"), dateEnd, "invalid date format, should be YYYY-MM-DD"))
		}
	}

	// If both dates are provided and valid, check their order
	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("dateStart"), dateStart, "dateStart must be before dateEnd"))
		}
	}
}

// validateSrc validates the sources field.
// It ensures that the specified sources exist in the cluster.
func (r *HotNews) validateSrc(sources []string, allErrs *field.ErrorList) {
	if len(sources) <= 0 {
		return
	}

	sourceList := &SourceList{}
	err := HotNewsClient.Client.List(context.Background(), sourceList, &client.ListOptions{Namespace: r.Namespace})
	logrus.Println("SourceList: ", sourceList.Items)
	if err != nil {
		*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("sources"), r.Spec.Sources, "unable to fetch SourceList"))
		return
	}

	validSources := make(map[string]bool)
	for _, source := range sourceList.Items {
		validSources[source.Spec.ShortName] = true
	}

	for i, source := range r.Spec.Sources {
		if !validSources[source] {
			*allErrs = append(*allErrs, field.NotFound(field.NewPath("spec").Child("sources").Index(i), source))
		}
	}
}

// validateFeedGroups validates that the specified feed groups exist in the ConfigMap.
func (r *HotNews) validateFeedGroups(feedGroups []string, allErrs *field.ErrorList) {
	if len(feedGroups) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	configMap := &corev1.ConfigMap{}
	err := HotNewsClient.Client.Get(ctx, client.ObjectKey{
		Namespace: r.Namespace,
		Name:      feedGroupConfigMapName,
	}, configMap)
	if err != nil {
		*allErrs = append(*allErrs, field.Invalid(field.NewPath("spec").Child("feedGroups"), feedGroups, fmt.Sprintf("unable to fetch ConfigMap %s/%s: %v", r.Namespace, feedGroupConfigMapName, err)))
		return
	}
	for i, group := range feedGroups {
		if _, found := configMap.Data[group]; !found {
			*allErrs = append(*allErrs, field.NotFound(field.NewPath("spec").Child("feedGroups").Index(i), group))
		}
	}
}
