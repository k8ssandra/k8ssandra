package unit_test

import (
	"k8s.io/helm/pkg/chartutil"
	"os"
	"path/filepath"
	"strings"
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

// Common labels for k8ssandra chart
func GetK8ssandraRequiredLabels() map[string]string {

	k8ssandraChart, err := chartutil.Load(chartsPath)
	Expect(err).ToNot(BeNil())
	Expect(k8ssandraChart).ToNot(BeNil())

	meta := k8ssandraChart.Metadata
	requiredLabels := map[string]string{
		"helm.sh/chart":                "k8ssandra-" + meta.Version,
		"app.kubernetes.io/name":       "k8ssandra",
		"app.kubernetes.io/instance":   helmReleaseName,
		"app.kubernetes.io/version":    meta.AppVersion,
		"app.kubernetes.io/managed-by": "Helm",
	}
	return requiredLabels
}

func GetK8ssandraTemplates(k8ssandraChartPath string) []string {

	var templates []string
	templatesPath := k8ssandraChartPath + string(os.PathSeparator) + "templates"

	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			absPath, _ := filepath.Abs(path)
			templates = append(templates, absPath)
		}
		return nil
	})
	Expect(err).To(BeNil())
	return templates
}

func TestTemplateUnitTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit tests suite")
}

var _ = BeforeSuite(func(done Done) {
	close(done)
})
