/*
Licensed under the Apache License, Version 2.0 (the "License");
*/

package v1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Source Webhook", func() {

	Context("When creating Source under Validating Webhook", func() {
		It("Should deny if the name, short_name, or link field is empty", func() {
			ctx := context.Background()

			// Testing when the name field is empty
			source := &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-invalid-name",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "",
					ShortName: "short",
					Link:      "http://example.com/rss",
				},
			}
			err := k8sClient.Create(ctx, source)
			Expect(err).Should(HaveOccurred())

			// Testing when the short_name field is empty
			source = &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-invalid-shortname",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "Name",
					ShortName: "",
					Link:      "http://example.com/rss",
				},
			}
			err = k8sClient.Create(ctx, source)
			Expect(err).Should(HaveOccurred())

			// Testing when the link field is empty
			source = &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-invalid-link",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "Name",
					ShortName: "short",
					Link:      "",
				},
			}
			err = k8sClient.Create(ctx, source)
			Expect(err).Should(HaveOccurred())
		})

		It("Should deny if the name or short_name field is longer than 20 characters", func() {
			ctx := context.Background()
			source := &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-long-name",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "asdbhasdbhsadjsdhfaskdasdasdasdhvkjadskhads",
					ShortName: "short",
					Link:      "http://example.com/rss",
				},
			}
			err := k8sClient.Create(ctx, source)
			Expect(err).Should(HaveOccurred())

			source = &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-long-shortname",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "Name",
					ShortName: "asdbhasdbhsadjsdhfaskdasdasdasdhvkjadskhads",
					Link:      "http://example.com/rss",
				},
			}
			err = k8sClient.Create(ctx, source)
			Expect(err).Should(HaveOccurred())
		})

		It("Should deny if the link field is not a valid URL", func() {
			ctx := context.Background()

			source := &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-invalid-url",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "Name",
					ShortName: "short",
					Link:      "invalid-url",
				},
			}
			err := k8sClient.Create(ctx, source)
			Expect(err).Should(HaveOccurred())
		})

		It("Should admit if all required fields are provided and valid", func() {
			ctx := context.Background()

			source := &Source{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-source-valid",
					Namespace: "default",
				},
				Spec: SourceSpec{
					Name:      "ValidName",
					ShortName: "ValidShortName",
					Link:      "http://example.com/rss",
				},
			}

			Expect(k8sClient.Create(ctx, source)).Should(Succeed())

			createdSource := &Source{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "test-source-valid", Namespace: "default"}, createdSource)).Should(Succeed())
		})
	})

	AfterEach(func() {
		ctx := context.Background()
		names := []string{"test-source-invalid-name", "test-source-invalid-shortname", "test-source-invalid-link", "test-source-long-name", "test-source-long-shortname", "test-source-invalid-url", "test-source-valid"}
		for _, name := range names {
			source := &Source{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: "default"}, source)
			if err == nil {
				k8sClient.Delete(ctx, source)
			}
		}
		time.Sleep(1 * time.Second)
	})
})
