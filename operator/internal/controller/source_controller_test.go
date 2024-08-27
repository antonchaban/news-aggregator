package controller_test

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	"context"
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TestSourceReconciler runs all the tests
var _ = Describe("Source Controller", func() {
	const (
		resourceName = "test-resource"
		//timeout      = time.Second * 30
		//interval     = time.Millisecond * 500
	)

	var (
		typeNamespacedName types.NamespacedName
		reconciler         controller.SourceReconciler
		server             *httptest.Server
		fakeClient         client.Client
		source             aggregatorv1.Source
	)

	BeforeEach(func() {

		typeNamespacedName = types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		source = aggregatorv1.Source{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: "default",
			},
			Spec: aggregatorv1.SourceSpec{
				Name: "Test Source",
				Link: "http://example.com/rss",
			},
			Status: aggregatorv1.SourceStatus{
				ID: 0,
			},
		}

		fakeClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).WithStatusSubresource(&aggregatorv1.Source{}).Build()

		// Create fake HTTP server to mock an external service
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec aggregatorv1.SourceSpec
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				spec.Id = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
			} else if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))

		reconciler = controller.SourceReconciler{
			Client:                      fakeClient,
			Scheme:                      scheme.Scheme,
			HTTPClient:                  server.Client(),
			NewsAggregatorSrcServiceURL: server.URL,
		}
	})

	AfterEach(func() {
	})

	FIt("should successfully create the resource in the external service", func() {

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).ToNot(HaveOccurred())
		//
		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedSource.Status.ID).To(Equal(1))

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceAdded))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionTrue))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonSuccessCr))
	})

	/*It("should successfully update the resource in the external service", func() {
		source.Status.ID = 1
		err := k8sClient.Status().Update(ctx, source)
		Expect(err).NotTo(HaveOccurred())

		source.Spec.Name = "Updated Source"
		err = k8sClient.Update(ctx, source)
		Expect(err).NotTo(HaveOccurred())

		_, err = sourceReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).NotTo(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = k8sClient.Get(ctx, typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceUpdated))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionTrue))
		Expect(conditions[0].Reason).To(Equal(ReasonSuccessUpd))
	})

	It("should successfully delete the resource from the external service", func() {
		source.Status.ID = 1
		source.Finalizers = append(source.Finalizers, SrcFinalizer)
		err := k8sClient.Update(ctx, source)
		Expect(err).NotTo(HaveOccurred())

		// Initiate the deletion
		err = k8sClient.Delete(ctx, source)
		Expect(err).NotTo(HaveOccurred())

		// Reconcile to handle finalizers
		_, err = sourceReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).NotTo(HaveOccurred())

		// Verify the source has been removed from the external service
		Eventually(func() bool {
			err := k8sClient.Get(ctx, typeNamespacedName, source)
			return errors.IsNotFound(err)
		}, timeout, interval).Should(BeTrue())

		deletedSource := &aggregatorv1.Source{}
		err = k8sClient.Get(ctx, typeNamespacedName, deletedSource)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})

	It("should add the finalizer if not present", func() {
		source.Finalizers = []string{}
		err := k8sClient.Update(ctx, source)
		Expect(err).NotTo(HaveOccurred())

		_, err = sourceReconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).NotTo(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = k8sClient.Get(ctx, typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedSource.Finalizers).To(ContainElement(SrcFinalizer))
	})*/

})
