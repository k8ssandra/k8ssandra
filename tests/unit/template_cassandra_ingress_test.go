package unit_test

import (
	. "fmt"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/traefik"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

var _ = Describe("Verify Cassandra ingress template", func() {
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
		return helmUtils.RenderAndUnmarshall("templates/cassandra/ingress.yaml",
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
					"cassandra.ingress.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
		It("method is not traefik", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.ingress.enabled": "true",
					"cassandra.ingress.method":  "somethingElse",
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
					"stargate.enabled":                   "true",
					"stargate.ingress.enabled":           "true",
					"stargate.ingress.cassandra.enabled": "true",
				},
			}
			renderErr := renderTemplate(options)
			Expect(renderErr).ToNot(BeNil())
			Expect(renderErr.Error()).To(ContainSubstring("stargate.ingress.cassandra.enabled and cassandra.ingress.enabled cannot both be true"))
		})
	})

	Context("by rendering it when", func() {
		It("it is enabled and Stargate Cassandra ingress is disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.ingress.enabled":          "true",
					"stargate.ingress.cassandra.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra", "HostSNI(`*`)", Sprintf("%s-%s-service", HelmReleaseName, "dc1"), 9042)
		})

		It("it is enabled and Stargate ingress is disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.ingress.enabled":          "true",
					"stargate.enabled":                   "true",
					"stargate.ingress.enabled":           "false",
					"stargate.ingress.cassandra.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra", "HostSNI(`*`)", Sprintf("%s-%s-service", HelmReleaseName, "dc1"), 9042)
		})

		It("it is enabled and Stargate is disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.ingress.enabled":          "true",
					"stargate.enabled":                   "false",
					"stargate.ingress.enabled":           "true",
					"stargate.ingress.cassandra.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra", "HostSNI(`*`)", Sprintf("%s-%s-service", HelmReleaseName, "dc1"), 9042)
		})

		It("it is enabled with a custom host", func() {
			cassandraHost := "cassandra.host"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":          "false",
					"cassandra.ingress.enabled": "true",
					"cassandra.ingress.host":    cassandraHost,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra",
				Sprintf("HostSNI(`%s`)", cassandraHost),
				Sprintf("%s-%s-service", HelmReleaseName, "dc1"),
				9042)
		})

		It("it is enabled with host=empty string", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":          "false",
					"cassandra.ingress.enabled": "true",
					"cassandra.ingress.host":    "",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra",
				"HostSNI(`*`)",
				Sprintf("%s-%s-service", HelmReleaseName, "dc1"),
				9042)
		})

		It("it is enabled with host=nil", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/cassandra-ingress-nil-host.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra",
				"HostSNI(`*`)",
				Sprintf("%s-%s-service", HelmReleaseName, "dc1"),
				9042)
		})

		It("it is enabled and Stargate Cassandra ingress is disabled with release name != cluster name", func() {
			clusterName := Sprintf("k8ssandracluster%s", UniqueIdSuffix)
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":              clusterName,
					"cassandra.ingress.enabled":          "true",
					"stargate.ingress.cassandra.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("IngressRouteTCP"))

			VerifyTraefikTCPIngressRoute(ingress, "cassandra", "HostSNI(`*`)", Sprintf("%s-%s-service", clusterName, "dc1"), 9042)
		})
	})
})
