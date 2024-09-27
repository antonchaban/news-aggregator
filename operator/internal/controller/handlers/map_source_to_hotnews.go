package handlers

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// MapSourceToHotNews maps Source updates to HotNews resources.
func MapSourceToHotNews(c client.Client, fetchConfigMap func(context.Context) (*corev1.ConfigMap, error)) handler.MapFunc {
	return func(ctx context.Context, object client.Object) []reconcile.Request {
		logrus.Println("Mapping Source to HotNews")
		source := object.(*aggregatorv1.Source)

		var hotNewsList aggregatorv1.HotNewsList
		if err := c.List(ctx, &hotNewsList); err != nil {
			logrus.Errorf("Failed to list HotNews resources: %v", err)
			return nil
		}

		// Fetch the ConfigMap containing feed groups
		configMap, err := fetchConfigMap(ctx)
		if err != nil {
			logrus.Warn("ConfigMap not found or error fetching it")
			configMap = nil // Proceed without ConfigMap
		}

		var requests []reconcile.Request
		for _, hotNews := range hotNewsList.Items {
			// Combine sources from Spec.Sources and resolved FeedGroups
			combinedSources := make([]string, 0)
			combinedSources = append(combinedSources, hotNews.Spec.Sources...)

			if len(hotNews.Spec.FeedGroups) > 0 && configMap != nil {
				resolvedSources := resolveFeedGroups(hotNews.Spec.FeedGroups, configMap)
				combinedSources = append(combinedSources, resolvedSources...)
			}

			// Check if the updated Source's ShortName is in combinedSources
			if contains(combinedSources, source.Spec.ShortName) {
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

// resolveFeedGroups resolves feed groups to feed names using the ConfigMap.
func resolveFeedGroups(feedGroups []string, configMap *corev1.ConfigMap) []string {
	var feedNames []string
	for _, group := range feedGroups {
		if feeds, found := configMap.Data[group]; found {
			feedNames = append(feedNames, strings.Split(feeds, ",")...)
		}
	}
	return feedNames
}

// contains checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
