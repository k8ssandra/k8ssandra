package unit_test

import (
	. "fmt"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	DefaultStargate3Image          = "stargateio/stargate-3_11:v1.0.9"
	DefaultStargate4Image          = "stargateio/stargate-4_0:v1.0.9"
	DefaultStargate3ClusterVersion = "3.11"
	DefaultStargate4ClusterVersion = "4.0"
	DefaultStargateImage           = DefaultStargate3Image
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
			Expect(initContainer.Args[1]).To(ContainSubstring(Sprintf("nslookup %s-seed-service.%s.svc.cluster.local;", HelmReleaseName, DefaultTestNamespace)))

			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(Equal(DefaultStargateImage))
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
			image := "stargateio/stargate-4_0:v1.0.5"
			clusterVersion := "3.0"

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":        "true",
					"stargate.image":          image,
					"stargate.clusterVersion": clusterVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			templateSpec := deployment.Spec.Template.Spec
			Expect(len(templateSpec.Containers)).To(Equal(1))
			container := templateSpec.Containers[0]
			Expect(container.Image).To(Equal(image))
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
			Expect(container.Image).To(Equal(DefaultStargateImage))
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
			Expect(container.Image).To(Equal(DefaultStargate4Image))
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
			Expect(container.Image).To(Equal(DefaultStargate3Image))
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
			Expect(initContainer.Args[1]).To(ContainSubstring(Sprintf("nslookup %s-seed-service.%s.svc.cluster.local;", clusterName, DefaultTestNamespace)))

			container := deployment.Spec.Template.Spec.Containers[0]
			seed := kubeapi.FindEnvVarByName(container, "SEED")
			Expect(seed.Value).To(Equal(Sprintf("%s-seed-service.%s.svc.cluster.local", clusterName, DefaultTestNamespace)))
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
			alternateImage := "stargateio/stargate-3_11:v1.0.3"
			alternatePullPolicy := "Always"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":         "true",
					"stargate.image":           alternateImage,
					"stargate.imagePullPolicy": alternatePullPolicy,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Image).To(Equal(alternateImage))
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
	})
})
