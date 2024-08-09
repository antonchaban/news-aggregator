/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

const hotNewsFinalizer = "hotnews.finalizers.teamdev.com"

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
	logger := log.FromContext(ctx)

	// Fetch the HotNews instance
	hotNews := &aggregatorv1.HotNews{}
	err := r.Get(ctx, req.NamespacedName, hotNews)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Info("HotNews resource not found, possibly deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		if !slices.Contains(hotNews.Finalizers, hotNewsFinalizer) {
			hotNews.Finalizers = append(hotNews.Finalizers, hotNewsFinalizer)
			if err := r.Update(ctx, hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if slices.Contains(hotNews.Finalizers, hotNewsFinalizer) {
			// Remove the finalizer and update the resource
			hotNews.Finalizers = slices.Delete(hotNews.Finalizers, slices.Index(hotNews.Finalizers, hotNewsFinalizer), slices.Index(hotNews.Finalizers, hotNewsFinalizer)+1)
			if err := r.Update(ctx, hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	queryParams := make(map[string]string)
	if len(hotNews.Spec.Keywords) > 0 {
		logrus.Println("Keywords: ", hotNews.Spec.Keywords)
		queryParams["keywords"] = strings.Join(hotNews.Spec.Keywords, ",")
	}
	if hotNews.Spec.DateStart != "" {
		queryParams["date_start"] = hotNews.Spec.DateStart
	}
	if hotNews.Spec.DateEnd != "" {
		queryParams["date_end"] = hotNews.Spec.DateEnd
	}
	if len(hotNews.Spec.Sources) > 0 {
		logrus.Println("Sources: ", hotNews.Spec.Sources)
		// todo add validation of short names to webhook
		queryParams["sources"] = strings.Join(hotNews.Spec.Sources, ",")
	}

	reqURL := fmt.Sprintf("%s?%s", r.ArticleSvcURL, buildQuery(queryParams))
	logrus.Println("Request URL: ", reqURL)

	// Fetch news from the aggregator
	articles, err := r.fetchArticles(reqURL)
	if err != nil {
		logger.Error(err, "unable to fetch articles")
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
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
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

// unique removes duplicate strings from a slice
func unique(strings []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strings {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		Complete(r)
}
