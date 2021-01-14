package unit_test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("Verify Grafana template", func() {

	var (
		helmReleaseName       = "k8ssandra-test"
		defaultTestNamespace  = "k8ssandra"
		defaultKubeCtlOptions = k8s.NewKubectlOptions(
			"", "", defaultTestNamespace)

		helmChartPath string
		err           error
		dashBoard     map[string]interface{}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options, yamlFile string) {

		renderedOutput := helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{yamlFile},
		)
		helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, &dashBoard)
	}

	Context("by rendering dashboards with options", func() {

		options := &helm.Options{
			KubectlOptions: defaultKubeCtlOptions,
		}

		It("having required labels in cassandra-condensed dashboard", func() {
			renderTemplate(options, "templates/grafana/dashboards/cassandra-condensed.dashboard-helm-template.yaml")
			Expect(dashBoard["metadata"]).ToNot(BeNil())
			validateRequiredLabels(dashBoard["metadata"].(map[string]interface{})["labels"])
		})
		It("having required labels in cassandra-condensed dashboard", func() {
			renderTemplate(options, "templates/grafana/dashboards/overview-with-plugin.dashboard-helm-template.yaml")
			Expect(dashBoard["metadata"]).ToNot(BeNil())
			validateRequiredLabels(dashBoard["metadata"].(map[string]interface{})["labels"])
		})
		It("having required labels in cassandra-condensed dashboard", func() {
			renderTemplate(options, "templates/grafana/dashboards/system-metrics.dashboard-helm-template.yaml")
			Expect(dashBoard["metadata"]).ToNot(BeNil())
			validateRequiredLabels(dashBoard["metadata"].(map[string]interface{})["labels"])
		})
	})
})
