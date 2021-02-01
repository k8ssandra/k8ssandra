package unit_test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	DefaultTestNamespace        = "k8ssandra"
	HelmHookAnnotation          = "helm.sh/hook"
	HelmHookPreDeleteAnnotation = "helm.sh/hook-delete-policy"
	HelmReleaseName             = "k8ssandra-test"
	ReaperInstanceAnnotation    = "reaper.cassandra-reaper.io/instance"
	CommonsTemplate             = "commons.yaml"
	CommonLabelKey              = "common_labels"
)

var (
	defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", DefaultTestNamespace)
)

// Uses commons template to obtain list of required labels for verification.
func GetRequiredLabels(targetChartsPath string) map[string]interface{} {

	Expect(targetChartsPath).ToNot(BeNil())
	var configMap map[string]interface{}
	options := &helm.Options{
		KubectlOptions: defaultKubeCtlOptions,
	}

	renderedTemplate, renderedTemplateErr := helm.RenderTemplateE(GinkgoT(), options,
		targetChartsPath,
		HelmReleaseName,
		[]string{"templates/" + CommonsTemplate})
	Expect(renderedTemplateErr).To(BeNil())
	Expect(renderedTemplate).ToNot(BeNil())

	Expect(helm.UnmarshalK8SYamlE(GinkgoT(), renderedTemplate, &configMap)).To(BeNil())

	requiredLabels := configMap["data"].(map[string]interface{})[CommonLabelKey].(map[string]interface{})
	Expect(requiredLabels).ToNot(BeEmpty())

	return requiredLabels
}

// Returns templates ignoring helpers and commons.yaml used as expected results.
func GetTemplates(targetChartsPath string) []string {

	Expect(targetChartsPath).ToNot(BeNil())
	var templates []string

	err := filepath.Walk(filepath.Join(targetChartsPath, "templates"),
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") && CommonsTemplate != info.Name() {
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
