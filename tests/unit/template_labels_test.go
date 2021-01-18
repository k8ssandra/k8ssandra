package unit_test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
	"strings"
)

var _ = Describe("Verify k8ssandra template labels", func() {

	var (
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "",
			defaultTestNamespace)
		k8ssandraChartPath string
	)

	BeforeEach(func() {
		path, err := filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		k8ssandraChartPath = path
	})

	Context("by rendering it with options", func() {

		It("using only default options", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                              "true",
					"repair.reaper.enabled":                         "true",
					"backupRestore.medusa.enabled":                  "true",
					"ingress.traefik.enabled":                       "true",
					"ingress.traefik.monitoring.grafana.enabled":    "true",
					"ingress.traefik.monitoring.prometheus.enabled": "true",
				},
			}

			// Verify required labels for ea. k8ssandra template
			requiredLabels := GetK8ssandraRequiredLabels()
			templates := GetK8ssandraTemplates(k8ssandraChartPath)
			for _, template := range templates {

				var k8ssandraTemplates map[string]interface{}
				idx := strings.Index(template, "templates")

				targetTemplate := template[idx:]
				templateOutput, err := helm.RenderTemplateE(GinkgoT(), options,
					k8ssandraChartPath, helmReleaseName, []string{targetTemplate})

				Expect(err).To(BeNil())
				Expect(templateOutput).ToNot(BeEmpty())
				Expect(helm.UnmarshalK8SYamlE(GinkgoT(), templateOutput, &k8ssandraTemplates)).To(BeNil())

				Expect(k8ssandraTemplates["metadata"]).ToNot(BeNil())
				for k, v := range requiredLabels {
					Expect(k8ssandraTemplates["metadata"].(map[string]interface{})["labels"]).To(HaveKeyWithValue(k, v))
				}
			}

		})
	})
})
