package v1

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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
	if !ok {
		return nil, fmt.Errorf("expected a Config Map but got a %T", obj)
	}

	logrus.Println("Validating ConfigMap", cm.Name)
	// todo add validation logic here
	return nil, nil
}

func (v *CfgMapValidatorWebHook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Create")
	return v.validate(ctx, obj)
}

func (v *CfgMapValidatorWebHook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	logrus.Println("Validating ConfigMap Update")
	return v.validate(ctx, newObj)
}

func (v *CfgMapValidatorWebHook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return v.validate(ctx, obj)
}

func (v *CfgMapValidatorWebHook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		WithValidator(v).
		Complete()
}
