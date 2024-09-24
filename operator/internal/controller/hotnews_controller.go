package controller

import (
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"com.teamdev/news-aggregator/internal/controller/predicates"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"slices"
	"strings"
	"time"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	Client             client.Client
	Scheme             *runtime.Scheme
	HTTPClient         *http.Client // HTTP client for making external requests
	ArticleSvcURL      string       // URL of the news aggregator source service
	ConfigMapName      string       // Name of the ConfigMap that contains feed groups
	ConfigMapNamespace string       // Namespace of the ConfigMap
}

const HotNewsFinalizer = "hotnews.finalizers.teamdev.com"

type Article struct {
	Id          int       `json:"Id"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Link        string    `json:"Link"`
	Source      Source    `json:"Source"`
	PubDate     time.Time `json:"PubDate"`
}
type Source struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	ShortName string `json:"short_name"`
}

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
	hotNews := &aggregatorv1.HotNews{}
	err := r.Client.Get(ctx, req.NamespacedName, hotNews)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.FromContext(ctx).Error(err, "Failed to get HotNews")
			return ctrl.Result{}, err
		}
		logrus.Infof("HotNews resource not found, possibly deleted. In namespace: %s", req.Namespace)
		return ctrl.Result{}, nil
	}

	// Handle finalizer logic
	if hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		if !slices.Contains(hotNews.Finalizers, HotNewsFinalizer) {
			hotNews.Finalizers = append(hotNews.Finalizers, HotNewsFinalizer)
			if err := r.Client.Update(ctx, hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if slices.Contains(hotNews.Finalizers, HotNewsFinalizer) {
			if err := r.deleteOwnerRef(ctx, hotNews.Namespace, hotNews.Name); err != nil {
				return ctrl.Result{}, err
			}
			hotNews.Finalizers = slices.Delete(hotNews.Finalizers, slices.Index(hotNews.Finalizers, HotNewsFinalizer), slices.Index(hotNews.Finalizers, HotNewsFinalizer)+1)
			if err := r.Client.Update(ctx, hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Proceed with normal reconciliation for HotNews
	return r.reconcileHotNews(ctx, hotNews)
}

func (r *HotNewsReconciler) deleteOwnerRef(ctx context.Context, namespace, hotNewsName string) error {
	var sources aggregatorv1.SourceList

	if err := r.Client.List(ctx, &sources, &client.ListOptions{Namespace: namespace}); err != nil {
		logrus.Errorf("Failed to list Sources: %v", err)
		return fmt.Errorf("failed to list Sources: %w", err)
	}

	for _, src := range sources.Items {
		var newOwnerReferences []metav1.OwnerReference
		isRemoved := false

		for _, ref := range src.OwnerReferences {
			if ref.Name == hotNewsName && ref.Kind == "HotNews" {
				isRemoved = true
				continue
			}
			newOwnerReferences = append(newOwnerReferences, ref)
		}

		if isRemoved {
			src.OwnerReferences = newOwnerReferences
			if err := r.Client.Update(ctx, &src); err != nil {
				logrus.Errorf("Failed to update Source %s: %v", src.Name, err)
				return fmt.Errorf("failed to update Source %s: %w", src.Name, err)
			}
			logrus.Infof("Owner ref removed from src %s", src.Name)
		}
	}

	return nil
}

// reconcileHotNews performs the actual reconciliation logic for HotNews
func (r *HotNewsReconciler) reconcileHotNews(ctx context.Context, hotNews *aggregatorv1.HotNews) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ConfigMap containing feed groups
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: r.ConfigMapNamespace, Name: r.ConfigMapName}, configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "ConfigMap not found")
			errUp := r.selectStatus(ctx, hotNews, metav1.ConditionFalse, "ConfigMapNotFound", "ConfigMap not found")
			if errUp != nil {
				return ctrl.Result{}, errUp
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	// Resolve FeedGroups to actual source names
	if len(hotNews.Spec.FeedGroups) > 0 {
		resolvedSources := r.resolveFeedGroups(hotNews.Spec.FeedGroups, configMap)
		hotNews.Spec.Sources = append(hotNews.Spec.Sources, resolvedSources...)
	}

	// Set OwnerReferences for each source using ShortName
	for _, sourceShortName := range hotNews.Spec.Sources {
		var sourceList aggregatorv1.SourceList
		err := r.Client.List(ctx, &sourceList, &client.ListOptions{Namespace: hotNews.Namespace})
		if err != nil {
			logger.Error(err, "Failed to list Sources")
			continue
		}

		var source *aggregatorv1.Source
		for _, src := range sourceList.Items {
			if src.Spec.ShortName == sourceShortName {
				source = &src
				break
			}
		}

		if source == nil {
			logger.Error(fmt.Errorf("source with short_name %s not found", sourceShortName), "Failed to find Source")
			continue
		}

		// Check if the owner reference already exists
		ownerExists := false
		for _, owner := range source.OwnerReferences {
			if owner.UID == hotNews.UID {
				ownerExists = true
				break
			}
		}

		// Add owner reference if not exists
		if !ownerExists {
			source.OwnerReferences = append(source.OwnerReferences, metav1.OwnerReference{
				APIVersion: hotNews.APIVersion,
				Kind:       hotNews.Kind,
				Name:       hotNews.Name,
				UID:        hotNews.UID,
			})
			if err := r.Client.Update(ctx, source); err != nil {
				logger.Error(err, "Failed to update Source with owner reference", "ShortName", sourceShortName)
				return ctrl.Result{}, err
			}
		}
	}

	// Build query parameters for the HTTP request
	queryParams := make(map[string]string)
	if len(hotNews.Spec.Keywords) > 0 {
		queryParams["keywords"] = strings.Join(hotNews.Spec.Keywords, ",")
	}
	if hotNews.Spec.DateStart != "" {
		queryParams["date_start"] = hotNews.Spec.DateStart
	}
	if hotNews.Spec.DateEnd != "" {
		queryParams["date_end"] = hotNews.Spec.DateEnd
	}
	if len(hotNews.Spec.Sources) > 0 {
		queryParams["sources"] = strings.Join(hotNews.Spec.Sources, ",")
	}
	reqURL := fmt.Sprintf("%s?%s", r.ArticleSvcURL, buildQuery(queryParams))
	logrus.Println("Request URL: ", reqURL)
	// Fetch news from the aggregator
	articles, err := r.fetchArticles(reqURL)
	if err != nil {
		logger.Error(err, "unable to fetch articles")
		errUp := r.selectStatus(ctx, hotNews, metav1.ConditionFalse, "FetchArticlesFailed", "Failed to fetch articles")
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}

	// Update status
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = reqURL
	hotNews.Status.ArticlesTitles = getTitles(articles, hotNews.Spec.SummaryConfig.TitlesCount)

	// Update the HotNews status
	err = r.Client.Status().Update(ctx, hotNews)
	if err != nil {
		logger.Error(err, "unable to update HotNews status")
		errUp := r.selectStatus(ctx, hotNews, metav1.ConditionFalse, "UpdateStatusFailed", "Failed to update HotNews status")
		if errUp != nil {
			return ctrl.Result{}, errUp
		}
		return ctrl.Result{}, err
	}

	errUp := r.selectStatus(ctx, hotNews, metav1.ConditionTrue, "Reconciled", "Successfully reconciled HotNews")
	if errUp != nil {
		return ctrl.Result{}, errUp
	}
	logger.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
	return ctrl.Result{}, nil
}

// MapConfigMapToHotNews implement the event handler to map ConfigMap updates to HotNews resources
func (r *HotNewsReconciler) MapConfigMapToHotNews(ctx context.Context, object client.Object) []reconcile.Request {
	logrus.Println("Mapping ConfigMap to HotNews")
	var hotNewsList aggregatorv1.HotNewsList
	err := r.Client.List(ctx, &hotNewsList, &client.ListOptions{Namespace: r.ConfigMapNamespace})
	if err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return nil
	}

	var requests []reconcile.Request
	configMap := object.(*corev1.ConfigMap)

	// Iterate over each HotNews resource in the same namespace
	for _, hotNews := range hotNewsList.Items {
		shouldReconcile := false

		// Check if the ConfigMap name matches the one specified in the HotNews resource
		if r.ConfigMapName == configMap.Name {
			// If the ConfigMap directly referenced by HotNews was updated
			resolvedSources := r.resolveFeedGroups(hotNews.Spec.FeedGroups, configMap)
			logrus.Println("Resolved sources: ", resolvedSources)
			if len(resolvedSources) > 0 {
				hotNews.Spec.Sources = append(hotNews.Spec.Sources, resolvedSources...)
				shouldReconcile = true
			}
		} else {
			// Check if any feed group in HotNews references a source in this ConfigMap
			for _, feedGroup := range hotNews.Spec.FeedGroups {
				if feeds, found := configMap.Data[feedGroup]; found {
					feedSources := strings.Split(feeds, ",")
					if len(feedSources) > 0 {
						shouldReconcile = true
					}
				}
			}
		}

		if shouldReconcile {
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

// MapSourceToHotNews implement the event handler to map Source updates to HotNews resources
func (r *HotNewsReconciler) MapSourceToHotNews(ctx context.Context, object client.Object) []reconcile.Request {
	// Fetch the Source object
	logrus.Println("Mapping Source to HotNews")
	source := &aggregatorv1.Source{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: object.GetNamespace(), Name: object.GetName()}, source)
	if err != nil {
		logrus.Errorf("Failed to fetch Source %s: %v", object.GetName(), err)
		return nil
	}

	// Now list all HotNews objects in the same namespace
	var hotNewsList aggregatorv1.HotNewsList
	err = r.Client.List(ctx, &hotNewsList, &client.ListOptions{Namespace: object.GetNamespace()})
	if err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return nil
	}

	var requests []reconcile.Request
	for _, hotNews := range hotNewsList.Items {
		// Check if the Source's ShortName is in the HotNews Spec.Sources list
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

func buildQuery(params map[string]string) string {
	var query []string
	for k, v := range params {
		query = append(query, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(query, "&")
}

// resolveFeedGroups resolves feed groups to feed names
func (r *HotNewsReconciler) resolveFeedGroups(feedGroups []string, configMap *corev1.ConfigMap) []string {
	var feedNames []string
	for _, group := range feedGroups {
		if feeds, found := configMap.Data[group]; found {
			feedNames = append(feedNames, strings.Split(feeds, ",")...)
		}
	}
	return feedNames
}

// fetchArticles fetches articles from the given URL
func (r *HotNewsReconciler) fetchArticles(url string) ([]Article, error) {
	resp, err := r.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get articles from news aggregator: %s", resp.Status)
	}
	var articles []Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}
	return articles, nil
}

// getTitles extracts the titles of the articles
func getTitles(articles []Article, count int) []string {
	var titles []string
	for i, article := range articles {
		if i >= count {
			break
		}
		titles = append(titles, article.Title)
	}
	return titles
}

func (r *HotNewsReconciler) selectStatus(ctx context.Context, hotNews *aggregatorv1.HotNews, cndStatus metav1.ConditionStatus, reason, message string) error {
	if len(hotNews.Status.Conditions) == 0 {
		err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsAdded, cndStatus, reason, message)
		if err != nil {
			return err
		}
	} else {
		err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, cndStatus, reason, message)
		if err != nil {
			return err
		}
	}
	return nil
}

// updateHotNewsStatus updates the SourceStatus of a source resource with the given condition.
func (r *HotNewsReconciler) updateHotNewsStatus(ctx context.Context, hotNews *aggregatorv1.HotNews, conditionType aggregatorv1.HNewsConditionType, status metav1.ConditionStatus, reason, message string) error {
	logrus.Println("Updating HotNews status")
	logrus.Println("HotNews status: ", hotNews.Status)
	newCondition := aggregatorv1.HotNewsCondition{
		Type:           conditionType,
		Status:         status,
		LastUpdateTime: metav1.Time{Time: time.Now()},
		Reason:         reason,
		Message:        message,
	}

	hotNews.Status.Conditions = append(hotNews.Status.Conditions, newCondition)
	logrus.Println("Appending new condition")
	logrus.Println("HotNews status: ", hotNews.Status)
	return r.Client.Status().Update(ctx, hotNews)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		WithEventFilter(predicates.HotNews()).
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
