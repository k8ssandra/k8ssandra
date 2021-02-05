package unit_test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	networking "k8s.io/api/networking/v1beta1"
	"path/filepath"
)

var _ = Describe("Verify Stargate ingress template", func() {
	var (
		helmChartPath string
		err           error
		ingress       networking.Ingress
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		ingress = networking.Ingress{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/stargate/ingress.yaml",
			options, helmChartPath, helmReleaseName,
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
					"ingress.traefik.enabled": "true",
					"ingress.traefik.stargate.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is explicitly disabled at the Ingress level even when enabled at the Stargate level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "false",
					"ingress.traefik.stargate.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by confirming it fails when", func() {
		It("stargate host is not defined", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled":          "true",
					"ingress.traefik.cassandra.enabled": "true",
					"ingress.traefik.stargate.enabled": "true",
					"ingress.traefik.stargate.host": "",
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
					"ingress.traefik.cassandra.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, "*", true, false, true)
		})

		It("with everything enabled and default settings", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "true",
					"ingress.traefik.cassandra.enabled": "false",
					"ingress.traefik.stargate.graphql.playground.enabled": "true",
					"ingress.traefik.stargate.cassandra.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, "*", true, true, true)
		})

		It("with everything enabled and custom settings", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"ingress.traefik.enabled": "true",
					"ingress.traefik.cassandra.enabled": "false",
					"ingress.traefik.stargate.host": "stargate.host",
					"ingress.traefik.stargate.graphql.playground.enabled": "true",
					"ingress.traefik.stargate.cassandra.enabled": "true",
					"ingress.traefik.stargate.rest.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, "stargate.host", true, true, true)

		})
	})
})

func verifyIngressRules(ingress networking.Ingress, host string, graphEnabled bool, playgroundEnabled bool, restEnabled bool) {
	rules := ingress.Spec.Rules
	kubeapi.VerifyIngressRule(rules, helmReleaseName,"/v1/auth", nil, host, 8081)
	if graphEnabled {
		kubeapi.VerifyIngressRule(rules, helmReleaseName,"/graphql/", nil, host, 8080)
		kubeapi.VerifyIngressRule(rules, helmReleaseName,"/graphql-schema", nil, host, 8080)
		if playgroundEnabled {
			pathType := networking.PathTypeExact
			kubeapi.VerifyIngressRule(rules, helmReleaseName,"/playground", &pathType, host, 8080)
		} else {
			kubeapi.VerifyNoRuleWithPath(rules, "/playground")
		}
	} else {
		kubeapi.VerifyNoRuleWithPath(rules, "/graphql/")
		kubeapi.VerifyNoRuleWithPath(rules, "/graphql-schema")
		kubeapi.VerifyNoRuleWithPath(rules, "/playground")
	}
	if restEnabled {
		kubeapi.VerifyIngressRule(rules, helmReleaseName,"/v2/", nil, host, 8082)
	} else {
		kubeapi.VerifyNoRuleWithPath(rules, "/v2/")
	}
}

