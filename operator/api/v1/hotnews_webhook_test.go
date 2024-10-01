package v1_test

import (
	"com.teamdev/news-aggregator/api/v1"
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("HotNews Webhook Tests", func() {
	var (
		fakeClient client.Client
		scheme     *runtime.Scheme
		ctx        context.Context
		hotNews    *v1.HotNews
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(corev1.AddToScheme(scheme)).To(Succeed())
		Expect(v1.AddToScheme(scheme)).To(Succeed())

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme).
			Build()

		v1.HotNewsClient.Client = fakeClient

		ctx = context.TODO()
		source1 := &v1.Source{
			ObjectMeta: metav1.ObjectMeta{
				Name: "source1",
			},
			Spec: v1.SourceSpec{
				ShortName: "source1",
			},
		}

		source2 := &v1.Source{
			ObjectMeta: metav1.ObjectMeta{
				Name: "source2",
			},
			Spec: v1.SourceSpec{
				ShortName: "source2",
			},
		}

		Expect(fakeClient.Create(ctx, source1)).To(Succeed())
		Expect(fakeClient.Create(ctx, source2)).To(Succeed())

		hotNews = &v1.HotNews{
			Spec: v1.HotNewsSpec{
				Keywords:  []string{"breaking", "urgent"},
				DateStart: "2023-09-01",
				DateEnd:   "2023-09-10",
				Sources:   []string{"source1", "source2"},
				SummaryConfig: v1.SummaryConfig{
					TitlesCount: 5,
				},
			},
		}
	})

	Context("Defaulting", func() {
		It("should set TitlesCount to 10 if not provided", func() {
			hotNews.Spec.SummaryConfig.TitlesCount = 0
			hotNews.Default()
			Expect(hotNews.Spec.SummaryConfig.TitlesCount).To(Equal(10))
		})

		It("should not change TitlesCount if already set", func() {
			hotNews.Spec.SummaryConfig.TitlesCount = 5
			hotNews.Default()
			Expect(hotNews.Spec.SummaryConfig.TitlesCount).To(Equal(5))
		})
	})

	Context("ValidateCreate", func() {
		It("should pass validation with valid HotNews", func() {
			warnings, err := hotNews.ValidateCreate()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail validation with missing keywords", func() {
			hotNews.Spec.Keywords = nil
			warnings, err := hotNews.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("keywords must be present"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation with invalid date format", func() {
			hotNews.Spec.DateStart = "09-01-2023" // Invalid format
			warnings, err := hotNews.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid date format, should be YYYY-MM-DD"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation if dateStart is after dateEnd", func() {
			hotNews.Spec.DateStart = "2023-09-10"
			hotNews.Spec.DateEnd = "2023-09-01"
			warnings, err := hotNews.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("dateStart must be before dateEnd"))
			Expect(warnings).To(BeNil())
		})

		It("should append an error when List call fails", func() {
			// Replace the fake client with one that will simulate an error
			v1.HotNewsClient.Client = &errorFakeClient{}
			warnings, err := hotNews.ValidateCreate()

			// Assert that an error occurred and contains the expected message
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to fetch SourceList"))
			Expect(warnings).To(BeNil())
		})

		It("should fail validation if a source is not found", func() {
			existingSource := &v1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "existing-source",
					Namespace: "default",
				},
				Spec: v1.SourceSpec{
					ShortName: "source1",
				},
			}
			Expect(fakeClient.Create(ctx, existingSource)).To(Succeed())

			hotNews.Spec.Sources = []string{"source1", "invalid-source"}

			warnings, err := hotNews.ValidateCreate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid-source"))
			Expect(warnings).To(BeNil())
		})
	})

	Context("ValidateUpdate", func() {
		It("should pass validation with valid updated HotNews", func() {
			oldHotNews := hotNews.DeepCopy()
			warnings, err := hotNews.ValidateUpdate(oldHotNews)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("should fail validation on update with invalid dates", func() {
			oldHotNews := hotNews.DeepCopy()
			hotNews.Spec.DateStart = "2023-09-10"
			hotNews.Spec.DateEnd = "2023-09-01"
			warnings, err := hotNews.ValidateUpdate(oldHotNews)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("dateStart must be before dateEnd"))
			Expect(warnings).To(BeNil())
		})
	})

	Context("ValidateDelete", func() {
		It("should pass validation on delete without errors", func() {
			warnings, err := hotNews.ValidateDelete()
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeNil())
		})
	})
})

type errorFakeClient struct {
	client.Client
}

func (e *errorFakeClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return errors.New("simulated list error")
}
