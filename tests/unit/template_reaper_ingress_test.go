package unit_test

import (
	. "fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/traefik"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"path/filepath"
)

var _ = Describe("Verify Reaper ingress template", func() {
	var (
		helmChartPath string
		err           error
		ingress       traefik.IngressRoute
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		ingress = traefik.IngressRoute{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/reaper/ingress.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &ingress)
			})
	}

	Context("by confirming it does not render when", func() {
		It("is implicitly disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is explicitly disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRoute"))
			VerifyTraefikHTTPIngressRoute(ingress, "web", "Host(`repair.k8ssandra.cluster.local`)", Sprintf("%s-reaper-k8ssandra-reaper-service", HelmReleaseName), 8080)
		})

		It("with custom host", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":     "true",
					"ingress.traefik.repair.host": "reaper.host",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRoute"))
			VerifyTraefikHTTPIngressRoute(ingress, "web", "Host(`reaper.host`)", Sprintf("%s-reaper-k8ssandra-reaper-service", HelmReleaseName), 8080)
		})
	})
})
