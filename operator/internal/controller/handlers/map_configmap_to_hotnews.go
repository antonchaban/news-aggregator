package handlers

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MapConfigMapToHotNews maps ConfigMap updates to HotNews resources in the same namespace.
func MapConfigMapToHotNews(c client.Client, configMapName string) handler.MapFunc {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		logrus.Println("Mapping ConfigMap to HotNews")

		// Only proceed if the ConfigMap is the one we're interested in
		if object.GetName() != configMapName {
			return nil
		}

		namespace := object.GetNamespace()

		var hotNewsList aggregatorv1.HotNewsList
		// List all HotNews resources in the same namespace as the ConfigMap
		if err := c.List(ctx, &hotNewsList, &client.ListOptions{Namespace: namespace}); err != nil {
			logrus.Errorf("Failed to list HotNews resources: %v", err)
			return nil
		}

		var requests []reconcile.Request

		for _, hotNews := range hotNewsList.Items {
			// If the HotNews uses any feedGroups, we need to reconcile it
			if len(hotNews.Spec.FeedGroups) > 0 {
				requests = append(requests, reconcile.Request{
					NamespacedName: client.ObjectKey{
						Namespace: hotNews.Namespace,
						Name:      hotNews.Name,
					},
				})
			}
		}
		return requests
	}
}
