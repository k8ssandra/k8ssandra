package unit_test

import (
	. "fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	DefaultStargate3ImagePrefix    = "stargateio/stargate-3_11:"
	DefaultStargate4ImagePrefix    = "stargateio/stargate-4_0:"
	DefaultStargate3ClusterVersion = "3.11"
	DefaultStargate4ClusterVersion = "4.0"
	DefaultStargateImagePrefix     = DefaultStargate3ImagePrefix
	DefaultStargateClusterVersion  = DefaultStargate3ClusterVersion
)

var _ = Describe("Verify Stargate template", func() {
	var (
		helmChartPath string
		err           error
		deployment    *appsv1.Deployment
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		deployment = &appsv1.Deployment{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/stargate/stargate.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, deployment)
			})
	}

	Context("by confirming it does not render when", func() {
		It("is explicitly disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
		It("cassandra is explicitly disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by confirming it does render when", func() {
		It("is implicitly enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(renderTemplate(options)).Should(Succeed())
		})

		It("is explicitly enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "true",
				},
			}
			Expect(renderTemplate(options)).Should(Succeed())
		})
	})

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(deployment.Kind).To(Equal("Deployment"))

			templateSpec := deployment.Spec.Template.Spec
			Expect(len(templateSpec.InitContainers)).To(Equal(1))
			initContainer := templateSpec.InitContainers[0]
			Expect(string(initContainer.ImagePullPolicy)).To(Equal("IfNotPresent"))

			Expect(initContainer.Args[0]).To(Equal("-c"))
			Expect(initContainer.Args[1]).To(ContainSubstring(
				Sprintf("nslookup %s-dc1-service.%s.svc.cluster.local", HelmReleaseName, DefaultTestNamespace)))

			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(HavePrefix(DefaultStargateImagePrefix))
			Expect(container.Name).To(Equal(Sprintf("%s-dc1-stargate", HelmReleaseName)))
			Expect(string(container.ImagePullPolicy)).To(Equal("IfNotPresent"))

			oneMegabyte := 1024 * 1024
			limits := container.Resources.Limits
			Expect(limits.Memory().Value()).To(Equal(int64(1024 * oneMegabyte)))
			Expect(limits.Cpu().MilliValue()).To(Equal(int64(1000)))

			requests := container.Resources.Requests
			Expect(requests.Memory().Value()).To(Equal(int64(512 * oneMegabyte)))
			Expect(requests.Cpu().MilliValue()).To(Equal(int64(200)))

			javaOpts := kubeapi.FindEnvVarByName(container, "JAVA_OPTS")
			Expect(javaOpts.Value).To(ContainSubstring("-Xms256M"))
			Expect(javaOpts.Value).To(ContainSubstring("-Xmx256M"))

			clusterName := kubeapi.FindEnvVarByName(container, "CLUSTER_NAME")
			Expect(clusterName.Value).To(Equal(HelmReleaseName))

			clusterVersion := kubeapi.FindEnvVarByName(container, "CLUSTER_VERSION")
			Expect(clusterVersion.Value).To(Equal(DefaultStargateClusterVersion))

			seed := kubeapi.FindEnvVarByName(container, "SEED")
			Expect(seed.Value).To(Equal(Sprintf("%s-seed-service.%s.svc.cluster.local", HelmReleaseName, DefaultTestNamespace)))

			datacenterName := kubeapi.FindEnvVarByName(container, "DATACENTER_NAME")
			Expect(datacenterName.Value).To(Equal("dc1"))
		})

		It("using custom image and clusterVersion", func() {
			// This combination of values makes no real sense and would not work
			// but this tests that the defaults are avoided when a specific value is provided
			repo := "stargateio/stargate-4_0"
			tag := "v1.0.5"
			clusterVersion := "3.0"

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":          "true",
					"stargate.image.repository": repo,
					"stargate.image.tag":        tag,
					"stargate.clusterVersion":   clusterVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			templateSpec := deployment.Spec.Template.Spec
			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(Equal(DefaultRegistry + "/" + repo + ":" + tag))
			clusterVersionEnv := kubeapi.FindEnvVarByName(container, "CLUSTER_VERSION")
			Expect(clusterVersionEnv.Value).To(Equal(clusterVersion))
		})

		It("using defaults and empty clusterVersion", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":        "true",
					"stargate.clusterVersion": "",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			templateSpec := deployment.Spec.Template.Spec
			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(HavePrefix(DefaultStargateImagePrefix))
			clusterVersionEnv := kubeapi.FindEnvVarByName(container, "CLUSTER_VERSION")
			Expect(clusterVersionEnv.Value).To(Equal(DefaultStargateClusterVersion))
		})

		It("using cassandra version 4.0.0", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":  "true",
					"cassandra.version": "4.0.0",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			templateSpec := deployment.Spec.Template.Spec
			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(HavePrefix(DefaultStargate4ImagePrefix))
			clusterVersionEnv := kubeapi.FindEnvVarByName(container, "CLUSTER_VERSION")
			Expect(clusterVersionEnv.Value).To(Equal(DefaultStargate4ClusterVersion))
		})

		It("using cassandra version 3.11.10", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":  "true",
					"cassandra.version": "3.11.10",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			templateSpec := deployment.Spec.Template.Spec
			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(HavePrefix(DefaultStargate3ImagePrefix))
			clusterVersionEnv := kubeapi.FindEnvVarByName(container, "CLUSTER_VERSION")
			Expect(clusterVersionEnv.Value).To(Equal(DefaultStargate3ClusterVersion))
		})

		It("changing cluster name", func() {
			clusterName := Sprintf("k8ssandracluster%s", UniqueIdSuffix)
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName": clusterName,
					"stargate.enabled":      "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(deployment.Kind).To(Equal("Deployment"))

			initContainer := deployment.Spec.Template.Spec.InitContainers[0]
			Expect(initContainer.Args[0]).To(Equal("-c"))
			Expect(initContainer.Args[1]).To(ContainSubstring(
				Sprintf("nslookup %s-dc1-service.%s.svc.cluster.local", clusterName, DefaultTestNamespace)))

			container := deployment.Spec.Template.Spec.Containers[0]
			seed := kubeapi.FindEnvVarByName(container, "SEED")
			Expect(seed.Value).To(Equal(
				Sprintf("%s-seed-service.%s.svc.cluster.local", clusterName, DefaultTestNamespace)))
		})

		It("changing datacenter name", func() {
			targetDcName := "testDataCenter"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "true",
				},
				SetStrValues: map[string]string{
					"cassandra.datacenters[0].name": targetDcName,
					"cassandra.datacenters[0].size": "1",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			datacenterName := kubeapi.FindEnvVarByName(container, "DATACENTER_NAME")
			Expect(datacenterName.Value).To(Equal(targetDcName))

			clusterName := kubeapi.FindEnvVarByName(container, "CLUSTER_NAME")
			initContainer := deployment.Spec.Template.Spec.InitContainers[0]
			Expect(initContainer.Args[0]).To(Equal("-c"))
			Expect(initContainer.Args[1]).To(ContainSubstring(
				Sprintf("nslookup %s-%s-service.%s.svc.cluster.local",
					clusterName.Value, targetDcName, DefaultTestNamespace)))
		})

		It("changing memory allocation", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "true",
					"stargate.heapMB":  "512",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			oneGigabyte := 1024 * 1024 * 1024
			limits := container.Resources.Limits
			Expect(limits.Memory().Value()).To(Equal(int64(2 * oneGigabyte)))

			requests := container.Resources.Requests
			Expect(requests.Memory().Value()).To(Equal(int64(oneGigabyte)))

			javaOpts := kubeapi.FindEnvVarByName(container, "JAVA_OPTS")
			Expect(javaOpts.Value).To(ContainSubstring("-Xms512M"))
			Expect(javaOpts.Value).To(ContainSubstring("-Xmx512M"))
		})

		It("changing CPU allocation", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":          "true",
					"stargate.cpuReqMillicores": "2000",
					"stargate.cpuLimMillicores": "1500",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			limits := container.Resources.Limits
			Expect(limits.Cpu().MilliValue()).To(Equal(int64(1500)))

			requests := container.Resources.Requests
			Expect(requests.Cpu().MilliValue()).To(Equal(int64(2000)))
		})

		It("changing container image and imagePullPolicy", func() {
			repository := "stargateio/stargate-3_11"
			tag := "v1.0.3"
			alternatePullPolicy := "Always"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":          "true",
					"stargate.image.repository": repository,
					"stargate.image.tag":        tag,
					"stargate.image.pullPolicy": alternatePullPolicy,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Image).To(Equal(DefaultRegistry + "/" + repository + ":" + tag))
			Expect(string(container.ImagePullPolicy)).To(Equal(alternatePullPolicy))
		})

		It("changing stargate version", func() {
			alternateImage := "stargateio/stargate-3_11:v1.0.5"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled": "true",
					"stargate.version": "1.0.5",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Image).To(Equal(alternateImage))
		})

		It("changing stargate version with Cassandra 4.0", func() {
			alternateImage := "stargateio/stargate-4_0:v1.0.6"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":  "true",
					"stargate.version":  "1.0.6",
					"cassandra.version": "4.0.0",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Image).To(Equal(alternateImage))
		})

		It("changing probe initial delay", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":                      "true",
					"stargate.livenessInitialDelaySeconds":  "60",
					"stargate.readinessInitialDelaySeconds": "90",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.LivenessProbe.InitialDelaySeconds).To(Equal(int32(60)))
			Expect(container.ReadinessProbe.InitialDelaySeconds).To(Equal(int32(90)))
		})

		It("using tolerations", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/tolerations-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			tolerations := deployment.Spec.Template.Spec.Tolerations
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

			affinity := deployment.Spec.Template.Spec.Affinity
			Expect(affinity).To(Equal(expected))
		})
	})
})
