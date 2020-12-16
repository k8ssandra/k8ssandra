package tests

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/json"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

var _ = Describe("Verify Datasource template", func() {

	var (
		helmReleaseName       = "k8ssandra-test"
		defaultTestNamespace  = "k8ssandra"
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", defaultTestNamespace)

		helmChartPath string
		err           error
		ds            map[string]interface{}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs("../charts/k8ssandra-cluster")
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) {

		renderedOutput := helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/grafana/datasource.yaml"},
		)
		jsonOutput, err := yaml.YAMLToJSON([]byte(renderedOutput))

		Ω(err).To(BeNil(), "Must convert to json.")
		Ω(json.Unmarshal(jsonOutput, &ds)).To(BeNil(), "Must unmarshal cleanly.")
	}

	Context("by rendering it with options", func() {

		It("using default empty value for routePrefix", func() {

			expectedUrlNoRoutePrefix := "http://" + helmReleaseName + "-prometheus-k8ssandra." + defaultTestNamespace + ":9090"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetStrValues:   map[string]string{"monitoring.prometheus.routePrefix": ""},
			}

			renderTemplate(options)

			spec := ds["spec"].(map[string]interface{})
			Ω(spec).ToNot(BeNil())
			dataSources := spec["datasources"].([]interface{})
			Ω(dataSources).ToNot(BeEmpty())
			Ω(dataSources[0].(map[string]interface{})["url"]).To(BeIdenticalTo(expectedUrlNoRoutePrefix))
		})

		It("using specified routePrefix", func() {

			routePrefix := "prommy"
			expectedUrlWithRoutePrefix := "http://" + helmReleaseName + "-prometheus-k8ssandra." + defaultTestNamespace + ":9090/" + routePrefix
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetStrValues:   map[string]string{"monitoring.prometheus.routePrefix": routePrefix},
			}

			renderTemplate(options)

			spec := ds["spec"].(map[string]interface{})
			Ω(spec).ToNot(BeNil())
			dataSources := spec["datasources"].([]interface{})
			Ω(dataSources).ToNot(BeEmpty())
			Ω(dataSources[0].(map[string]interface{})["url"]).To(BeIdenticalTo(expectedUrlWithRoutePrefix))
		})
	})
})
