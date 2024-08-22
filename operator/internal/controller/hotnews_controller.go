package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"slices"
	"strings"
	"time"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	HTTPClient    *http.Client // HTTP client for making external requests
	ArticleSvcURL string       // URL of the news aggregator source service
	ConfigMapName string       // Name of the ConfigMap that contains feed groups
}

type Article struct {
	Id          int       `json:"Iіd"`
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
	// Check if the reconcile was triggered by a ConfigMap or Source update
	var configMap corev1.ConfigMap
	var source aggregatorv1.Source
	err := k8sClient.Get(ctx, req.NamespacedName, &configMap)
	if err == nil {
		logrus.Println("Handling ConfigMap update")
		// ConfigMap update detected, handle all HotNews that reference this ConfigMap
		return r.handleConfigMapUpdate(ctx, &configMap)
	}

	err = k8sClient.Get(ctx, req.NamespacedName, &source)
	if err == nil {
		logrus.Println("Handling Source update")
		// Source update detected, handle all HotNews that reference this Source
		return r.handleSourceUpdate(ctx, &source)
	}

	// If not above - it is a standard HotNews
	hotNews := &aggregatorv1.HotNews{}
	err = k8sClient.Get(ctx, req.NamespacedName, hotNews)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info("HotNews resource not found, possibly deleted. In namespace: ", req.Namespace)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Proceed with normal reconciliation for HotNews
	return r.reconcileHotNews(ctx, hotNews)
}

// reconcileHotNews performs the actual reconciliation logic for HotNews
func (r *HotNewsReconciler) reconcileHotNews(ctx context.Context, hotNews *aggregatorv1.HotNews) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ConfigMap containing feed groups
	configMap := &corev1.ConfigMap{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: hotNews.Namespace, Name: r.ConfigMapName}, configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "ConfigMap not found")
			err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, metav1.ConditionFalse, "FailedConfigMap", err.Error())
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Resolve FeedGroups to actual source names
	if len(hotNews.Spec.FeedGroups) > 0 {
		resolvedSources := r.resolveFeedGroups(hotNews.Spec.FeedGroups, configMap)
		hotNews.Spec.Sources = append(hotNews.Spec.Sources, resolvedSources...)
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
		err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, metav1.ConditionFalse, "FailedFetchArticles", err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		logger.Error(err, "unable to fetch articles")
		return ctrl.Result{}, err
	}

	// Update status
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = reqURL
	hotNews.Status.ArticlesTitles = getTitles(articles, hotNews.Spec.SummaryConfig.TitlesCount)

	// Update the HotNews status
	err = k8sClient.Status().Update(ctx, hotNews)
	if err != nil {
		err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, metav1.ConditionFalse, "FailedUpdateStatus", err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		logger.Error(err, "unable to update HotNews status")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
	if len(hotNews.Status.Conditions) == 0 {
		err = r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsAdded, metav1.ConditionTrue, "Success", "HotNews created successfully")
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		err = r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, metav1.ConditionTrue, "Success", "HotNews updated successfully")
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// handleConfigMapUpdate handles updates to the ConfigMap and reconciles all relevant HotNews resources
func (r *HotNewsReconciler) handleConfigMapUpdate(ctx context.Context, configMap *corev1.ConfigMap) (ctrl.Result, error) {
	var hotNewsList aggregatorv1.HotNewsList
	err := k8sClient.List(ctx, &hotNewsList, &client.ListOptions{Namespace: configMap.Namespace})
	if err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return ctrl.Result{}, err
	}

	for _, hotNews := range hotNewsList.Items {
		if r.ConfigMapName == configMap.Name {
			_, err := r.reconcileHotNews(ctx, &hotNews)
			if err != nil {
				logrus.Errorf("Failed to reconcile HotNews: %v", err)
			}
		}
	}
	return ctrl.Result{}, nil
}

// handleSourceUpdate handles updates to the Source and reconciles all relevant HotNews resources
func (r *HotNewsReconciler) handleSourceUpdate(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
	var hotNewsList aggregatorv1.HotNewsList
	err := k8sClient.List(ctx, &hotNewsList, &client.ListOptions{Namespace: source.Namespace})
	if err != nil {
		logrus.Errorf("Failed to list HotNews resources: %v", err)
		return ctrl.Result{}, err
	}

	for _, hotNews := range hotNewsList.Items {
		logrus.Println("HotNews: ", hotNews.Name)
		logrus.Println("Sources: ", hotNews.Spec.Sources)
		logrus.Println("Source: ", source.Spec)

		// First check if the Source is directly in HotNews.Spec.Sources
		if slices.Contains(hotNews.Spec.Sources, source.Spec.ShortName) {
			_, err := r.reconcileHotNews(ctx, &hotNews)
			if err != nil {
				logrus.Errorf("Failed to reconcile HotNews: %v", err)
			}
			continue
		}

		// Fetch the ConfigMap referenced by HotNews
		configMap := &corev1.ConfigMap{}
		err := k8sClient.Get(ctx, client.ObjectKey{Namespace: hotNews.Namespace, Name: r.ConfigMapName}, configMap)
		if err != nil {
			logrus.Errorf("Failed to get ConfigMap %s: %v", r.ConfigMapName, err)
			continue
		}

		// Check if the Source is part of any feed groups in the ConfigMap
		for _, feedGroup := range hotNews.Spec.FeedGroups {
			if feeds, found := configMap.Data[feedGroup]; found {
				feedSources := strings.Split(feeds, ",")
				if slices.Contains(feedSources, source.Spec.ShortName) {
					_, err := r.reconcileHotNews(ctx, &hotNews)
					if err != nil {
						logrus.Errorf("Failed to reconcile HotNews: %v", err)
					}
					break
				}
			}
		}
	}
	return ctrl.Result{}, nil
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

// updateHotNewsStatus updates the SourceStatus of a source resource with the given condition.
func (r *HotNewsReconciler) updateHotNewsStatus(ctx context.Context, hotNews *aggregatorv1.HotNews, conditionType aggregatorv1.HNewsConditionType, status metav1.ConditionStatus, reason, message string) error {
	newCondition := aggregatorv1.HotNewsCondition{
		Type:           conditionType,
		Status:         status,
		LastUpdateTime: metav1.Time{Time: time.Now()},
		Reason:         reason,
		Message:        message,
	}

	hotNews.Status.Conditions = append(hotNews.Status.Conditions, newCondition)
	return k8sClient.Status().Update(ctx, hotNews)
}

// todo validate configmaps and sources
// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient() // Set the global k8sClient variable
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return !e.DeleteStateUnknown
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
			},
		}).
		Watches(
			&corev1.ConfigMap{},
			&handler.EnqueueRequestForObject{},
		).
		Watches(
			&aggregatorv1.Source{},
			&handler.EnqueueRequestForObject{},
		).
		Complete(r)
}
