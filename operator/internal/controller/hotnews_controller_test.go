package controller

/*
import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

func TestHotNewsReconcile(t *testing.T) {
	// Create a new scheme and add the necessary schemes
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	// Create the initial HotNews object
	initialHotNews := &aggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-hotnews",
			Namespace: "default",
		},
		Spec: aggregatorv1.HotNewsSpec{
			Keywords:   []string{"news"},
			FeedGroups: []string{"group1"},
		},
		Status: aggregatorv1.HotNewsStatus{
			ArticlesCount:  1,
			NewsLink:       "http://127.0.0.1:59641?keywords=news&sources=bbc,abcnews",
			ArticlesTitles: []string{},
		},
	}

	// Create a ConfigMap with feed groups
	initialConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "feed-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"group1": "bbc,abcnews",
		},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(initialHotNews, initialConfigMap).Build()

	// Mock HTTP server to simulate external news service
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `[{"id": 1, "Title": "Mock Article", "Description": "Mock Description", "Link": "http://mock-link", "Source": {"id": 1, "name": "BBC News", "link": "https://feeds.bbci.co.uk/news/rss.xml", "short_name": "bbc"}, "PubDate": "2024-08-13T00:00:00Z"}]`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	// Create the reconciler with the mock server URL
	r := &HotNewsReconciler{
		Client:        k8sClient,
		Scheme:        scheme,
		HTTPClient:    http.DefaultClient,
		ArticleSvcURL: mockServer.URL,
		ConfigMapName: "feed-config",
	}

	// Create the reconcile request
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-hotnews",
			Namespace: "default",
		},
	}

	// Perform the reconciliation
	_, err := r.Reconcile(context.Background(), req)

	// Retrieve the HotNews object after reconciliation
	updatedHotNews := &aggregatorv1.HotNews{}
	err = k8sClient.Get(context.Background(), req.NamespacedName, updatedHotNews)
	assert.NoError(t, err, "HotNews object should be found after reconciliation")

	// Verify that the status was updated correctly
	assert.Equal(t, 1, updatedHotNews.Status.ArticlesCount, "Articles count should be 1 after reconciliation")
	assert.NotEmpty(t, updatedHotNews.Status.NewsLink, "NewsLink should be set after reconciliation")
}
*/
