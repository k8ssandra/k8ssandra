package unit_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	chartsPath                  = ("../../charts/k8ssandra-cluster")
	reaperInstanceAnnotation    = "reaper.cassandra-reaper.io/instance"
	helmHookAnnotation          = "helm.sh/hook"
	helmHookPreDeleteAnnotation = "helm.sh/hook-delete-policy"
	defaultTestNamespace        = "k8ssandra"
	helmReleaseName             = "k8ssandra-test"
)

func TestTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tests Suite")
}

var _ = BeforeSuite(func(done Done) {
	close(done)
})
