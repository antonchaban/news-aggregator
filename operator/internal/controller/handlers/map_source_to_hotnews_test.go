package handlers

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestMapSourceToHotNews(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	hotNews1 := &aggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hotnews1",
			Namespace: "default",
		},
		Spec: aggregatorv1.HotNewsSpec{
			Sources:    []string{"source1"},
			FeedGroups: []string{"group1"},
		},
	}

	// Create a fake client with the HotNews resource
	c := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(hotNews1).
		Build()

	fetchConfigMap := func(ctx context.Context) (*corev1.ConfigMap, error) {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "feed-group-source",
				Namespace: "default",
			},
			Data: map[string]string{
				"group1": "source2,source3",
			},
		}
		return configMap, nil
	}

	mapFunc := MapSourceToHotNews(c, fetchConfigMap)

	source1 := &aggregatorv1.Source{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "source1",
			Namespace: "default",
		},
		Spec: aggregatorv1.SourceSpec{
			ShortName: "source1",
		},
	}

	reqs := mapFunc(context.Background(), source1)

	// expect that hotNews1 will be returned because it specifies source1 in its Sources
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(reqs))
	}

	expectedRequest := reconcile.Request{
		NamespacedName: client.ObjectKey{
			Namespace: "default",
			Name:      "hotnews1",
		},
	}

	if reqs[0] != expectedRequest {
		t.Errorf("Expected request %+v, got %+v", expectedRequest, reqs[0])
	}

	source2 := &aggregatorv1.Source{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "source2",
			Namespace: "default",
		},
		Spec: aggregatorv1.SourceSpec{
			ShortName: "source2",
		},
	}

	reqs = mapFunc(context.Background(), source2)

	// expect that hotNews1 will be returned because source2 is in the resolved FeedGroups
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(reqs))
	}

	if reqs[0] != expectedRequest {
		t.Errorf("Expected request %+v, got %+v", expectedRequest, reqs[0])
	}

	source3 := &aggregatorv1.Source{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "source4",
			Namespace: "default",
		},
		Spec: aggregatorv1.SourceSpec{
			ShortName: "source4",
		},
	}

	reqs = mapFunc(context.Background(), source3)

	// expect that no HotNews will be returned
	if len(reqs) != 0 {
		t.Errorf("Expected 0 requests, got %d", len(reqs))
	}
}
