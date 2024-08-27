package controller_test

import (
	v1 "com.teamdev/news-aggregator/api/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	_ = v1.AddToScheme(scheme.Scheme)
})
