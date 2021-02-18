package unit_test

import (
	. "fmt"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	networking "k8s.io/api/networking/v1beta1"
)

var _ = Describe("Verify Stargate ingress template", func() {
	var (
		helmChartPath string
		err           error
		ingress       networking.Ingress
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		ingress = networking.Ingress{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/stargate/ingress.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &ingress)
			})
	}

	Context("by confirming it does not render when", func() {
		It("is disabled at the Stargate level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "false",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "true",
					"stargate.ingress.rest.enabled":    "true",
					"stargate.ingress.graphql.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is disabled at the Stargate ingress level", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "false",
					"stargate.ingress.auth.enabled":    "true",
					"stargate.ingress.rest.enabled":    "true",
					"stargate.ingress.graphql.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("is disabled at the individual ingress levels", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "false",
					"stargate.ingress.rest.enabled":    "false",
					"stargate.ingress.graphql.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by confirming", func() {
		It("is disabled by default", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by rendering it with options", func() {
		It("using only default options and Stargate ingress enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.ingress.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, true, nil, true, true, nil, true, nil)
		})

		It("with everything explicitly enabled and custom host", func() {
			stargateHost := "stargate.host"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                            "true",
					"stargate.ingress.enabled":                    "true",
					"stargate.ingress.cassandra.enabled":          "false",
					"stargate.ingress.host":                       stargateHost,
					"stargate.ingress.graphql.playground.enabled": "true",
					"stargate.ingress.rest.enabled":               "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, true, &stargateHost, true, true, &stargateHost, true, &stargateHost)

		})

		It("with everything explicitly enabled and wildcard host", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                            "true",
					"stargate.ingress.enabled":                    "true",
					"stargate.ingress.host":                       "*",
					"stargate.ingress.graphql.playground.enabled": "true",
					"stargate.ingress.rest.enabled":               "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, true, nil, true, true, nil, true, nil)

		})

		It("with only auth enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "true",
					"stargate.ingress.rest.enabled":    "false",
					"stargate.ingress.graphql.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, true, nil, false, false, nil, false, nil)

		})

		It("with only rest enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "false",
					"stargate.ingress.rest.enabled":    "true",
					"stargate.ingress.graphql.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, false, nil, false, false, nil, true, nil)

		})

		It("with only graphql enabled and playground implicitly enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "false",
					"stargate.ingress.rest.enabled":    "false",
					"stargate.ingress.graphql.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, false, nil, true, true, nil, false, nil)

		})

		It("with only graphql enabled and playground explicitly enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                            "true",
					"stargate.ingress.enabled":                    "true",
					"stargate.ingress.auth.enabled":               "false",
					"stargate.ingress.rest.enabled":               "false",
					"stargate.ingress.graphql.enabled":            "true",
					"stargate.ingress.graphql.playground.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, false, nil, true, true, nil, false, nil)

		})

		It("with only graphql and playground disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                            "true",
					"stargate.ingress.enabled":                    "true",
					"stargate.ingress.auth.enabled":               "false",
					"stargate.ingress.rest.enabled":               "false",
					"stargate.ingress.graphql.enabled":            "true",
					"stargate.ingress.graphql.playground.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, false, nil, true, false, nil, false, nil)

		})

		It("with everything enabled and varying hosts", func() {
			stargateDefaultHost := "stargate.host" // we're not setting a host for rest, so it should inherit this
			stargateGraphqlHost := "graphql.host"  // this should override the default
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "true",
					"stargate.ingress.rest.enabled":    "true",
					"stargate.ingress.graphql.enabled": "true",
					"stargate.ingress.host":            stargateDefaultHost,
					"stargate.ingress.auth.host":       "", // accept auth requests from any host, overriding the default
					"stargate.ingress.graphql.host":    stargateGraphqlHost,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, true, nil, true, true, &stargateGraphqlHost, true, &stargateDefaultHost)

		})

		It("with everything enabled, no default host, and some overridden hosts", func() {
			stargateRestHost := "rest.host"
			stargateGraphqlHost := "graphql.host"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                 "true",
					"stargate.ingress.enabled":         "true",
					"stargate.ingress.auth.enabled":    "true",
					"stargate.ingress.rest.enabled":    "true",
					"stargate.ingress.graphql.enabled": "true",
					"stargate.ingress.rest.host":       stargateRestHost,
					"stargate.ingress.graphql.host":    stargateGraphqlHost,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(ingress.Kind).To(Equal("Ingress"))
			verifyIngressRules(ingress, true, nil, true, true, &stargateGraphqlHost, true, &stargateRestHost)

		})
	})
})

func verifyIngressRules(ingress networking.Ingress, authEnabled bool, authHost *string, graphEnabled bool, playgroundEnabled bool, graphHost *string, restEnabled bool, restHost *string) {
	rules := ingress.Spec.Rules
	serviceName := Sprintf("%s-%s-stargate-service", HelmReleaseName, "dc1")
	if authEnabled {
		kubeapi.VerifyIngressRule(rules, "/v1/auth", authHost, serviceName, 8081)
	} else {
		kubeapi.VerifyNoRuleWithPath(rules, "/v1/auth")
	}
	if graphEnabled {
		kubeapi.VerifyIngressRule(rules, "/graphql/", graphHost, serviceName, 8080)
		kubeapi.VerifyIngressRule(rules, "/graphql-schema", graphHost, serviceName, 8080)
		if playgroundEnabled {
			kubeapi.VerifyIngressRule(rules, "/playground", graphHost, serviceName, 8080)
		} else {
			kubeapi.VerifyNoRuleWithPath(rules, "/playground")
		}
	} else {
		kubeapi.VerifyNoRuleWithPath(rules, "/graphql/")
		kubeapi.VerifyNoRuleWithPath(rules, "/graphql-schema")
		kubeapi.VerifyNoRuleWithPath(rules, "/playground")
	}
	if restEnabled {
		kubeapi.VerifyIngressRule(rules, "/v2/", restHost, serviceName, 8082)
	} else {
		kubeapi.VerifyNoRuleWithPath(rules, "/v2/")
	}
}
