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
	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// SourceReconciler reconciles a Source object
type SourceReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

const (
	newsAggregatorSrcServiceURL = "https://news-alligator-service.news-alligator.svc.cluster.local:8443/sources"
	newsAggregatorArtServiceURL = "https://news-alligator-service.news-alligator.svc.cluster.local:8443/articles"
)

// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aggregator.com.teamdev,resources=sources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Source object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *SourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var source aggregatorv1.Source
	err := r.Client.Get(ctx, req.NamespacedName, &source)
	if err != nil {
		return ctrl.Result{}, err
	}

	logrus.Printf("name is: %s", source.Spec.Name)
	logrus.Printf("link is: %s", source.Spec.Link)

	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	testReq, err := http.NewRequest(http.MethodGet, newsAggregatorArtServiceURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ctrl.Result{}, nil
	}
	testReq.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.Do(testReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ctrl.Result{}, nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ctrl.Result{}, nil
	}

	logrus.Printf("body is: %s", body)

	logrus.Println("############################")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aggregatorv1.Source{}).
		Complete(r)
}
