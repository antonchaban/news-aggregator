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
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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
	Scheme          *runtime.Scheme
	HTTPClient      *http.Client // HTTP client for making external requests
	ArticleSvcURL   string       // URL of the news aggregator source service
	ConfigMapName   string       // Name of the ConfigMap that contains feed groups
	CfgMapNameSpace string       // Namespace of the ConfigMap
}

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
	// If not above - it is a standard HotNews
	hotNews := &aggregatorv1.HotNews{}
	err := r.Client.Get(ctx, req.NamespacedName, hotNews)
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
	// Fetch the ConfigMap containing feed groups
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: r.CfgMapNameSpace, Name: r.ConfigMapName}, configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Println("ConfigMap not found")
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
	queryParams := r.buildParams(hotNews)
	logrus.Println("Query params: ", queryParams)

	reqURL := fmt.Sprintf("%s?%s", r.ArticleSvcURL, buildQuery(queryParams))
	logrus.Println("Request URL: ", reqURL)

	// Fetch news from the aggregator
	articles, err := r.fetchArticles(reqURL)
	if err != nil {
		err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, metav1.ConditionFalse, "FailedFetchArticles", err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		logrus.Error(err, "unable to fetch articles")
		return ctrl.Result{}, err
	}

	logrus.Println("Articles ttl: ", getTitles(articles, hotNews.Spec.SummaryConfig.TitlesCount))
	logrus.Println("Articles count: ", len(articles))
	logrus.Println("News link: ", reqURL)
	// Update status
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = reqURL
	hotNews.Status.ArticlesTitles = getTitles(articles, hotNews.Spec.SummaryConfig.TitlesCount)

	// Update the HotNews status
	err = r.Client.Status().Update(ctx, hotNews)
	if err != nil {
		err := r.updateHotNewsStatus(ctx, hotNews, aggregatorv1.HNewsUpdated, metav1.ConditionFalse, "FailedUpdateStatus", err.Error())
		if err != nil {
			return ctrl.Result{}, err
		}
		logrus.Error(err, "unable to update HotNews status")
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
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

func (r *HotNewsReconciler) buildParams(hotNews *aggregatorv1.HotNews) map[string]string {
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
	return queryParams
}

// handleConfigMapUpdate handles updates to the ConfigMap and reconciles all relevant HotNews resources
/*func (r *HotNewsReconciler) handleConfigMapUpdate(ctx context.Context, configMap *corev1.ConfigMap) (ctrl.Result, error) {
  var hotNewsList aggregatorv1.HotNewsList
  err := r.Client.List(ctx, &hotNewsList, &client.ListOptions{Namespace: configMap.Namespace})
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
}*/

// handleSourceUpdate handles updates to the Source and reconciles all relevant HotNews resources
/*func (r *HotNewsReconciler) handleSourceUpdate(ctx context.Context, source *aggregatorv1.Source) (ctrl.Result, error) {
  var hotNewsList aggregatorv1.HotNewsList
  err := r.Client.List(ctx, &hotNewsList, &client.ListOptions{Namespace: source.Namespace})
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
    err := r.Client.Get(ctx, client.ObjectKey{Namespace: hotNews.Namespace, Name: r.ConfigMapName}, configMap)
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
*/
func buildQuery(params map[string]string) string {
	var query []string
	for k, v := range params {
		query = append(query, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(query, "&")
}

// resolveFeedGroups resolves feed groups to feed names
func (r *HotNewsReconciler) resolveFeedGroups(feedGroups []string, configMap *corev1.ConfigMap) []string {
	logrus.Println("Resolving feed groups")
	var feedNames []string
	for _, group := range feedGroups {
		if feeds, found := configMap.Data[group]; found {
			feedNames = append(feedNames, strings.Split(feeds, ",")...)
		}
	}
	logrus.Println("Resolved feed names: ", feedNames)
	return feedNames
}

// fetchArticles fetches articles from the given URL
func (r *HotNewsReconciler) fetchArticles(url string) ([]Article, error) {
	logrus.Println("Fetching articles from: ", url)
	resp, err := r.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logrus.Println("Response status: ", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get articles from news aggregator: %s", resp.Status)
	}

	var articles []Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		logrus.Println("Error decoding response: ", err)
		return nil, err
	}

	logrus.Println("Articles fetched in fetchArticles(): ", articles)
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

// todo validate configmaps and sources
// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
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
