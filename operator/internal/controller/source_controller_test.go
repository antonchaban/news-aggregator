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
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aggregatorv1 "com.teamdev/news-aggregator/api/v1"
)

var _ = Describe("Source Controller", func() {
	const (
		resourceName = "test-resource"
		timeout      = time.Second * 30
		interval     = time.Millisecond * 500
	)

	var (
		ctx                  context.Context
		typeNamespacedName   types.NamespacedName
		source               *aggregatorv1.Source
		controllerReconciler *SourceReconciler
		server               *httptest.Server
	)

	cleanupResource := func(typeNamespacedName types.NamespacedName) {
		resource := &aggregatorv1.Source{}
		err := k8sClient.Get(ctx, typeNamespacedName, resource)
		if err == nil {
			// remove finalizers
			if len(resource.Finalizers) > 0 {
				resource.Finalizers = []string{}
				err = k8sClient.Update(ctx, resource)
				Expect(err).NotTo(HaveOccurred())
			}

			// Resource exists, delete it
			err = k8sClient.Delete(ctx, resource)
			Expect(err).NotTo(HaveOccurred())

			// w8 for the resource to be actually deleted
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, resource)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		} else {
			Expect(errors.IsNotFound(err)).To(BeTrue())
		}
	}

	BeforeEach(func() {
		ctx = context.Background()
		typeNamespacedName = types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		source = &aggregatorv1.Source{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: "default",
			},
			Spec: aggregatorv1.SourceSpec{
				Name: "Test Source",
				Link: "http://example.com/rss",
			},
			Status: aggregatorv1.SourceStatus{
				ID: 0,
			},
		}

		cleanupResource(typeNamespacedName)

		// Create fake HTTP server to mock an external service
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var spec aggregatorv1.SourceSpec
				err := json.NewDecoder(r.Body).Decode(&spec)
				Expect(err).NotTo(HaveOccurred())

				spec.Id = 1
				respData, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())

				w.WriteHeader(http.StatusOK)
				w.Write(respData)
			} else if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))

		controllerReconciler = &SourceReconciler{
			Client:                      k8sClient,
			Scheme:                      k8sClient.Scheme(),
			HTTPClient:                  server.Client(),
			NewsAggregatorSrcServiceURL: server.URL,
		}

		err := k8sClient.Create(ctx, source)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		cleanupResource(typeNamespacedName)
		server.Close()
	})

	Context("When reconciling a resource", func() {
		It("should successfully create the resource in the external service", func() {
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			updatedSource := &aggregatorv1.Source{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedSource)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedSource.Status.ID).To(Equal(1))

			conditions := updatedSource.Status.Conditions
			Expect(len(conditions)).To(Equal(1))
			Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceAdded))
			Expect(conditions[0].Status).To(Equal(metav1.ConditionTrue))
			Expect(conditions[0].Reason).To(Equal(reasonSuccessCr))
		})

		It("should successfully update the resource in the external service", func() {
			source.Status.ID = 1
			err := k8sClient.Status().Update(ctx, source)
			Expect(err).NotTo(HaveOccurred())

			source.Spec.Name = "Updated Source"
			err = k8sClient.Update(ctx, source)
			Expect(err).NotTo(HaveOccurred())

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			updatedSource := &aggregatorv1.Source{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedSource)
			Expect(err).NotTo(HaveOccurred())

			conditions := updatedSource.Status.Conditions
			Expect(len(conditions)).To(Equal(1))
			Expect(conditions[0].Type).To(Equal(aggregatorv1.SourceUpdated))
			Expect(conditions[0].Status).To(Equal(metav1.ConditionTrue))
			Expect(conditions[0].Reason).To(Equal(reasonSuccessUpd))
		})

		It("should successfully delete the resource from the external service", func() {
			source.Status.ID = 1
			source.Finalizers = append(source.Finalizers, srcFinalizer)
			err := k8sClient.Update(ctx, source)
			Expect(err).NotTo(HaveOccurred())

			// Initiate the deletion
			err = k8sClient.Delete(ctx, source)
			Expect(err).NotTo(HaveOccurred())

			// Reconcile to handle finalizers
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the source has been removed from the external service
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, source)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			deletedSource := &aggregatorv1.Source{}
			err = k8sClient.Get(ctx, typeNamespacedName, deletedSource)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should add the finalizer if not present", func() {
			source.Finalizers = []string{}
			err := k8sClient.Update(ctx, source)
			Expect(err).NotTo(HaveOccurred())

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			updatedSource := &aggregatorv1.Source{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedSource)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedSource.Finalizers).To(ContainElement(srcFinalizer))
		})
	})
})
