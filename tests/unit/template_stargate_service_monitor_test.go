package unit_test

import (
	"path/filepath"

	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify Service Monitor template", func() {
	var (
		helmChartPath     string
		err               error
		k8ssandraTemplate map[string]interface{}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		k8ssandraTemplate = map[string]interface{}{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/stargate/service_monitor.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &k8ssandraTemplate)
			})
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"monitoring.serviceMonitors.namespace": "test",
					"stargate.enabled":                     "true",
				},
			}

			err = renderTemplate(options)

			Expect(err).To(BeNil())
			Expect(k8ssandraTemplate["metadata"].(map[string]interface{})["namespace"]).To(Equal("test"))
		})
	})

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "true",
				},
			}

			err = renderTemplate(options)

			Expect(err).To(BeNil())
			Expect(k8ssandraTemplate["metadata"].(map[string]interface{})["namespace"]).To(Equal(DefaultTestNamespace))
		})
	})
})
