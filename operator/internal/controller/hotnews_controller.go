package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller/models"
	"com.teamdev/news-aggregator/internal/controller/predicates"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"strings"
	"time"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	Client             client.Client   // Client for interacting with the Kubernetes API.
	Scheme             *runtime.Scheme // Scheme for the reconciler.
	HTTPClient         *http.Client    // HTTP client for making external requests
	ArticleSvcURL      string          // URL of the news aggregator source service
	ConfigMapName      string          // Name of the ConfigMap that contains feed groups
	ConfigMapNamespace string          // Namespace of the ConfigMap
	WorkingNamespace   string          // Namespace where CRDs are created
}

const HotNewsFinalizer = "hotnews.finalizers.teamdev.com"

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Modify the Reconcile function to compare the state specified by
// the HotNews object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
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

func (r *HotNewsReconciler) handleFinalizer(ctx context.Context, hotNews *aggregatorv1.HotNews) (bool, error) {
	if !hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is being deleted
		if slices.Contains(hotNews.Finalizers, HotNewsFinalizer) {
			if err := r.deleteOwnerReferences(ctx, hotNews.Namespace, hotNews.Name); err != nil {
				return false, err
			}
			hotNews.Finalizers = slices.Delete(hotNews.Finalizers,
				slices.Index(hotNews.Finalizers, HotNewsFinalizer),
				slices.Index(hotNews.Finalizers, HotNewsFinalizer)+1)
			if err := r.Client.Update(ctx, hotNews); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	// Add finalizer if not present
	if !slices.Contains(hotNews.Finalizers, HotNewsFinalizer) {
		hotNews.Finalizers = append(hotNews.Finalizers, HotNewsFinalizer)
		if err := r.Client.Update(ctx, hotNews); err != nil {
			return false, err
		}
	}
	return false, nil
}

func (r *HotNewsReconciler) deleteOwnerReferences(ctx context.Context, namespace, hotNewsName string) error {
	var sources aggregatorv1.SourceList
	if err := r.Client.List(ctx, &sources, &client.ListOptions{Namespace: namespace}); err != nil {
		logrus.Errorf("Failed to list Sources: %v", err)
		return fmt.Errorf("failed to list Sources: %w", err)
	}

	for _, src := range sources.Items {
		updatedOwnerRefs := removeOwnerReference(src.OwnerReferences, hotNewsName, "HotNews")
		if len(updatedOwnerRefs) != len(src.OwnerReferences) {
			src.OwnerReferences = updatedOwnerRefs
			if err := r.Client.Update(ctx, &src); err != nil {
				logrus.Errorf("Failed to update Source %s: %v", src.Name, err)
				return fmt.Errorf("failed to update Source %s: %w", src.Name, err)
			}
			logrus.Infof("Owner reference removed from Source %s", src.Name)
		}
	}
	return nil
}

// removeOwnerReference removes a specific owner reference from the list.
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

// reconcileHotNews performs the actual reconciliation logic for HotNews
func (r *HotNewsReconciler) reconcileHotNews(ctx context.Context, hotNews *aggregatorv1.HotNews) (ctrl.Result, error) {
	// Fetch the ConfigMap containing feed groups
	configMap, err := r.fetchConfigMap(ctx)
	if err != nil {
		logrus.Error(err, "Failed to fetch ConfigMap")
		r.updateStatusCondition(ctx, hotNews, metav1.ConditionFalse, "ConfigMapNotFound", "ConfigMap not found")
		return ctrl.Result{}, err
	}

	// Resolve FeedGroups to actual source names
	if len(hotNews.Spec.FeedGroups) > 0 {
		resolvedSources := r.resolveFeedGroups(hotNews.Spec.FeedGroups, configMap)
		hotNews.Spec.Sources = append(hotNews.Spec.Sources, resolvedSources...)
	}

	// Set OwnerReferences for each source using ShortName
	if err := r.setOwnerReferences(ctx, hotNews, hotNews.Spec.Sources); err != nil {
		return ctrl.Result{}, err
	}

	// Build query parameters and fetch articles
	reqURL := fmt.Sprintf("%s?%s", r.ArticleSvcURL, buildQueryParams(hotNews))
	logrus.Println("Request URL:", reqURL)

	articles, err := r.fetchArticles(reqURL)
	if err != nil {
		logrus.Error(err, "Unable to fetch articles")
		r.updateStatusCondition(ctx, hotNews, metav1.ConditionFalse, "FetchArticlesFailed", "Failed to fetch articles")
		return ctrl.Result{}, err
	}

	// Update HotNews status
	r.updateHotNewsStatus(hotNews, articles, reqURL)
	r.updateStatusCondition(ctx, hotNews, metav1.ConditionTrue, "Reconciled", "Successfully reconciled HotNews")
	logrus.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
	return ctrl.Result{}, nil
}

// fetchConfigMap retrieves the ConfigMap containing feed groups.
func (r *HotNewsReconciler) fetchConfigMap(ctx context.Context) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: r.ConfigMapNamespace, Name: r.ConfigMapName}, configMap)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// setOwnerReferences sets OwnerReferences for Sources based on ShortNames.
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
				logrus.Error(err, "Failed to update Source with owner reference", "ShortName", sourceShortName)
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
func buildQueryParams(hotNews *aggregatorv1.HotNews) string {
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
	if len(hotNews.Spec.Sources) > 0 {
		params["sources"] = strings.Join(hotNews.Spec.Sources, ",")
	}
	return buildQuery(params)
}

