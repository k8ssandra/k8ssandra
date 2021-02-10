package unit_test

import (
	fmt "fmt"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/traefik"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

var _ = Describe("Verify Stargate Cassandra ingress template", func() {
	var (
		helmChartPath string
		err           error
		ingress       traefik.IngressRouteTCP
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		ingress = traefik.IngressRouteTCP{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/stargate/cassandra-ingress.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &ingress)
			})
	}

	Context("by confirming it does render when", func() {
		It("is implicitly enabled at the Ingress level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).Should(Succeed())
		})
	})

	Context("by confirming it does not render when", func() {
		It("is explicitly disabled at the Ingress level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is explicitly disabled at the Stargate level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":          "true",
					"ingress.traefik.stargate.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is explicitly disabled at the Stargate Cassandra level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":                    "true",
					"ingress.traefik.stargate.enabled":           "true",
					"ingress.traefik.stargate.cassandra.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is explicitly disabled at the Stargate level even when enabled at the Cassandra level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":           "true",
					"ingress.traefik.cassandra.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is explicitly disabled at the Ingress level even when enabled at the Cassandra level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":                    "false",
					"ingress.traefik.stargate.enabled":           "true",
					"ingress.traefik.stargate.cassandra.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by confirming it fails when", func() {
		It("cassandra ingress is enabled for both Stargate and non-Stargate", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":                    "true",
					"ingress.traefik.cassandra.enabled":          "true",
					"ingress.traefik.stargate.enabled":           "true",
					"ingress.traefik.stargate.cassandra.enabled": "true",
				},
			}
			renderErr := renderTemplate(options)
			Expect(renderErr).ToNot(BeNil())
			Expect(renderErr.Error()).To(ContainSubstring("ingress.traefik.stargate.cassandra.enabled and ingress.traefik.cassandra.enabled cannot both be enabled"))
		})
		It("stargate host is not defined", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":           "true",
					"ingress.traefik.cassandra.enabled": "true",
					"ingress.traefik.stargate.enabled":  "true",
					"ingress.traefik.stargate.host":     "",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by rendering it when", func() {
		It("is enabled with default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress,
				"cassandra",
				"HostSNI(`*`)",
				fmt.Sprintf("%s-%s-stargate-service", HelmReleaseName, "dc1"),
				9042)
		})

		It("is enabled with custom host", func() {
			stargateHost := "stargate.host"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":                    "true",
					"ingress.traefik.cassandra.enabled":          "false",
					"ingress.traefik.stargate.enabled":           "true",
					"ingress.traefik.stargate.cassandra.enabled": "true",
					"ingress.traefik.stargate.host":              stargateHost,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress,
				"cassandra",
				fmt.Sprintf("HostSNI(`%s`)", stargateHost),
				fmt.Sprintf("%s-%s-stargate-service", HelmReleaseName, "dc1"),
				9042)
		})
	})
})
