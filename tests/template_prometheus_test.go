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

var _ = Describe("Verify Prometheus template", func() {

	var (
		helmReleaseName       = "k8ssandra-test"
		defaultTestNamespace  = "k8ssandra"
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", defaultTestNamespace)

		helmChartPath string
		err           error
		prom          map[string]interface{}
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
			[]string{"templates/prometheus/prometheus.yaml"},
		)
		jsonOutput, err := yaml.YAMLToJSON([]byte(renderedOutput))
		Expect(err).To(BeNil(), "Must convert to json.")
		Expect(json.Unmarshal(jsonOutput, &prom)).To(BeNil(), "Must unmarshal cleanly.")
	}

	Context("by rendering it with options", func() {

		It("using defaults (empty values) for externalUrl and routePrefix", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			renderTemplate(options)
			spec := prom["spec"]
			Expect(spec.(map[string]interface{})["routePrefix"]).To(BeNil())
			Expect(spec.(map[string]interface{})["externalUrl"]).To(BeNil())
		})

		It("using specific externaUrl and routePrefix", func() {

			testExternalUrl := "http://foobar.com:8675"
			testRoutePrefix := "prommy"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetStrValues: map[string]string{"monitoring.prometheus.externalUrl": testExternalUrl,
					"monitoring.prometheus.routePrefix": testRoutePrefix},
			}

			renderTemplate(options)
			spec := prom["spec"]
			Expect(spec.(map[string]interface{})["routePrefix"]).To(BeIdenticalTo(testRoutePrefix))
			Expect(spec.(map[string]interface{})["externalUrl"]).To(BeIdenticalTo(testExternalUrl))
		})
	})
})
