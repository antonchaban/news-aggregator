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

func TestMapConfigMapToHotNews(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aggregatorv1.AddToScheme(scheme)

	hotNewsWithFeedGroups := &aggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hotnews-with-feedgroups",
			Namespace: "default",
		},
		Spec: aggregatorv1.HotNewsSpec{
			FeedGroups: []string{"group1"},
		},
	}

	hotNewsWithoutFeedGroups := &aggregatorv1.HotNews{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hotnews-without-feedgroups",
			Namespace: "default",
		},
		Spec: aggregatorv1.HotNewsSpec{
			FeedGroups: []string{},
		},
	}

	c := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(hotNewsWithFeedGroups, hotNewsWithoutFeedGroups).
		Build()

	configMapName := "feed-group-source"
	configMapNamespace := "default"

	mapFunc := MapConfigMapToHotNews(c, configMapName, configMapNamespace)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: configMapNamespace,
		},
	}

	reqs := mapFunc(context.Background(), configMap)

	// expect that only hotNewsWithFeedGroups will be returned
	if len(reqs) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(reqs))
	}

	expectedRequest := reconcile.Request{
		NamespacedName: client.ObjectKey{
			Namespace: "default",
			Name:      "hotnews-with-feedgroups",
		},
	}

	if reqs[0] != expectedRequest {
		t.Errorf("Expected request %+v, got %+v", expectedRequest, reqs[0])
	}

	// Test with a ConfigMap that does not match the name or namespace
	otherConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "other-configmap",
			Namespace: "default",
		},
	}

	reqs = mapFunc(context.Background(), otherConfigMap)
	if len(reqs) != 0 {
		t.Errorf("Expected 0 requests, got %d", len(reqs))
	}
}
