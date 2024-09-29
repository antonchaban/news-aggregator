package v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate--v1-configmap,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=configmaps,verbs=create;update;delete,versions=v1,name=vconfigmap.kb.io,admissionReviewVersions=v1

// CfgMapValidatorWebHook validates a specific ConfigMap.
type CfgMapValidatorWebHook struct {
	Client     client.Client
	CfgMapName string
}

// ValidateCreate validates the creation of the ConfigMap, checks does it use existing sources.
func (v *CfgMapValidatorWebHook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Create")
	return v.validate(ctx, obj)
}

// ValidateUpdate validates the update of the ConfigMap, checks does the ConfigMap use existing sources.
func (v *CfgMapValidatorWebHook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Update")
	return v.validate(ctx, newObj)
}

// ValidateDelete function for validating the deletion of the ConfigMap.
// ValidateDelete validates the deletion of the ConfigMap.
// It prevents deletion if there are any HotNews resources that depend on it.
func (v *CfgMapValidatorWebHook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Delete")
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("expected a ConfigMap but got a %T", obj)
	}

	// Check if the ConfigMap is for feed-groups before validating
	if cm.Name != v.CfgMapName {
		logrus.Printf("ConfigMap %s/%s is not the target ConfigMap; skipping delete validation", cm.Namespace, cm.Name)
		return nil, nil
	}

	logrus.Printf("Validating deletion of ConfigMap %s/%s", cm.Namespace, cm.Name)

	var hotNewsList HotNewsList
	if err := v.Client.List(ctx, &hotNewsList); err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return nil, fmt.Errorf("failed to list HotNews resources: %v", err)
	}

	// Check if any HotNews resource depends on this ConfigMap
	for _, hotNews := range hotNewsList.Items {
		if len(hotNews.Spec.FeedGroups) > 0 {
			// HotNews depends on the ConfigMap
			errMsg := fmt.Sprintf("Cannot delete ConfigMap %s/%s because HotNews %s/%s depends on it", cm.Namespace, cm.Name, hotNews.Namespace, hotNews.Name)
			return nil, fmt.Errorf(errMsg)
		}
	}
	return nil, nil
}

// validate checks the ConfigMap for valid sources.
func (v *CfgMapValidatorWebHook) validate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("expected a ConfigMap but got a %T", obj)
	}

	// Check if the ConfigMap is with required name and namespace before validating
	if cm.Name != v.CfgMapName {
		logrus.Printf("ConfigMap %s/%s is not the target ConfigMap; skipping validation", cm.Namespace, cm.Name)
		return nil, nil
	}

	logrus.Printf("Validating ConfigMap %s/%s", cm.Namespace, cm.Name)

	// Check if the ConfigMap has data
	if len(cm.Data) == 0 {
		return nil, fmt.Errorf("ConfigMap %s/%s has no data", cm.Namespace, cm.Name)
	}

	// List the sources in the specified namespace
	var sourceList SourceList
	err := v.Client.List(ctx, &sourceList, &client.ListOptions{Namespace: cm.Namespace})
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %v", err)
	}

	// Build a set of valid source names for quick lookup
	validSources := make(map[string]struct{})
	for _, source := range sourceList.Items {
		validSources[source.Spec.ShortName] = struct{}{}
	}

	// Validate each key-value pair in the ConfigMap data
	var allErrs field.ErrorList
	for key, value := range cm.Data {
		errs := isValidSource(key, value, validSources)
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(corev1.SchemeGroupVersion.WithKind("ConfigMap").GroupKind(), cm.Name, allErrs)
	}

	return nil, nil
}

// isValidSource validates the sources in a ConfigMap entry.
func isValidSource(key string, value string, validSources map[string]struct{}) field.ErrorList {
	var errs field.ErrorList
	logrus.Printf("Validating key '%s' with value '%s'", key, value)
	sources := strings.Split(value, ",")
	for _, sourceName := range sources {
		sourceName = strings.TrimSpace(sourceName)
		if sourceName == "" {
			continue // Skip empty source names
		}
		if _, found := validSources[sourceName]; !found {
			errs = append(errs, field.Invalid(field.NewPath("data").Key(key), sourceName, "source not found"))
		}
	}
	return errs
}

// SetupWebhookWithManager registers the webhook with the manager.
func (v *CfgMapValidatorWebHook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithValidator(v).
		Complete()
}
