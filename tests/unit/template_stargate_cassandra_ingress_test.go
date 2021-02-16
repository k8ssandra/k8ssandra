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

	Context("by confirming it does not render when", func() {
		It("is explicitly disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.ingress.cassandra.enabled": "false",
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
					"cassandra.ingress.enabled":          "true",
					"stargate.ingress.cassandra.enabled": "true",
				},
			}
			renderErr := renderTemplate(options)
			Expect(renderErr).ToNot(BeNil())
			Expect(renderErr.Error()).To(ContainSubstring("stargate.ingress.cassandra.enabled and cassandra.ingress.enabled cannot both be true"))
		})
	})

	Context("by confirming", func() {
		It("is enabled by default", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(renderTemplate(options)).Should(Succeed())
		})
	})

	Context("by rendering it when", func() {
		It("is enabled with default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.ingress.enabled":          "false",
					"stargate.ingress.cassandra.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress,
				"cassandra",
				"HostSNI(`*`)",
				Sprintf("%s-%s-stargate-service", HelmReleaseName, "dc1"),
				9042)
		})

		It("is enabled and release name != cluster name", func() {
			clusterName := Sprintf("k8ssandraclustername%s", UniqueIdSuffix)
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":              clusterName,
					"stargate.ingress.cassandra.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress,
				"cassandra",
				"HostSNI(`*`)",
				Sprintf("%s-%s-stargate-service", HelmReleaseName, "dc1"),
				9042)
		})

		It("is enabled with custom host", func() {
			stargateHost := "stargate.host"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.ingress.host": stargateHost,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress,
				"cassandra",
				Sprintf("HostSNI(`%s`)", stargateHost),
				Sprintf("%s-%s-stargate-service", HelmReleaseName, "dc1"),
				9042)
		})
	})
})
