package unit_test

import (
	corev1 "k8s.io/api/core/v1"
	"path/filepath"
	"reflect"

	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"

	"github.com/gruntwork-io/terratest/modules/helm"
	api "github.com/k8ssandra/reaper-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify Reaper template", func() {
	var (
		helmChartPath string
		err           error
		reaper        *api.Reaper
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		reaper = &api.Reaper{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/reaper/reaper.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, reaper)
			})
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)

			Expect(string(reaper.Spec.ServerConfig.StorageType)).To(Equal("cassandra"))
			Expect(reaper.Kind).To(Equal("Reaper"))
		})

		It("changing datacenter name", func() {
			targetDcName := "reaper-dc"
			options := &helm.Options{
				SetStrValues: map[string]string{
					"cassandra.datacenters[0].name": targetDcName,
					"cassandra.datacenters[0].size": "1",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.CassandraBackend.CassandraDatacenter.Name).To(Equal(targetDcName))
		})

		It("modifying autoscheduling option", func() {
			options := &helm.Options{
				SetStrValues:   map[string]string{"reaper.autoschedule": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.AutoScheduling).ToNot(BeNil())
		})

		It("modifying autoscheduling additional properties", func() {
			options := &helm.Options{
				SetStrValues: map[string]string{
					"reaper.autoschedule":                               "true",
					"reaper.autoschedule_properties.initialDelayPeriod": "PT10S",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.AutoScheduling).ToNot(BeNil())
			Expect(reaper.Spec.ServerConfig.AutoScheduling.InitialDelay).To(Equal("PT10S"))
		})

		It("modifying secret options", func() {
			options := &helm.Options{
				SetStrValues: map[string]string{
					"reaper.jmx.secret":           "somethingelse",
					"reaper.cassandraUser.secret": "cassandraSpecial",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.CassandraBackend.CassandraUserSecretName).To(Equal("cassandraSpecial"))
			Expect(reaper.Spec.ServerConfig.JmxUserSecretName).To(Equal("somethingelse"))
		})

		It("verifying default secret values", func() {
			options := &helm.Options{
				SetStrValues: map[string]string{
					"cassandra.clusterName": "nowyouseeme",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.JmxUserSecretName).To(HavePrefix("nowyouseeme"))
		})

		It("using affinity", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/affinity-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			expected := &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/e2e-az-name",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"e2e-az1", "e2e-az2"},
									},
								},
							},
						},
					},
				},
			}

			affinity := reaper.Spec.Affinity
			Expect(affinity).To(Equal(expected))
		})

		It("using tolerations", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/tolerations-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			tolerations := reaper.Spec.Tolerations
			Expect(tolerations).To(ConsistOf(
				corev1.Toleration{
					Key:      "key1",
					Operator: corev1.TolerationOpEqual,
					Value:    "value1",
					Effect:   corev1.TaintEffectNoSchedule,
				},
				corev1.Toleration{
					Key:      "key2",
					Operator: corev1.TolerationOpEqual,
					Value:    "value2",
					Effect:   corev1.TaintEffectNoSchedule,
				},
			))
		})

		It("has default container security context", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			_ = renderTemplate(options)

			Expect(reaper.Spec.SecurityContext).ToNot(BeNil())
			Expect(*reaper.Spec.SecurityContext.ReadOnlyRootFilesystem).To(BeTrue())
		})

		It("has custom container security context", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/reaper-security-context-custom-values.yaml"},
			}
			_ = renderTemplate(options)

			Expect(reaper.Spec.SecurityContext).ToNot(BeNil())
			Expect(*reaper.Spec.SecurityContext.RunAsGroup).
				To(BeIdenticalTo(int64(8675309)))
		})

		It("has default init-container security context", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			_ = renderTemplate(options)

			Expect(reaper.Spec.SchemaInitContainerConfig).ToNot(BeNil())
			Expect(reaper.Spec.SchemaInitContainerConfig.SecurityContext).ToNot(BeNil())
			Expect(*reaper.Spec.SchemaInitContainerConfig.SecurityContext.ReadOnlyRootFilesystem).To(BeTrue())
		})

		It("has custom init-container security context", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/reaper-security-context-custom-values.yaml"},
			}
			_ = renderTemplate(options)

			Expect(reaper.Spec.SchemaInitContainerConfig).ToNot(BeNil())
			Expect(reaper.Spec.SchemaInitContainerConfig.SecurityContext).
				ToNot(BeNil())
			Expect(*reaper.Spec.SchemaInitContainerConfig.SecurityContext.RunAsGroup).
				To(BeIdenticalTo(int64(8675309)))
		})

		It("has default pod security context", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			_ = renderTemplate(options)

			Expect(reaper.Spec.PodSecurityContext).ToNot(BeNil())
			Expect(reflect.DeepEqual(*reaper.Spec.PodSecurityContext, corev1.PodSecurityContext{})).To(BeTrue())
		})

		It("has custom pod security context", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/reaper-security-context-custom-values.yaml"},
			}
			_ = renderTemplate(options)

			Expect(*reaper.Spec.PodSecurityContext.FSGroup).To(BeIdenticalTo(int64(1)))
		})
	})
})
