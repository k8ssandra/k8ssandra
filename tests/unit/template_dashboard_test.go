package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify dashboard Config Map template", func() {

	var (
		helmReleaseName       = "k8ssandra-test"
		defaultTestNamespace  = "k8ssandra"
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", defaultTestNamespace)

		helmChartPath string
		err           error
		cm            map[string]interface{}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {

		renderedOutput, err := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/grafana/configmap.yaml"},
		)

		if err == nil {
			helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, &cm)
		}

		return err
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
