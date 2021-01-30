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
	DefaultTestNamespace        = "k8ssandra"
	HelmHookAnnotation          = "helm.sh/hook"
	HelmHookPreDeleteAnnotation = "helm.sh/hook-delete-policy"
	HelmReleaseName             = "k8ssandra-test"
	ReaperInstanceAnnotation    = "reaper.cassandra-reaper.io/instance"
	CommonsTemplate             = "commons.yaml"
)

var (
	defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", DefaultTestNamespace)
)

// Uses commons template to obtain list of required labels for verification.
func GetK8ssandraRequiredLabels() map[string]interface{} {

	var cm map[string]interface{}
	options := &helm.Options{
		KubectlOptions: defaultKubeCtlOptions,
	}

	commonsTemplate, commonsTemplateErr := helm.RenderTemplateE(GinkgoT(), options, ChartsPath,
		HelmReleaseName,
		[]string{"templates/" + CommonsTemplate})

	Expect(commonsTemplateErr).To(BeNil())

	unmarshalErr := helm.UnmarshalK8SYamlE(GinkgoT(), commonsTemplate, &cm)
	Expect(unmarshalErr).To(BeNil())

	requiredLabels := cm["data"].(map[string]interface{})["k8ssandra_labels"].(map[string]interface{})
	Expect(requiredLabels).ToNot(BeEmpty())

	return requiredLabels
}

// Returns k8ssandra templates ignoring helpers and commons.yaml used as expected results.
func GetK8ssandraTemplates(k8ssandraChartPath string) []string {

	var templates []string
	err := filepath.Walk(filepath.Join(k8ssandraChartPath, "templates"),
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
