package v1_test

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Source Webhook Tests", func() {
	var (
		source       *aggregatorv1.Source
		fakeClient   client.Client
		scheme       *runtime.Scheme
		ctx          context.Context
		namespace    = "default"
		existingName = "existing-source"
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(aggregatorv1.AddToScheme(scheme)).To(Succeed())

		// Create a fake client with the scheme and no initial objects
		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme).
			Build()

		// Inject the fake client into the SourceClient wrapper
		aggregatorv1.SourceClient = aggregatorv1.SourceClientWrapper{Client: fakeClient}

		source = &aggregatorv1.Source{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
			},
			Spec: aggregatorv1.SourceSpec{
				Name:      "ValidName",
				ShortName: "ValidShort",
				Link:      "http://example.com/rss",
			},
		}

		ctx = context.TODO()
	})

	Context("ValidateCreate", func() {
		It("should pass validation with valid fields", func() {
			warnings, err := source.ValidateCreate()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when name, short_name, or link is empty", func() {
			source.Spec.Name = ""
			warnings, err := source.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name must be present"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when name or short_name exceeds 20 characters", func() {
			source.Spec.Name = "ThisIsAVeryLongNameExceedingTheLimit"
			warnings, err := source.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot be more than 20 characters"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation when link is invalid", func() {
			source.Spec.Link = "invalid-url"
			warnings, err := source.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("link must be a valid URL"))
			Expect(warnings).To(BeNil())
		})
	})

	Context("ValidateUpdate", func() {
		It("should pass ValidateUpdate with valid fields", func() {
			oldSource := &aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "old-source",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					Name:      "OldName",
					ShortName: "OldShort",
					Link:      "http://old.com/rss",
				},
			}
			warnings, err := source.ValidateUpdate(oldSource)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail ValidateUpdate when new name exceeds 20 characters", func() {
			source.Spec.Name = "ThisIsAVeryLongNameExceedingTheLimit"
			oldSource := &aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "old-source",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					Name:      "OldName",
					ShortName: "OldShort",
					Link:      "http://old.com/rss",
				},
			}
			warnings, err := source.ValidateUpdate(oldSource)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot be more than 20 characters"))
			Expect(warnings).To(BeNil())
		})
	})

	Context("ValidateDelete", func() {
		It("should pass ValidateDelete without errors", func() {
			warnings, err := source.ValidateDelete()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})
	})

	Context("CheckUniqueFields", func() {
		BeforeEach(func() {
			// Create an existing source in the fake client
			existingSource := &aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      existingName,
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					Name:      "ExistingName",
					ShortName: "ExistingShort",
					Link:      "http://example.com/rss",
				},
			}
			Expect(fakeClient.Create(ctx, existingSource)).To(Succeed())
		})

		It("should pass uniqueness check with unique fields", func() {
			source.Spec.Name = "UniqueName"
			source.Spec.ShortName = "UniqueShort"
			source.Spec.Link = "http://unique.com/rss"

			warnings, err := source.ValidateCreate()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail uniqueness check with duplicate name", func() {
			source.Spec.Name = "ExistingName"
			warnings, err := source.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name must be unique in the namespace"))
			Expect(warnings).To(BeNil())
		})

		It("should fail uniqueness check with duplicate short_name", func() {
			source.Spec.ShortName = "ExistingShort"
			_, err := source.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unique in the namespace"))
		})

		It("should fail uniqueness check with duplicate link", func() {
			source.Spec.Link = "http://example.com/rss"
			warnings, err := source.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("link must be unique in the namespace"))
			Expect(warnings).To(BeNil())
		})
	})
})
