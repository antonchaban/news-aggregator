package v1_test

import (
	"com.teamdev/news-aggregator/api/v1"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ = Describe("CfgMapValidatorWebHook Tests", func() {
	var (
		fakeClient client.Client
		scheme     *runtime.Scheme
		ctx        context.Context
		validator  *v1.CfgMapValidatorWebHook
		configMap  *corev1.ConfigMap
		namespace  = "news-alligator"
	)
	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(corev1.AddToScheme(scheme)).To(Succeed())
		Expect(v1.AddToScheme(scheme)).To(Succeed())

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme).
			Build()

		ctx = context.TODO()

		validator = &v1.CfgMapValidatorWebHook{
			Client:     fakeClient,
			CfgMapName: "test-configmap",
		}

		configMap = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "test-configmap",
			},
			Data: map[string]string{
				"source1": "ValidSource",
			},
		}

		source := &v1.Source{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "valid-source",
			},
			Spec: v1.SourceSpec{
				Name:      "ValidName",
				ShortName: "ValidSource",
			},
		}

		Expect(fakeClient.Create(ctx, source)).To(Succeed())
		logrus.SetLevel(logrus.DebugLevel)
	})

	Context("ValidateCreate", func() {
		It("should pass validation when ConfigMap has valid data", func() {
			warnings, err := validator.ValidateCreate(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when ConfigMap has no data", func() {
			configMap.Data = map[string]string{}
			warnings, err := validator.ValidateCreate(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ConfigMap news-alligator/test-configmap has no data"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when source in ConfigMap does not exist", func() {
			configMap.Data = map[string]string{
				"source1": "NonExistentSource",
			}
			warnings, err := validator.ValidateCreate(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`Invalid value: "NonExistentSource": source not found`))
			Expect(warnings).To(BeNil())
		})

		It("should skip validation for non-target ConfigMap", func() {
			configMap.Name = "other-configmap"
			warnings, err := validator.ValidateCreate(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})
	})

	Context("ValidateUpdate", func() {
		It("should pass validation when updated ConfigMap has valid data", func() {
			warnings, err := validator.ValidateUpdate(ctx, configMap, configMap)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when updated ConfigMap has invalid data", func() {
			configMap.Data = map[string]string{
				"source1": "InvalidSource",
			}

			warnings, err := validator.ValidateUpdate(ctx, configMap, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`Invalid value: "InvalidSource": source not found`))
			Expect(warnings).To(BeNil())
		})
	})

	Context("ValidateDelete", func() {
		It("should pass validation when deleting a ConfigMap", func() {
			warnings, err := validator.ValidateDelete(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail deletion, because configmap is in use", func() {
			hotNews := &v1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "hot-news",
				},
				Spec: v1.HotNewsSpec{
					FeedGroups: []string{"source1"},
				},
			}
			Expect(fakeClient.Create(ctx, hotNews)).To(Succeed())

			warnings, err := validator.ValidateDelete(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot delete ConfigMap news-alligator/test-configmap because HotNews news-alligator/hot-news depends on it"))
			Expect(warnings).To(BeNil())
		})

	})

	Context("validate function", func() {
		It("should return error when passed object is not a ConfigMap", func() {
			warnings, err := validator.ValidateCreate(ctx, &corev1.Pod{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected a ConfigMap but got a *v1.Pod"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when no data found in ConfigMap", func() {
			configMap.Data = map[string]string{}
			warnings, err := validator.ValidateCreate(ctx, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ConfigMap news-alligator/test-configmap has no data"))
			Expect(warnings).To(BeNil())
		})

		It("should skip validation when ConfigMap is not the target one", func() {
			validator.CfgMapName = "other-configmap"
			warnings, err := validator.ValidateCreate(ctx, configMap)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})
	})
})

var _ = Describe("CfgMapValidatorWebHook SetupWebhookWithManager", func() {
	var (
		mgr       manager.Manager
		validator *v1.CfgMapValidatorWebHook
		scheme    *runtime.Scheme
		setupErr  error
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(corev1.AddToScheme(scheme)).To(Succeed())
		Expect(v1.AddToScheme(scheme)).To(Succeed())
		mgr, setupErr = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme: scheme,
		})
		Expect(setupErr).ToNot(HaveOccurred())

		validator = &v1.CfgMapValidatorWebHook{}
	})

	Context("SetupWebhookWithManager", func() {
		It("should set up the webhook without errors", func() {
			err := validator.SetupWebhookWithManager(mgr)
			Expect(err).ToNot(HaveOccurred())

			webhookServer := mgr.GetWebhookServer()
			Expect(webhookServer).ToNot(BeNil())
		})

		It("should panic if the manager is nil", func() {
			defer func() {
				if r := recover(); r != nil {
					Expect(r).To(HaveOccurred())
				} else {
					Fail("Expected panic, but code did not panic")
				}
			}()
			err := validator.SetupWebhookWithManager(nil)
			if err != nil {
				return
			}
		})
	})
})
