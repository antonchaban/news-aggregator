package v1

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// +kubebuilder:webhook:path=/validate--v1-configmap,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=configmaps,verbs=create;update,versions=v1,name=vconfigmap.kb.io,admissionReviewVersions=v1

type CfgMapValidatorWebHook struct {
	Client          client.Client
	CfgMapName      string
	CfgMapNamespace string
}

// validate admits a pod if a specific annotation exists.
func (v *CfgMapValidatorWebHook) validate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	var allErrs field.ErrorList
	var warnings admission.Warnings
	if !ok {
		return nil, fmt.Errorf("expected a Config Map but got a %T", obj)
	}

	logrus.Println("Validating ConfigMap", cm.Name)
	if cm.Data == nil || len(cm.Data) == 0 {
		return nil, fmt.Errorf("ConfigMap %s has no data", cm.Name)
	} else {
		var sourceList SourceList
		err := v.Client.List(ctx, &sourceList, &client.ListOptions{Namespace: v.CfgMapNamespace})
		if err != nil {
			return nil, fmt.Errorf("failed to list sources: %v", err)
		}

		logrus.Println("ConfigMap has data, beginning validation")
		for key, value := range cm.Data {
			logrus.Printf("Key: %s, Value: %s\n", key, value)
			sources := strings.Split(value, ",")
			logrus.Println("Sources: ", sources)
			if len(sources) == 0 {
				allErrs = append(allErrs, field.Invalid(field.NewPath(key), value, "No sources found"))
			} else {
				logrus.Printf("ConfigMap %s has sources: %v\n", cm.Name, sources)
				for _, sourceName := range sources {
					found := false
					for _, source := range sourceList.Items {
						if source.Spec.ShortName == sourceName {
							found = true
							break
						}
					}
					if !found {
						allErrs = append(allErrs, field.Invalid(field.NewPath(key), value, "Source "+sourceName+" not found"))
					}
				}
			}
		}

		if len(allErrs) > 0 {
			return warnings, fmt.Errorf("validation failed: %v", allErrs)
		}
	}
	return nil, nil
}

func (v *CfgMapValidatorWebHook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Create")
	return v.validate(ctx, obj)
}

func (v *CfgMapValidatorWebHook) ValidateUpdate(ctx context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Update")
	return v.validate(ctx, newObj)
}

func (v *CfgMapValidatorWebHook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *CfgMapValidatorWebHook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithValidator(v).
		Complete()
}
