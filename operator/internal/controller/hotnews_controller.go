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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"slices"
	"strings"
	"time"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	newsAggregatorServiceURL = "https://news-alligator-service.news-alligator.svc.cluster.local:8443/articles"
	hotNewsFinalizer         = "hotnews.finalizers.teamdev.com"
)

type Article struct {
	Id          int       `json:"Id"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Link        string    `json:"Link"`
	Source      Source    `json:"Source"`
	PubDate     time.Time `json:"PubDate"`
}

type Source struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=hotnews/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HotNews object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var hotNews aggregatorv1.HotNews
	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
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
			if err := r.Client.Update(ctx, &hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if slices.Contains(hotNews.Finalizers, hotNewsFinalizer) {
			// Simply remove the finalizer and update the resource
			hotNews.Finalizers = slices.Delete(hotNews.Finalizers,
				slices.Index(hotNews.Finalizers, hotNewsFinalizer), slices.Index(hotNews.Finalizers, hotNewsFinalizer)+1)
			if err := r.Client.Update(ctx, &hotNews); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	logrus.Info("Reconciling HotNews", "Name", hotNews.Name, "Keywords", hotNews.Spec.Keywords)

	// Construct query parameters
	queryParams := make(map[string]string)
	if len(hotNews.Spec.Keywords) > 0 {
		queryParams["keywords"] = strings.Join(hotNews.Spec.Keywords, ",")
	}
	if hotNews.Spec.StartDate != "" {
		queryParams["date_start"] = hotNews.Spec.StartDate
	}
	if hotNews.Spec.EndDate != "" {
		queryParams["date_end"] = hotNews.Spec.EndDate
	}
	if len(hotNews.Spec.Sources) > 0 {
		var validSources []string
		for _, source := range hotNews.Spec.Sources {
			if isValidSource(source) {
				validSources = append(validSources, source)
			}
		}
		queryParams["sources"] = strings.Join(validSources, ",")
	}

	// Generate the request URL
	reqURL := fmt.Sprintf("%s?%s", newsAggregatorServiceURL, buildQuery(queryParams))
	logrus.Println("Request URL: ", reqURL)

	// Send the HTTP request
	c := &http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	resp, err := c.Get(reqURL)
	if err != nil {
		logrus.Error("Failed to fetch articles: ", err)
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error(fmt.Errorf("failed to get articles from news aggregator"), "Status", resp.Status)
		return ctrl.Result{}, fmt.Errorf("failed to get articles from news aggregator: %s", resp.Status)
	}

	var articles []Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		logrus.Error("Failed to decode response: ", err)
		return ctrl.Result{}, err
	}

	hotNews.Status.ArticlesTitles = make([]string, 0, len(articles))
	// Update HotNews status
	hotNews.Status.ArticlesCount = len(articles)
	hotNews.Status.NewsLink = reqURL
	for _, article := range articles {
		hotNews.Status.ArticlesTitles = append(hotNews.Status.ArticlesTitles, article.Title)
	}

	if err := r.Status().Update(ctx, &hotNews); err != nil {
		logrus.Error("Failed to update HotNews status: ", err)
		return ctrl.Result{}, err
	}

	logrus.Info("Successfully reconciled HotNews", "Name", hotNews.Name)
	return ctrl.Result{}, nil
}

func buildQuery(params map[string]string) string {
	var query []string
	for k, v := range params {
		query = append(query, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(query, "&")
}

func isValidSource(source string) bool {
	validSources := []string{"abcnews", "bbc", "washingtontimes", "nbc", "usatoday", "other"}
	for _, s := range validSources {
		if source == s {
			return true
		}
	}
	return false
}

func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.HotNews{}).
		//Watches(&corev1.ConfigMap{}, handler.TypedEventHandler())).
		Complete(r)
}
