package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller/handlers"
	"com.teamdev/news-aggregator/internal/controller/models"
	"com.teamdev/news-aggregator/internal/controller/predicates"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"slices"
	"strings"
	"time"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	Client        client.Client   // Client for interacting with the Kubernetes API.
	Scheme        *runtime.Scheme // Scheme for the reconciler.
	HTTPClient    *http.Client    // HTTP client for making external requests.
	ArticleSvcURL string          // URL of the news aggregator service.
	ConfigMapName string          // Name of the ConfigMap that contains feed groups.
}

const HotNewsFinalizer = "hotnews.finalizers.teamdev.com"

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// It receives specified HotNews resource and
// tries to receive articles from the news aggregator service due to specified parameters.
// If resource is being deleted, it removes owner references from the sources and removes finalizer.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logrus.Println("Reconciling HotNews")
	var hotNews aggregatorv1.HotNews
	if err := r.Client.Get(ctx, req.NamespacedName, &hotNews); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		logrus.Infof("HotNews resource not found, possibly deleted. Namespace: %s", req.Namespace)
		return ctrl.Result{}, nil
	}

	// Handle finalizer logic
	deleted, err := r.handleFinalizer(ctx, &hotNews)
	if err != nil {
		return ctrl.Result{}, err
	}
	if deleted {
		return ctrl.Result{}, nil
	}

	// Proceed with normal reconciliation for HotNews
	return r.reconcileHotNews(ctx, &hotNews)
}

// Manages the addition and removal of the finalizer on the HotNews resource.
// Ensures that owner references are cleaned up before the resource is deleted.
func (r *HotNewsReconciler) handleFinalizer(ctx context.Context, hotNews *aggregatorv1.HotNews) (bool, error) {
	if hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.isFinalizerPresent(ctx, hotNews)
	}
	return r.handleDeletion(ctx, hotNews)
}

// Checks if the finalizer is present on the HotNews resource.
func (r *HotNewsReconciler) isFinalizerPresent(ctx context.Context, hotNews *aggregatorv1.HotNews) (bool, error) {
	if slices.Contains(hotNews.Finalizers, HotNewsFinalizer) {
		return false, nil
	}

	hotNews.Finalizers = append(hotNews.Finalizers, HotNewsFinalizer)
	if err := r.Client.Update(ctx, hotNews); err != nil {
		return false, err
	}
	return false, nil
}

// Handles the deletion of the HotNews resource.
func (r *HotNewsReconciler) handleDeletion(ctx context.Context, hotNews *aggregatorv1.HotNews) (bool, error) {
	if !slices.Contains(hotNews.Finalizers, HotNewsFinalizer) {
		return true, nil
	}

	if err := r.deleteOwnerReferences(ctx, hotNews.Namespace, hotNews.Name); err != nil {
		return false, err
	}

	index := slices.Index(hotNews.Finalizers, HotNewsFinalizer)
	hotNews.Finalizers = slices.Delete(hotNews.Finalizers, index, index+1)
	if err := r.Client.Update(ctx, hotNews); err != nil {
		return false, err
	}
	return true, nil
}

// Removes the owner reference to the HotNews resource from all Source resources in the namespace.
// Ensures that Source resources are not left with dangling owner references after the HotNews resource is deleted.
func (r *HotNewsReconciler) deleteOwnerReferences(ctx context.Context, namespace, hotNewsName string) error {
	var sources aggregatorv1.SourceList
	if err := r.Client.List(ctx, &sources, &client.ListOptions{Namespace: namespace}); err != nil {
		return fmt.Errorf("failed to list Sources: %w", err)
	}

	for _, src := range sources.Items {
		if err := r.removeOwnerReferenceFromSource(ctx, &src, hotNewsName); err != nil {
			return err
		}
	}
	return nil
}

// Removes the owner reference to the HotNews resource from a single Source resource.
func (r *HotNewsReconciler) removeOwnerReferenceFromSource(ctx context.Context, src *aggregatorv1.Source, hotNewsName string) error {
	updatedOwnerRefs := removeOwnerReference(src.OwnerReferences, hotNewsName, "HotNews")
	if len(updatedOwnerRefs) == len(src.OwnerReferences) {
		// No change in owner references so skip update
		return nil
	}

	src.OwnerReferences = updatedOwnerRefs
	if err := r.Client.Update(ctx, src); err != nil {
		logrus.Errorf("Failed to update Source %s: %v", src.Name, err)
		return fmt.Errorf("failed to update Source %s: %w", src.Name, err)
	}
	logrus.Infof("Owner reference removed from Source %s", src.Name)
	return nil
}

