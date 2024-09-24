package controller_test

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller"
	"com.teamdev/news-aggregator/internal/controller/models"
	"context"
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"
)

var _ = Describe("HotNewsReconciler Tests", func() {
	const (
		hotNewsName = "test-hotnews"
		namespace   = "default"
	)

	var (
		typeNamespacedName types.NamespacedName
		reconciler         controller.HotNewsReconciler
		server             *httptest.Server
		fakeClient         client.Client
		hotNews            aggregatorv1.HotNews
		configMap          corev1.ConfigMap
	)

	BeforeEach(func() {
		typeNamespacedName = types.NamespacedName{
			Name:      hotNewsName,
			Namespace: namespace,
		}
		hotNews = aggregatorv1.HotNews{
			ObjectMeta: metav1.ObjectMeta{
				Name:      hotNewsName,
				Namespace: namespace,
			},
			Spec: aggregatorv1.HotNewsSpec{
				FeedGroups: []string{"group1"},
				Keywords:   []string{"breaking", "news"},
				DateStart:  "2023-01-01",
				DateEnd:    "2023-01-02",
				Sources:    []string{"source1"},
			},
		}

		configMap = corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "feed-config",
				Namespace: namespace,
			},
			Data: map[string]string{
				"group1": "source1,source2",
			},
		}
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
		}
	})

	Context("Reconciliation", func() {
		It("should successfully reconcile a HotNews resource", func() {
			source1 := aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "source1",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					ShortName: "source1",
				},
			}

			source2 := aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "source2",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					ShortName: "source2",
				},
			}

			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap, &source1, &source2).
				Build()

			// Create a fake HTTP server to simulate the external service
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.String(), "/articles") {
					articles := []models.Article{
						{Title: "Breaking News", Link: "http://example.com/article1"},
						{Title: "More News", Link: "http://example.com/article2"},
					}
					respData, err := json.Marshal(articles)
					Expect(err).NotTo(HaveOccurred())
					w.WriteHeader(http.StatusOK)
					w.Write(respData)
				} else {
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				HTTPClient:         server.Client(),
				ArticleSvcURL:      server.URL + "/articles",
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			result, err := reconciler.Reconcile(context.TODO(), reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())

			updatedHotNews := &aggregatorv1.HotNews{}
			err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedHotNews)
			Expect(err).NotTo(HaveOccurred())

			// Check that the status was updated correctly
			Expect(updatedHotNews.Status.ArticlesCount).To(Equal(2))
			Expect(updatedHotNews.Status.NewsLink).To(ContainSubstring(server.URL + "/articles"))

			conditions := updatedHotNews.Status.Conditions
			Expect(len(conditions)).To(Equal(1))
			Expect(conditions[0].Type).To(Equal(aggregatorv1.HNewsAdded))
			Expect(conditions[0].Status).To(Equal(metav1.ConditionTrue))
		})

		It("should handle reconciliation failure due to missing ConfigMap", func() {
			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews).
				Build()

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				ArticleSvcURL:      "http://example.com/articles",
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			result, err := reconciler.Reconcile(context.TODO(), reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).To(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())

			updatedHotNews := &aggregatorv1.HotNews{}
			err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedHotNews)
			Expect(err).NotTo(HaveOccurred())

			conditions := updatedHotNews.Status.Conditions
			Expect(len(conditions)).To(Equal(1))
			Expect(conditions[0].Type).To(Equal(aggregatorv1.HNewsAdded))
			Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(conditions[0].Reason).To(Equal("ConfigMapNotFound"))
		})

		It("should handle HTTP request failure during reconciliation", func() {
			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap).
				Build()

			// Simulate HTTP failure by not starting the server
			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				HTTPClient:         http.DefaultClient, // DefaultClient without a server will fail
				ArticleSvcURL:      "http://localhost:1234/articles",
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			result, err := reconciler.Reconcile(context.TODO(), reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).To(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())

			updatedHotNews := &aggregatorv1.HotNews{}
			err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedHotNews)
			Expect(err).NotTo(HaveOccurred())

			conditions := updatedHotNews.Status.Conditions
			Expect(len(conditions)).To(Equal(1))
			Expect(conditions[0].Type).To(Equal(aggregatorv1.HNewsAdded))
			Expect(conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(conditions[0].Reason).To(Equal("FetchArticlesFailed"))
		})
	})

	Context("Finalizer Handling", func() {
		It("should add a finalizer if not present", func() {

			source1 := aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "source1",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					ShortName: "source1",
				},
			}

			source2 := aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "source2",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					ShortName: "source2",
				},
			}

			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap, &source1, &source2).
				Build()
			// Create a fake HTTP server to simulate the external service
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.String(), "/articles") {
					articles := []models.Article{
						{Title: "Breaking News", Link: "http://example.com/article1"},
						{Title: "More News", Link: "http://example.com/article2"},
					}
					respData, err := json.Marshal(articles)
					Expect(err).NotTo(HaveOccurred())
					w.WriteHeader(http.StatusOK)
					w.Write(respData)
				} else {
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				HTTPClient:         server.Client(),
				ArticleSvcURL:      server.URL + "/articles",
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			result, err := reconciler.Reconcile(context.TODO(), reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Requeue).To(BeFalse())

			updatedHotNews := &aggregatorv1.HotNews{}
			err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedHotNews)
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedHotNews.Finalizers).To(ContainElement(controller.HotNewsFinalizer))
		})

		It("should remove the finalizer when the resource is deleted", func() {
			hotNews = aggregatorv1.HotNews{
				ObjectMeta: metav1.ObjectMeta{
					Name:              hotNewsName,
					Namespace:         namespace,
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
					Finalizers:        []string{controller.HotNewsFinalizer},
				},
				Spec: aggregatorv1.HotNewsSpec{
					FeedGroups: []string{"group1"},
					Keywords:   []string{"breaking", "news"},
					DateStart:  "2023-01-01",
					DateEnd:    "2023-01-02",
					Sources:    []string{"source1"},
				},
			}

			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap).
				Build()

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				ArticleSvcURL:      "http://example.com/articles",
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			_, err := reconciler.Reconcile(context.TODO(), reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())

			updatedHotNews := &aggregatorv1.HotNews{}
			err = reconciler.Client.Get(context.TODO(), typeNamespacedName, updatedHotNews)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})

	Context("Mapping Functions", func() {
		It("should map ConfigMap updates to the appropriate HotNews resources", func() {
			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap).
				Build()

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			configMapUpdate := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "feed-config",
					Namespace: namespace,
				},
				Data: map[string]string{
					"group1": "source1",
				},
			}
			Expect(fakeClient.Update(context.TODO(), configMapUpdate)).To(Succeed())

			requests := reconciler.MapConfigMapToHotNews(context.TODO(), configMapUpdate)
			Expect(len(requests)).To(Equal(1))
			Expect(requests[0].NamespacedName.Name).To(Equal(hotNewsName))
		})

		It("should map Source updates to the appropriate HotNews resources", func() {
			source1 := aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "source1",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					ShortName: "source1",
				},
			}

			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap, &source1).
				Build()

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			latestSource := &aggregatorv1.Source{}
			err := fakeClient.Get(context.TODO(), types.NamespacedName{
				Name:      "source1",
				Namespace: namespace,
			}, latestSource)
			Expect(err).NotTo(HaveOccurred())

			latestSource.Spec.ShortName = "newsource1"
			Expect(fakeClient.Update(context.TODO(), latestSource)).To(Succeed())

			updatedHotNews := &aggregatorv1.HotNews{}
			err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: hotNewsName, Namespace: namespace}, updatedHotNews)
			Expect(err).NotTo(HaveOccurred())
			updatedHotNews.Spec.Sources = []string{"newsource1"}
			Expect(fakeClient.Update(context.TODO(), updatedHotNews)).To(Succeed())

			requests := reconciler.MapSourceToHotNews(context.TODO(), latestSource)
			Expect(len(requests)).To(Equal(1))
			Expect(requests[0].NamespacedName.Name).To(Equal(hotNewsName))
		})

		It("should not map unrelated ConfigMap updates", func() {
			unrelatedConfigMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unrelated-config",
					Namespace: namespace,
				},
				Data: map[string]string{
					"unrelatedKey": "unrelatedValue",
				},
			}

			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, unrelatedConfigMap).
				Build()

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			requests := reconciler.MapConfigMapToHotNews(context.TODO(), unrelatedConfigMap)
			Expect(len(requests)).To(Equal(0))
		})

		FIt("should handle source deletion properly", func() {
			source1 := aggregatorv1.Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "source1",
					Namespace: namespace,
				},
				Spec: aggregatorv1.SourceSpec{
					ShortName: "source1",
				},
			}

			fakeClient = fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithStatusSubresource(&aggregatorv1.HotNews{}).
				WithObjects(&hotNews, &configMap, &source1).
				Build()

			reconciler = controller.HotNewsReconciler{
				Client:             fakeClient,
				Scheme:             scheme.Scheme,
				ConfigMapName:      configMap.Name,
				ConfigMapNamespace: configMap.Namespace,
			}

			Expect(fakeClient.Delete(context.TODO(), &source1)).To(Succeed())

			requests := reconciler.MapSourceToHotNews(context.TODO(), &source1)
			Expect(len(requests)).To(Equal(0))
		})
	})
})
