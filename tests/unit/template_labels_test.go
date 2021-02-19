package unit_test

import (
	"path/filepath"
	"strings"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify k8ssandra and dependent template labels", func() {

	var (
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "",
			DefaultTestNamespace)
		localChartsPath string
	)

	BeforeEach(func() {
		localChartsPath = ""
	})

	Context("by rendering k8ssandra templates having common labels", func() {
		It("using all enabled options", func() {

			path, err := filepath.Abs(ChartsPath)
			Expect(err).To(BeNil())

			localChartsPath = path

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                              "true",
					"reaper.enabled":                                "true",
					"medusa.enabled":                                "true",
					"reaper.ingress.enabled":                        "true",
					"reaper.ingress.host":                           "reaper.host",
					"cassandra.ingress.enabled":                     "false",
					"ingress.traefik.monitoring.grafana.enabled":    "true",
					"ingress.traefik.monitoring.prometheus.enabled": "true",
					"stargate.ingress.enabled":                      "true",
					"stargate.ingress.auth.enabled":                 "true",
					"stargate.ingress.rest.enabled":                 "true",
					"stargate.ingress.graphql.enabled":              "true",
					"stargate.ingress.graphql.playground.enabled":   "true",
					"cassandra.auth.enabled":                        "true",
					"cassandra.auth.superuser.username":             "admin",
					"cassandra.clusterName":                         "test-cluster",
				},
			}

			// Verify required labels for ea. template
			requiredLabels := GetRequiredLabels(localChartsPath)
			templates := GetTemplates(localChartsPath)
			for _, template := range templates {

				var k8ssandraTemplates map[string]interface{}
				idx := strings.Index(template, "templates")

				if template[idx:] == "templates/cassandra/ingress.yaml" {
					continue
				}

				templateOutput, err := helm.RenderTemplateE(GinkgoT(), options,
					localChartsPath, HelmReleaseName, []string{filepath.Join(".", template[idx:])})

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

	Context("by rendering cass-operator templates having k8ssandra common labels", func() {
		It("using default options", func() {

			path, err := filepath.Abs(CassOperatorChartsPath)
			Expect(err).To(BeNil())
			localChartsPath = path

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			// Verify required labels for ea. template
			requiredLabels := GetRequiredLabels(localChartsPath)
			templates := GetTemplates(localChartsPath)
			for _, template := range templates {

				var k8ssandraTemplates map[string]interface{}
				idx := strings.Index(template, "templates")

				templateOutput, err := helm.RenderTemplateE(GinkgoT(), options,
					localChartsPath, HelmReleaseName, []string{filepath.Join(".", template[idx:])})

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

	Context("by rendering medusa-operator templates having k8ssandra common labels", func() {
		It("using default options", func() {

			path, err := filepath.Abs(MedusaOperatorChartsPath)
			Expect(err).To(BeNil())
			localChartsPath = path

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			// Verify required labels for ea. template
			requiredLabels := GetRequiredLabels(localChartsPath)
			templates := GetTemplates(localChartsPath)
			for _, template := range templates {

				var k8ssandraTemplates map[string]interface{}
				idx := strings.Index(template, "templates")

				templateOutput, err := helm.RenderTemplateE(GinkgoT(), options,
					localChartsPath, HelmReleaseName, []string{filepath.Join(".", template[idx:])})

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

	Context("by rendering reaper-operator templates having k8ssandra common labels", func() {
		It("using default options", func() {

			path, err := filepath.Abs(ReaperOperatorChartsPath)
			Expect(err).To(BeNil())
			localChartsPath = path

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			// Verify required labels for ea. template
			requiredLabels := GetRequiredLabels(localChartsPath)
			templates := GetTemplates(localChartsPath)
			for _, template := range templates {

				var k8ssandraTemplates map[string]interface{}
				idx := strings.Index(template, "templates")

				templateOutput, err := helm.RenderTemplateE(GinkgoT(), options,
					localChartsPath, HelmReleaseName, []string{filepath.Join(".", template[idx:])})

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
