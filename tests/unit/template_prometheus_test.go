package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify Service Monitor template", func() {

	var (
		HelmReleaseName       = "k8ssandra-test"
		DefaultTestNamespace  = "k8ssandra"
		defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", DefaultTestNamespace)

		helmChartPath string
		err           error
		sm            map[string]interface{}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {

		renderedOutput, err := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, HelmReleaseName,
			[]string{"templates/prometheus/service_monitor.yaml"},
		)

		if err == nil {
			helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, &sm)
		}

		return err
	}

	Context("by rendering it with options", func() {

		It("using provision_service_monitors false", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues:      map[string]string{"monitoring.prometheus.provision_service_monitors": "false"},
			}

			err = renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})

		It("using provision_service_monitors true", func() {

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues:      map[string]string{"monitoring.prometheus.provision_service_monitors": "true"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			meta := sm["metadata"]
			Expect(meta.(map[string]interface{})["labels"].(map[string]interface{})["release"]).To(BeIdenticalTo(HelmReleaseName))

			spec := sm["spec"]
			Expect(spec).ToNot(BeEmpty())
		})
	})
})
