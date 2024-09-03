package controller_test

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	"context"
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Create source tests", func() {
	const (
		resourceName = "test-resource"
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

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&aggregatorv1.Source{}).
			Build()
	})

	AfterEach(func() {
		server.Close()
	})

	It("should successfully create the resource in the external service", func() {
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

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).ToNot(HaveOccurred())

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

	It("should handle creation failure, when error in request to server", func() {
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

		server.Close() // Simulate failure by closing the server
		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceAdded))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonFailedCr))
	})

	It("should handle creation failure, when error in resp.StatusCode", func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec aggregatorv1.SourceSpec
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				spec.Id = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusMethodNotAllowed)
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

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceAdded))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonFailedCr))
	})

	It("should add finalizers if not present", func() {
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

		source.Finalizers = []string{}

		reconciler = controller.SourceReconciler{
			Client:                      fakeClient,
			Scheme:                      scheme.Scheme,
			HTTPClient:                  server.Client(),
			NewsAggregatorSrcServiceURL: server.URL,
		}

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())

		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).NotTo(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedSource.Finalizers).To(ContainElement(controller.SrcFinalizer))
	})

	It("should handle creation failure, when error in decoding", func() {
		// Create fake HTTP server to mock an external service
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec int64
				spec = 1
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

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

		server.Close() // Simulate failure by closing the server
		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceAdded))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonFailedCr))
	})
})

var _ = Describe("Update source tests", func() {
	const (
		resourceName = "test-resource"
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
				ID: 1,
			},
		}

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&aggregatorv1.Source{}).
			Build()
	})

	AfterEach(func() {})

	It("should successfully update the resource in the external service", func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec int64
				spec = 1
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
			} else if r.Method == http.MethodPut {
				var spec aggregatorv1.SourceSpec
				spec.Name = "Updated Source"
				spec.Link = "http://example.com/rss"
				spec.Id = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
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

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).ToNot(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceUpdated))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionTrue))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonSuccessUpd))
	})

	It("should handle upd failure, when error in request to server", func() {
		// Create fake HTTP server to mock an external service
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec int64
				spec = 1
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
			} else if r.Method == http.MethodPut {
				var spec aggregatorv1.SourceSpec
				spec.Name = "Updated Source"
				spec.Link = "http://example.com/rss"
				spec.Id = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
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

		server.Close() // Simulate failure by closing the server
		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceUpdated))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonFailedUpd))
	})

	It("should handle upd failure, when error in resp.StatusCode", func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec int64
				spec = 1
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
			} else if r.Method == http.MethodPut {
				var spec aggregatorv1.SourceSpec
				spec.Name = "Updated Source"
				spec.Link = "http://example.com/rss"
				spec.Id = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write(respData)
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

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceUpdated))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonFailedUpd))
	})

	It("should handle upd failure, when error in decoding", func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec int64
				spec = 1
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
			} else if r.Method == http.MethodPut {
				var spec int64
				spec = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
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

		server.Close() // Simulate failure by closing the server
		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())

		updatedSource := &aggregatorv1.Source{}
		err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedSource)
		Expect(err).NotTo(HaveOccurred())

		conditions := updatedSource.Status.Conditions
		Expect(len(conditions)).To(Equal(1))
		Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceUpdated))
		Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
		Expect(conditions[0].Reason).To(Equal(controller.ReasonFailedUpd))
	})

})

var _ = Describe("Delete source tests", func() {
	const (
		resourceName = "test-resource"
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
				Name:       resourceName,
				Namespace:  "default",
				Finalizers: []string{controller.SrcFinalizer},
			},
			Spec: aggregatorv1.SourceSpec{
				Name: "Test Source",
				Link: "http://example.com/rss",
			},
			Status: aggregatorv1.SourceStatus{
				ID: 1,
			},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	It("should successfully delete the resource from the external service", func() {
		// Create fake HTTP server to mock an external service
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&aggregatorv1.Source{}).
			Build()

		reconciler = controller.SourceReconciler{
			Client:                      fakeClient,
			Scheme:                      scheme.Scheme,
			HTTPClient:                  server.Client(),
			NewsAggregatorSrcServiceURL: server.URL,
		}

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		Expect(fakeClient.Delete(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			err := fakeClient.Get(context.TODO(), typeNamespacedName, &source)
			return errors.IsNotFound(err)
		}).Should(BeTrue())
	})

	It("should handle deletion failure, when error in request to server", func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&aggregatorv1.Source{}).
			Build()

		reconciler = controller.SourceReconciler{
			Client:                      fakeClient,
			Scheme:                      scheme.Scheme,
			HTTPClient:                  server.Client(),
			NewsAggregatorSrcServiceURL: server.URL,
		}

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		Expect(fakeClient.Delete(context.TODO(), &source)).ToNot(HaveOccurred())
		server.Close() // Simulate failure by closing the server
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())
	})

	It("should handle deletion failure, when error in resp.StatusCode", func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusMethodNotAllowed)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))

		fakeClient = fake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithStatusSubresource(&aggregatorv1.Source{}).
			Build()

		reconciler = controller.SourceReconciler{
			Client:                      fakeClient,
			Scheme:                      scheme.Scheme,
			HTTPClient:                  server.Client(),
			NewsAggregatorSrcServiceURL: server.URL,
		}

		Expect(fakeClient.Create(context.TODO(), &source)).ToNot(HaveOccurred())
		Expect(fakeClient.Delete(context.TODO(), &source)).ToNot(HaveOccurred())
		_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
			NamespacedName: typeNamespacedName,
		})
		Expect(err).To(HaveOccurred())
	})
})
