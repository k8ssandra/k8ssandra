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

var (
	requiredLabels = map[string]string{
		"helm.sh/chart":                "k8ssandra-0.24.0",
		"app.kubernetes.io/name":       "k8ssandra",
		"app.kubernetes.io/instance":   "k8ssandra-test",
		"app.kubernetes.io/version":    "3.11.7",
		"app.kubernetes.io/managed-by": "Helm",
	}
)

// validateRequiredLabels supports validation of k8ssandra required labels
func validateRequiredLabels(existingLabels interface{}) {

	for k, v := range requiredLabels {
		Expect(existingLabels).To(HaveKeyWithValue(k, v))
	}
}

func TestTemplateUnitTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit tests suite")
}

var _ = BeforeSuite(func(done Done) {
	close(done)
})
