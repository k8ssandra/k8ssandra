package unit_test

import (
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify dashboard Config Map template", func() {

	var (
		HelmReleaseName       = "k8ssandra-test"
		DefaultTestNamespace  = "k8ssandra"
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", DefaultTestNamespace)

		helmChartPath string
		err           error
		cm            map[string]interface{}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/grafana/configmap.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &cm)
			})
	}

	Context("by rendering it with options", func() {

		It("using provision_dashboards false", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues:      map[string]string{"monitoring.grafana.provision_dashboards": "false"},
			}

			err = renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})

		It("using provision_dashboards true", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues:      map[string]string{"monitoring.grafana.provision_dashboards": "true"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			data := cm["data"].(map[string]interface{})
			Expect(data).ToNot(BeNil())
			Expect(data).ToNot(BeEmpty())
			Expect(data["cassandra-condensed.json"]).ToNot(BeEmpty())
			Expect(data["cassandra-overview.json"]).ToNot(BeEmpty())
			Expect(data["system-metrics.json"]).ToNot(BeEmpty())
		})
	})
})