// buildQuery builds a query string from a map of parameters.
func buildQuery(params map[string]string) string {
	var query []string
	for k, v := range params {
		query = append(query, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(query, "&")
}

// resolveFeedGroups resolves feed groups to feed names using the ConfigMap.
func (r *HotNewsReconciler) resolveFeedGroups(feedGroups []string, configMap *corev1.ConfigMap) []string {
	var feedNames []string
	for _, group := range feedGroups {
		if feeds, found := configMap.Data[group]; found {
			feedNames = append(feedNames, strings.Split(feeds, ",")...)
		}
	}
	return feedNames
}

// fetchArticles fetches articles from the given URL.
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

// updateHotNewsStatus updates the HotNews status with articles information.
func (r *HotNewsReconciler) updateHotNewsStatus(hotNews *aggregatorv1.HotNews, articles []models.Article, reqURL string) {
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = reqURL
	hotNews.Status.ArticlesTitles = getTitles(articles, hotNews.Spec.SummaryConfig.TitlesCount)
}

// getTitles extracts the titles of the articles.
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

// updateStatusCondition updates the HotNews status condition.
func (r *HotNewsReconciler) updateStatusCondition(ctx context.Context, hotNews *aggregatorv1.HotNews, status metav1.ConditionStatus, reason, message string) {
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
	if err := r.Client.Status().Update(ctx, hotNews); err != nil {
		logrus.Error(err, "Unable to update HotNews status")
		r.updateStatusCondition(ctx, hotNews, metav1.ConditionFalse, "UpdateStatusFailed", "Failed to update HotNews status")
	}
}

// MapConfigMapToHotNews maps ConfigMap updates to HotNews resources.
func (r *HotNewsReconciler) MapConfigMapToHotNews(ctx context.Context, object client.Object) []reconcile.Request {
	logrus.Println("Mapping ConfigMap to HotNews")
	var hotNewsList aggregatorv1.HotNewsList
	if err := r.Client.List(ctx, &hotNewsList, &client.ListOptions{Namespace: r.ConfigMapNamespace}); err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return nil
	}

	var requests []reconcile.Request
	configMap := object.(*corev1.ConfigMap)

	for _, hotNews := range hotNewsList.Items {
		if r.shouldReconcileConfigMap(hotNews, configMap) {
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

// shouldReconcileConfigMap determines if a HotNews resource should be reconciled based on a ConfigMap update.
func (r *HotNewsReconciler) shouldReconcileConfigMap(hotNews aggregatorv1.HotNews, configMap *corev1.ConfigMap) bool {
	if r.ConfigMapName == configMap.Name {
		return true
	}
	for _, feedGroup := range hotNews.Spec.FeedGroups {
		if _, found := configMap.Data[feedGroup]; found {
			return true
		}
	}
	return false
}

// MapSourceToHotNews maps Source updates to HotNews resources.
func (r *HotNewsReconciler) MapSourceToHotNews(ctx context.Context, object client.Object) []reconcile.Request {
	logrus.Println("Mapping Source to HotNews")
	source := object.(*aggregatorv1.Source)

	var hotNewsList aggregatorv1.HotNewsList
	if err := r.Client.List(ctx, &hotNewsList, &client.ListOptions{Namespace: source.Namespace}); err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return nil
	}

	var requests []reconcile.Request
	for _, hotNews := range hotNewsList.Items {
		if slices.Contains(hotNews.Spec.Sources, source.Spec.ShortName) {
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

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		WithEventFilter(predicates.HotNews(r.WorkingNamespace)).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.MapConfigMapToHotNews),
		).
		Watches(
			&aggregatorv1.Source{},
			handler.EnqueueRequestsFromMapFunc(r.MapSourceToHotNews),
		).
		Complete(r)
}
