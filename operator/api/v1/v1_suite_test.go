package v1_test

import (
	v1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestV1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V1 Suite")
}

var _ = BeforeSuite(func() {
	_ = v1.AddToScheme(scheme.Scheme)
})
