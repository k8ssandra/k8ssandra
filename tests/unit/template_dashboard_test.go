package unit_test

import (
	"path/filepath"

	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify dashboards Config Map templates", func() {

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

	renderTemplate := func(options *helm.Options, templatePath string) error {
		return helmUtils.RenderAndUnmarshall(templatePath,
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &cm)
			})
	}

	Context("by rendering it with options", func() {

		Context("provision_dashboards false", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues:      map[string]string{"monitoring.grafana.provision_dashboards": "false"},
			}

			It("does not render the Cassandra Overview dashboard", func() {
				err = renderTemplate(options, "templates/cassandra/overview-dashboard.yaml")

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("could not find template"))
			})

			It("does not render the Cassandra Condensed dashboard", func() {
				err = renderTemplate(options, "templates/cassandra/condensed-dashboard.yaml")

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("could not find template"))
			})

			It("does not render the System Metrics dashboard", func() {
				err = renderTemplate(options, "templates/cassandra/system-metrics-dashboard.yaml")

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("could not find template"))
			})

			It("does not render the Stargate dashboard", func() {
				err = renderTemplate(options, "templates/stargate/api-dashboard.yaml")

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("could not find template"))
			})
		})

		Context("provision_dashboards true", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues:      map[string]string{"monitoring.grafana.provision_dashboards": "true"},
			}

			It("renders the Cassandra Overview dashboard", func() {
				Expect(renderTemplate(options, "templates/cassandra/overview-dashboard.yaml")).To(Succeed())

				data := cm["data"].(map[string]interface{})
				Expect(data).ToNot(BeNil())
				Expect(data).ToNot(BeEmpty())
				Expect(data["cassandra-overview.json"]).ToNot(BeEmpty())
			})

			It("renders the Cassandra Condensed dashboard", func() {
				Expect(renderTemplate(options, "templates/cassandra/condensed-dashboard.yaml")).To(Succeed())

				data := cm["data"].(map[string]interface{})
				Expect(data).ToNot(BeNil())
				Expect(data).ToNot(BeEmpty())
				Expect(data["cassandra-condensed.json"]).ToNot(BeEmpty())
			})

			It("renders the System Metrics dashboard", func() {
				Expect(renderTemplate(options, "templates/cassandra/system-metrics-dashboard.yaml")).To(Succeed())

				data := cm["data"].(map[string]interface{})
				Expect(data).ToNot(BeNil())
				Expect(data).ToNot(BeEmpty())
				Expect(data["system-metrics.json"]).ToNot(BeEmpty())
			})

			Context("and stargate.enabled false", func() {
				options := &helm.Options{
					KubectlOptions: defaultKubeCtlOptions,
					SetValues: map[string]string{
						"monitoring.grafana.provision_dashboards": "true",
						"stargate.enabled":                        "false",
					},
				}

				It("does not render the Stargate dashboard", func() {
					err = renderTemplate(options, "templates/stargate/api-dashboard.yaml")

					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(ContainSubstring("could not find template"))
				})
			})

			Context("and stargate.enabled true", func() {
				options := &helm.Options{
					KubectlOptions: defaultKubeCtlOptions,
					SetValues: map[string]string{
						"monitoring.grafana.provision_dashboards": "true",
						"stargate.enabled":                        "true",
					},
				}

				It("renders the Stargate dashboard", func() {
					Expect(renderTemplate(options, "templates/stargate/api-dashboard.yaml")).To(Succeed())

					data := cm["data"].(map[string]interface{})
					Expect(data).ToNot(BeNil())
					Expect(data).ToNot(BeEmpty())
					Expect(data["stargate.json"]).ToNot(BeEmpty())
				})
			})
		})
	})
})
