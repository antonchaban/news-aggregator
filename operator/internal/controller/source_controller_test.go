package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
	"time"
)

func TestSourceReconcile(t *testing.T) {
	// Create a new scheme and add the necessary schemes
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	// Create the initial Source object
	initialSource := &aggregatorv1.Source{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-source",
			Namespace: "default",
		},
		Spec: aggregatorv1.SourceSpec{
			Name:      "Test Source",
			Link:      "https://example.com/rss",
			ShortName: "test",
		},
		Status: aggregatorv1.SourceStatus{
			ID: 123,
			Conditions: []aggregatorv1.SourceCondition{
				{
					Type:           aggregatorv1.SourceAdded,
					Status:         metav1.ConditionTrue,
					LastUpdateTime: metav1.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Reason:         "Added",
				},
			},
		},
	}
	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(initialSource).Build()

	// Set up an HTTP test server that returns a specific response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 123}`))
		}
	}))
	defer testServer.Close()

	r := &SourceReconciler{
		Client:                      k8sClient,
		Scheme:                      scheme,
		HTTPClient:                  testServer.Client(),
		NewsAggregatorSrcServiceURL: testServer.URL,
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-source",
			Namespace: "default",
		},
	}

	// Perform the reconciliation
	res, err := r.Reconcile(context.Background(), req)
	assert.False(t, res.Requeue)

	source := &aggregatorv1.Source{}
	err = k8sClient.Get(context.Background(), req.NamespacedName, source)
	assert.NoError(t, err, "Source object should be found after reconciliation")
	assert.Contains(t, source.Finalizers, "source.finalizers.teamdev.com", "Finalizer should be added")
	assert.Equal(t, 123, source.Status.ID, "Source ID should be updated to the value returned by the external service")
}
