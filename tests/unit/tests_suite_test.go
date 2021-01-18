package unit_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	chartsPath                  = ("../../charts/k8ssandra")
	reaperInstanceAnnotation    = "reaper.cassandra-reaper.io/instance"
	helmHookAnnotation          = "helm.sh/hook"
	helmHookPreDeleteAnnotation = "helm.sh/hook-delete-policy"
	defaultTestNamespace        = "k8ssandra"
	helmReleaseName             = "k8ssandra-test"
)

func TestTemplateUnitTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit tests suite")
}

var _ = BeforeSuite(func(done Done) {
	close(done)
})
