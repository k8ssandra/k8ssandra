package unit_test

import (
	. "fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"helm.sh/helm/v3/pkg/chartutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	ChartsPath                  = "../../charts/k8ssandra"
	CassOperatorChartsPath      = "../../charts/cass-operator"
	MedusaOperatorChartsPath    = "../../charts/medusa-operator"
	ReaperOperatorChartsPath    = "../../charts/reaper-operator"
	HelmHookAnnotation          = "helm.sh/hook"
	HelmHookPreDeleteAnnotation = "helm.sh/hook-delete-policy"
	ReaperInstanceAnnotation    = "reaper.cassandra-reaper.io/instance"
)

var (
	UniqueIdSuffix        = strings.ToLower(random.UniqueId())
	DefaultTestNamespace  = Sprintf("k8ssandranamespace%s", UniqueIdSuffix)
	HelmReleaseName       = Sprintf("k8ssandratestrelease%s", UniqueIdSuffix)
	defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", DefaultTestNamespace)
)

// Uses commons template to obtain list of required labels for verification.
func GetRequiredLabels(targetChartsPath string) map[string]interface{} {

	chartMetadata, _ := chartutil.LoadChartfile(filepath.Join(targetChartsPath, "Chart.yaml"))
	Expect(chartMetadata).ToNot(BeNil())

	// k8ssandra-common.labels
	commonLabels := map[string]interface{}{
		"helm.sh/chart":                chartMetadata.Name + "-" + chartMetadata.Version,
		"app.kubernetes.io/name":       chartMetadata.Name,
		"app.kubernetes.io/instance":   HelmReleaseName,
		"app.kubernetes.io/managed-by": "Helm",
		"app.kubernetes.io/part-of":    "k8ssandra" + "-" + HelmReleaseName + "-" + DefaultTestNamespace,
	}
	// k8ssandra.lables includes version label in addition to k8ssandra-common.labels
	if targetChartsPath == ChartsPath {
		commonLabels["app.kubernetes.io/version"] = chartMetadata.AppVersion
	}
	return commonLabels
}

// Returns templates ignoring helpers and commons.yaml used as expected results.
func GetTemplates(targetChartsPath string) []string {

	Expect(targetChartsPath).ToNot(BeNil())
	var templates []string

	err := filepath.Walk(filepath.Join(targetChartsPath, "templates"),
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
				absPath, err := filepath.Abs(path)
				Expect(err).To(BeNil())
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