// Utility function to remove a specific owner reference from a slice of owner references.
func removeOwnerReference(ownerRefs []metav1.OwnerReference, name, kind string) []metav1.OwnerReference {
	var updatedRefs []metav1.OwnerReference
	for _, ref := range ownerRefs {
		if ref.Name == name && ref.Kind == kind {
			continue
		}
		updatedRefs = append(updatedRefs, ref)
	}
	return updatedRefs
}

// Performs the main reconciliation logic for the HotNews resource.
// Resolves feed groups, sets owner references, fetches articles, and updates the status.
func (r *HotNewsReconciler) reconcileHotNews(ctx context.Context, hotNews *aggregatorv1.HotNews) (ctrl.Result, error) {
	configMap, err := r.fetchConfigMap(ctx, hotNews.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Warnf("ConfigMap %s not found in namespace %s, proceeding without resolving feed groups", r.ConfigMapName, hotNews.Namespace)
			r.updateStatusCondition(hotNews, metav1.ConditionFalse, "ConfigMapNotFound", fmt.Sprintf("ConfigMap %s not found in namespace %s", r.ConfigMapName, hotNews.Namespace))
			configMap = nil
		} else {
			logrus.Errorf("Failed to fetch ConfigMap: %v", err)
			r.updateStatusCondition(hotNews, metav1.ConditionFalse, "ConfigMapFetchFailed", "Failed to fetch ConfigMap")
			// Update the status before returning
			if err := r.Client.Status().Update(ctx, hotNews); err != nil {
				logrus.Error(err, "Failed to update HotNews status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}
	}

	// Combine sources from Spec.Sources and resolved FeedGroups
	combinedSources := append([]string{}, hotNews.Spec.Sources...)

	// Resolve FeedGroups to actual source names without modifying Spec.Sources
	if len(hotNews.Spec.FeedGroups) > 0 {
		if configMap != nil {
			resolvedSources := r.resolveFeedGroups(hotNews.Spec.FeedGroups, configMap)
			combinedSources = append(combinedSources, resolvedSources...)
		} else {
			logrus.Warn("ConfigMap not available, cannot resolve feed groups")
			r.updateStatusCondition(hotNews, metav1.ConditionFalse, "ConfigMapNotFound", "ConfigMap not found, feed groups cannot be resolved")
		}
	}

	// Set OwnerReferences for each source using combinedSources
	if err := r.setOwnerReferences(ctx, hotNews, combinedSources); err != nil {
		if errStatus := r.Client.Status().Update(ctx, hotNews); errStatus != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	// Build query parameters and fetch articles
	reqURL := fmt.Sprintf("%s?%s", r.ArticleSvcURL, r.buildQueryParams(hotNews, combinedSources))
	logrus.Println("Request URL:", reqURL)

	articles, err := r.fetchArticles(reqURL)
	if err != nil {
		logrus.Error(err, "Unable to fetch articles")
		r.updateStatusCondition(hotNews, metav1.ConditionFalse, "FetchArticlesFailed", "Failed to fetch articles")
		if err := r.Client.Status().Update(ctx, hotNews); err != nil { // upd status before return
			logrus.Error(err, "Failed to update HotNews status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	// Configure HotNews status fields
	r.updateHotNewsStatus(hotNews, articles, reqURL)
	r.updateStatusCondition(hotNews, metav1.ConditionTrue, "Reconciled", "Successfully reconciled HotNews")

	if err := r.Client.Status().Update(ctx, hotNews); err != nil {
		logrus.Error(err, "Failed to update HotNews status")
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
	return ctrl.Result{}, nil
}

// Retrieves the ConfigMap that contains feed group definitions.
func (r *HotNewsReconciler) fetchConfigMap(ctx context.Context, namespace string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: r.ConfigMapName}, configMap)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// Sets the HotNews resource as an owner reference on the relevant Source resources.
// Ensures that the Source resources are properly tracked and cannot be deleted while the HotNews resource exists.
func (r *HotNewsReconciler) setOwnerReferences(ctx context.Context, hotNews *aggregatorv1.HotNews, sources []string) error {
	// Fetch all sources once
	var sourceList aggregatorv1.SourceList
	if err := r.Client.List(ctx, &sourceList, &client.ListOptions{Namespace: hotNews.Namespace}); err != nil {
		logrus.Error(err, "Failed to list Sources")
		return err
	}

	// Build a map from ShortName to Source
	sourceMap := make(map[string]*aggregatorv1.Source)
	for i := range sourceList.Items {
		src := &sourceList.Items[i]
		sourceMap[src.Spec.ShortName] = src
	}

	for _, sourceShortName := range sources {
		source, found := sourceMap[sourceShortName]
		if !found {
			logrus.Errorf("Source with ShortName %s not found", sourceShortName)
			continue
		}

		// Add owner reference if not exists
		if !hasOwnerReference(source.OwnerReferences, hotNews.UID) {
			source.OwnerReferences = append(source.OwnerReferences, metav1.OwnerReference{
				APIVersion: hotNews.APIVersion,
				Kind:       hotNews.Kind,
				Name:       hotNews.Name,
				UID:        hotNews.UID,
			})
			if err := r.Client.Update(ctx, source); err != nil {
				return err
			}
		}
	}
	return nil
}

// hasOwnerReference checks if the owner reference already exists.
func hasOwnerReference(ownerRefs []metav1.OwnerReference, uid types.UID) bool {
	for _, owner := range ownerRefs {
		if owner.UID == uid {
			return true
		}
	}
	return false
}

// buildQueryParams constructs the query parameters for the HTTP request.
func (r *HotNewsReconciler) buildQueryParams(hotNews *aggregatorv1.HotNews, combinedSources []string) string {
	params := map[string]string{}
	if len(hotNews.Spec.Keywords) > 0 {
		params["keywords"] = strings.Join(hotNews.Spec.Keywords, ",")
	}
	if hotNews.Spec.DateStart != "" {
		params["date_start"] = hotNews.Spec.DateStart
	}
	if hotNews.Spec.DateEnd != "" {
		params["date_end"] = hotNews.Spec.DateEnd
	}
	if len(combinedSources) > 0 {
		params["sources"] = strings.Join(combinedSources, ",")
	}
	return buildQuery(params)
}

// Converts a map of query parameters into a URL-encoded query string.
func buildQuery(params map[string]string) string {
	var query []string
	for k, v := range params {
		query = append(query, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(query, "&")
}

// Resolves feed group names to actual source names using the data in the ConfigMap.
func (r *HotNewsReconciler) resolveFeedGroups(feedGroups []string, configMap *corev1.ConfigMap) []string {
	var feedNames []string
	for _, group := range feedGroups {
		if feeds, found := configMap.Data[group]; found {
			feedNames = append(feedNames, strings.Split(feeds, ",")...)
		}
	}
	return feedNames
}

// Fetches articles from the news aggregator service using the constructed URL.
func (r *HotNewsReconciler) fetchArticles(url string) ([]models.Article, error) {
	resp, err := r.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get articles from news aggregator: %s", resp.Status)
	}
	var articles []models.Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}
	return articles, nil
}

// Updates the HotNews status fields with information about the articles retrieved.
func (r *HotNewsReconciler) updateHotNewsStatus(hotNews *aggregatorv1.HotNews, articles []models.Article, reqURL string) {
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = reqURL
	hotNews.Status.ArticlesTitles = getTitles(articles, hotNews.Spec.SummaryConfig.TitlesCount)
}

// Extracts the titles of a specified number of articles.
func getTitles(articles []models.Article, count int) []string {
	titles := make([]string, 0, count)
	for i, article := range articles {
		if i >= count {
			break
		}
		titles = append(titles, article.Title)
	}
	return titles
}

// Updates the Conditions field in the HotNews status with the provided condition.
func (r *HotNewsReconciler) updateStatusCondition(hotNews *aggregatorv1.HotNews, status metav1.ConditionStatus, reason, message string) {
	conditionType := aggregatorv1.HNewsUpdated
	if len(hotNews.Status.Conditions) == 0 {
		conditionType = aggregatorv1.HNewsAdded
	}
	newCondition := aggregatorv1.HotNewsCondition{
		Type:           conditionType,
		Status:         status,
		LastUpdateTime: metav1.Time{Time: time.Now()},
		Reason:         reason,
		Message:        message,
	}
	hotNews.Status.SetCondition(newCondition)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		WithEventFilter(predicates.HotNews()).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(
				handlers.MapConfigMapToHotNews(r.Client, r.ConfigMapName),
			),
		).
		Watches(
			&aggregatorv1.Source{},
			handler.EnqueueRequestsFromMapFunc(
				handlers.MapSourceToHotNews(r.Client, r.fetchConfigMap),
			),
		).
		Complete(r)
}
