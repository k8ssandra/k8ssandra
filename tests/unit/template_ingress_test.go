package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify Ingress templates", func() {
	var (
		helmChartPath string
		err           error
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) {
		helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/traefik.ingressroutes.yaml"},
		)

		// helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, ingress)
	}

	Context("by rendering it with options", func() {
		It("using only enabled option", func() {
			options := &helm.Options{
				SetStrValues:   map[string]string{"ingress.traefik.enabled": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
		})

		It("enabling the TLS", func() {
			options := &helm.Options{
				SetStrValues: map[string]string{
					"ingress.traefik.enabled":               "true",
					"ingress.traefik.tls.options.name":      "custom-tls",
					"ingress.traefik.tls.options.namespace": "current",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
		})
	})
})
