package unit_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"

	"github.com/gruntwork-io/terratest/modules/helm"
	cassdcv1beta1 "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/cassdc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type CassandraConfig struct {
	Authenticator             string
	Authorizer                string
	RolesValidityMillis       int64 `json:"roles_validity_in_ms"`
	RolesUpdateMillis         int64 `json:"roles_update_interval_in_ms"`
	PermissionsValidityMillis int64 `json:"permissions_validity_in_ms"`
	PermissionsUpdateMillis   int64 `json:"permissions_update_interval_in_ms"`
	CredentialsValidityMillis int64 `json:"credentials_validity_in_ms"`
	CredentialsUpdateMillis   int64 `json:"credentials_update_interval_in_ms"`
	NumTokens                 int64 `json:"num_tokens"`
	//AllocateTokensForLocalRF  int64 `json:"allocate_tokens_for_local_replication_factor"`
	AllocateTokensForLocalRF int64 `json:"allocate_tokens_for_local_replication_factor"`
}

type JvmOptions struct {
	AdditionalJvmOptions           []string `json:"additional-jvm-opts"`
	InitialHeapSize                string   `json:"initial_heap_size"`
	MaxHeapSize                    string   `json:"max_heap_size"`
	YoungGenSize                   string   `json:"heap_size_young_generation"`
	GarbageCollector               string   `json:"garbage_collector"`
	SurvivorRatio                  *int64   `json:"survivor_ratio"`
	MaxTenuringThreshold           *int64   `json:"max_tenuring_threshold"`
	InitiatingOccupancyFraction    *int64   `json:"cms_initiating_occupancy_fraction"`
	CmsWaitDuration                *int64   `json:"cms_wait_duration"`
	SetUpdatingPauseTimePercent    *int64   `json:"g1r_set_updating_pause_time_percent"`
	MaxGcPauseMillis               *int64   `json:"max_gc_pause_millis"`
	InitiatingHeapOccupancyPercent *int64   `json:"initiating_heap_occupancy_percent"`
	ParallelGcThreads              *int64   `json:"parallel_gc_threads"`
	ConcurrentGcThreads            *int64   `json:"conc_gc_threads"`
}

type Config struct {
	CassandraConfig  CassandraConfig `json:"cassandra-yaml"`
	JvmOptions       *JvmOptions     `json:"jvm-options"`
	JvmServerOptions *JvmOptions     `json:"jvm-server-options"`
}

var (
	reaperInstanceValue    = fmt.Sprintf("%s-reaper", HelmReleaseName)
	medusaConfigVolumeName = fmt.Sprintf("%s-medusa", HelmReleaseName)
)

const (
	ConfigInitContainer         = "server-config-init"
	BaseConfigInitContainer     = "base-config-init"
	MedusaInitContainer         = "medusa-restore"
	JmxCredentialsInitContainer = "jmx-credentials"

	CassandraContainer = "cassandra"
	MedusaContainer    = "medusa"

	CassandraConfigVolumeName            = "cassandra-config"
	CassandraMetricsCollConfigVolumeName = "cassandra-metrics-coll-config"
	CassandraTmpVolumeName               = "cassandra-tmp"

	MedusaBucketKeyVolumeName = "medusa-bucket-key"
	PodInfoVolumeName         = "podinfo"
)

var _ = Describe("Verify CassandraDatacenter template", func() {
	var (
		helmChartPath string
		cassdc        *cassdcv1beta1.CassandraDatacenter
	)

	BeforeEach(func() {
		path, err := filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		helmChartPath = path
		cassdc = &cassdcv1beta1.CassandraDatacenter{}
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/cassandra/cassdc.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, cassdc)
			})
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Kind).To(Equal("CassandraDatacenter"))

			// Reaper should be enabled in default - verify
			// Verify reaper annotation is set
			Expect(cassdc.Annotations).Should(HaveKeyWithValue(ReaperInstanceAnnotation, reaperInstanceValue))

			initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
			Expect(len(initContainers)).To(Equal(3))
			// Verify initContainers includes the base Cassandra config
			Expect(initContainers[0].Name).To(Equal(BaseConfigInitContainer))
			// Verify initContainers includes config-builder config
			Expect(initContainers[1].Name).To(Equal(ConfigInitContainer))
			// Verify initContainers includes JMX credentials
			Expect(initContainers[2].Name).To(Equal(JmxCredentialsInitContainer))
			// Verify LOCAL_JMX value
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(1))
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Name).To(Equal("LOCAL_JMX"))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Value).To(Equal("no"))
			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(false))
			Expect(*cassdc.Spec.DockerImageRunsAsCassandra).To(BeTrue())

			// Server version and mgmt-api image specified
			Expect(cassdc.Spec.ServerVersion).ToNot(BeEmpty())
			Expect(cassdc.Spec.ServerImage).ToNot(BeEmpty())

			// JVM heap options -- default to settings as defined in cassdc.yaml
			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmServerOptions).ToNot(BeNil())
			Expect(config.JvmServerOptions.InitialHeapSize).To(BeEmpty())
			Expect(config.JvmServerOptions.MaxHeapSize).To(BeEmpty())
			Expect(config.JvmServerOptions.YoungGenSize).To(BeEmpty())

			// Default set of volume and volume mounts
			Expect(kubeapi.GetVolumeMountNames(&initContainers[0])).To(ConsistOf(CassandraConfigVolumeName,
				CassandraMetricsCollConfigVolumeName, CassandraTmpVolumeName))
			Expect(kubeapi.GetVolumeNames(cassdc.Spec.PodTemplateSpec)).To(ConsistOf(CassandraConfigVolumeName,
				CassandraMetricsCollConfigVolumeName, CassandraTmpVolumeName))

			// Default security context for containers
			AssertContainerSecurityContextExists(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			AssertContainerSecurityContextNotExists(cassdc, MedusaContainer, MedusaInitContainer)

			// TODO - this will change once mgmt-api root access needs are addressed.
			isReadOnlyRootFilesystemAllowed := false
			expectedCtx := corev1.SecurityContext{ReadOnlyRootFilesystem: &isReadOnlyRootFilesystemAllowed}
			AssertContainerSecurityContextExistsAndMatches(cassdc, CassandraContainer, expectedCtx)

			// Default pod security context for cassdc
			Expect(cassdc.Spec.PodTemplateSpec.Spec.SecurityContext).ToNot(BeNil())
		})

		It("is not rendered if disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.enabled": "false",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})

		It("override clusterName", func() {
			clusterName := "test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName": clusterName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ClusterName).To(Equal(clusterName))
		})

		It("default clusterName as release name", func() {
			clusterName := ""
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName": clusterName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ClusterName).To(Equal(HelmReleaseName))
		})

		It("override datacenter name", func() {
			dcName := "test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.datacenters[0].name": dcName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Name).To(Equal(dcName))
		})

		It("override datacenter size and name", func() {
			dcName := "dc1"
			size := "3"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.datacenters[0].size": size,
					// Not sure why, but if we do not specify the name here we get a
					// template rendering error in reaper.yaml.
					"cassandra.datacenters[0].name": dcName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Spec.Size, 3)
		})

		It("disabling the logging sidecar", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.loggingSidecar.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			AssertContainerNamesMatch(cassdc, CassandraContainer)
		})

		It("using private registry and non-default images", func() {
			registry := "localhost:5000"
			configBuilderRepo := "test/config-builder"
			configBuilderTag := "5.0"
			systemLoggerRepo := "test/system-logger"
			systemLoggerTag := "1.0"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.configBuilder.image.registry":    registry,
					"cassandra.configBuilder.image.repository":  configBuilderRepo,
					"cassandra.configBuilder.image.tag":         configBuilderTag,
					"cassandra.loggingSidecar.image.registry":   registry,
					"cassandra.loggingSidecar.image.repository": systemLoggerRepo,
					"cassandra.loggingSidecar.image.tag":        systemLoggerTag,
					"cassandra.serviceAccount":                  "k8ssandra",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Spec.ConfigBuilderImage).To(Equal("localhost:5000/test/config-builder:5.0"))
			Expect(cassdc.Spec.SystemLoggerImage).To(Equal("localhost:5000/test/system-logger:1.0"))
			Expect(cassdc.Spec.ServiceAccount).To(Equal("k8ssandra"))
		})

		It("using custom init containers", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/init-containers-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())
			initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
			Expect(len(initContainers)).To(Equal(5))
			Expect(initContainers[3].Name).To(Equal("foo"))
			Expect(initContainers[4].Name).To(Equal("bar"))
		})

		It("using multiple racks with no affinity labels", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/racks-no-affinity-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Spec.Racks).To(ConsistOf([]cassdcv1beta1.Rack{
				{
					Name: "r1",
				},
				{
					Name: "r2",
				},
				{
					Name: "r3",
				},
			}))
		})

		It("using multiple racks with affinity labels", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/racks-affinity-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Spec.Racks).To(ConsistOf([]cassdcv1beta1.Rack{
				{
					Name: "r1",
					NodeAffinityLabels: map[string]string{
						"topology.kubernetes.io/zone": "us-east1-b",
					},
				},
				{
					Name: "r2",
					NodeAffinityLabels: map[string]string{
						"topology.kubernetes.io/zone": "us-east1-a",
					},
				},
				{
					Name: "r3",
					NodeAffinityLabels: map[string]string{
						"topology.kubernetes.io/zone": "us-east1-c",
					},
				},
			}))
		})

		It("disabling reaper and medusa and stargate", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"stargate.enabled": "false",
					"reaper.enabled":   "false",
					"medusa.enabled":   "false",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Annotations).ShouldNot(HaveKeyWithValue(ReaperInstanceAnnotation, reaperInstanceValue))

			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(1))
			// No env slice should be present
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env).To(BeNil())

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			AssertContainerSecurityContextExists(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			// TODO - this will change once mgmt-api root access needs are addressed.
			isReadOnlyRootFilesystemAllowed := false
			expectedCtx := corev1.SecurityContext{ReadOnlyRootFilesystem: &isReadOnlyRootFilesystemAllowed}
			AssertContainerSecurityContextExistsAndMatches(cassdc, CassandraContainer, expectedCtx)

			AssertContainerSecurityContextNotExists(cassdc, MedusaContainer, MedusaInitContainer)

			// No users should exist
			Expect(cassdc.Spec.Users).To(BeNil())
		})

		It("enabling only medusa using custom securityContext", func() {
			storageSecret := HelmReleaseName + "-medusa-storage"
			options := &helm.Options{
				SetValues: map[string]string{
					"medusa.enabled":       "true",
					"medusa.storageSecret": storageSecret,
					"reaper.enabled":       "false",
				},
				ValuesFiles:    []string{"./testdata/medusa-security-context-custom-values.yaml"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer, JmxCredentialsInitContainer, MedusaInitContainer)
			AssertContainerSecurityContextExists(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			// Two containers, medusa and cassandra
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(2))
			// Second container should be medusa
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].Name).To(Equal(MedusaContainer))

			medusaContainer := GetContainer(cassdc, MedusaContainer)
			medusaConfigMap := HelmReleaseName + "-medusa"
			medusaRestoreInitContainer := GetInitContainer(cassdc, MedusaInitContainer)

			Expect(kubeapi.GetVolumeMountNames(medusaContainer)).To(ConsistOf(medusaConfigMap,
				CassandraConfigVolumeName, "server-data", storageSecret))
			Expect(kubeapi.GetVolumeNames(cassdc.Spec.PodTemplateSpec)).To(ConsistOf(medusaConfigMap,
				CassandraConfigVolumeName, CassandraMetricsCollConfigVolumeName, CassandraTmpVolumeName,
				storageSecret, PodInfoVolumeName))

			Expect(medusaRestoreInitContainer).ToNot(BeNil())
			Expect(medusaRestoreInitContainer.SecurityContext).ToNot(BeNil())

			// TODO - this will change once medusa root access needs are addressed.
			Expect(*medusaRestoreInitContainer.SecurityContext.ReadOnlyRootFilesystem).To(BeFalse())

			Expect(medusaContainer.SecurityContext).ToNot(BeNil())
			// TODO - this will change once medusa root access needs are addressed.
			Expect(*medusaContainer.SecurityContext.ReadOnlyRootFilesystem).To(BeFalse())

		})

		It("enabling only medusa with local storage", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"medusa.enabled":                 "true",
					"medusa.storage":                 "local",
					"medusa.podStorage.size":         "30Gi",
					"medusa.podStorage.storageClass": "slow",
					"reaper.enabled":                 "false",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer, MedusaInitContainer)
			AssertContainerSecurityContextExists(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			// TODO - this will change once medusa root access needs are addressed.
			isReadOnlyRootFilesystemAllowed := false
			expectedCtx := corev1.SecurityContext{ReadOnlyRootFilesystem: &isReadOnlyRootFilesystemAllowed}
			AssertContainerSecurityContextExistsAndMatches(cassdc, MedusaInitContainer, expectedCtx)
			AssertContainerSecurityContextExistsAndMatches(cassdc, MedusaContainer, expectedCtx)

			// Two containers, medusa and cassandra
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(2))
			// Second container should be medusa
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].Name).To(Equal(MedusaContainer))
			// AdditionalVolumes should have been created
			Expect(cassdc.Spec.StorageConfig.AdditionalVolumes).To(HaveLen(1))
			Expect(*cassdc.Spec.StorageConfig.AdditionalVolumes[0].PVCSpec.StorageClassName).To(Equal("slow"))
			Expect(len(cassdc.Spec.StorageConfig.AdditionalVolumes[0].PVCSpec.AccessModes)).To(Equal(1))
			Expect(cassdc.Spec.StorageConfig.AdditionalVolumes[0].PVCSpec.AccessModes[0]).To(Equal(corev1.ReadWriteOnce))
			Expect(cassdc.Spec.StorageConfig.AdditionalVolumes[0].Name).To(Equal("medusa-backups"))

			medusaContainer := GetContainer(cassdc, MedusaContainer)
			medusaConfigMap := HelmReleaseName + "-medusa"

			Expect(kubeapi.GetVolumeMountNames(medusaContainer)).To(ConsistOf(medusaConfigMap, CassandraConfigVolumeName,
				"server-data", "medusa-backups"))
			Expect(kubeapi.GetVolumeNames(cassdc.Spec.PodTemplateSpec)).To(ConsistOf(medusaConfigMap,
				CassandraConfigVolumeName, CassandraMetricsCollConfigVolumeName, CassandraTmpVolumeName, PodInfoVolumeName))

			medusaRestoreInitContainer := GetInitContainer(cassdc, MedusaInitContainer)

			Expect(kubeapi.GetVolumeMountNames(medusaRestoreInitContainer)).To(ConsistOf(medusaConfigMap,
				"server-config", "server-data", "medusa-backups", PodInfoVolumeName))
		})

		It("enabling only medusa with local storage with modified access modes", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"medusa.enabled":                   "true",
					"medusa.storage":                   "local",
					"medusa.podStorage.size":           "30Gi",
					"medusa.podStorage.storageClass":   "slow",
					"medusa.podStorage.accessModes[0]": "ReadWriteMany",
					"reaper.enabled":                   "false",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer, MedusaInitContainer)
			AssertContainerSecurityContextExists(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			// TODO - this will change once mgmt-api root access needs are addressed.
			isReadOnlyRootFilesystemAllowed := false
			expectedCtx := corev1.SecurityContext{ReadOnlyRootFilesystem: &isReadOnlyRootFilesystemAllowed}
			AssertContainerSecurityContextExistsAndMatches(cassdc, CassandraContainer, expectedCtx)
			// TODO - this will change once medusa root access needs are addressed.
			AssertContainerSecurityContextExistsAndMatches(cassdc, MedusaInitContainer, expectedCtx)
			AssertContainerSecurityContextExistsAndMatches(cassdc, MedusaContainer, expectedCtx)

			// Two containers, medusa and cassandra
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(2))
			// Second container should be medusa
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].Name).To(Equal(MedusaContainer))
			// AdditionalVolumes should have been created
			Expect(cassdc.Spec.StorageConfig.AdditionalVolumes).To(HaveLen(1))
			Expect(*cassdc.Spec.StorageConfig.AdditionalVolumes[0].PVCSpec.StorageClassName).To(Equal("slow"))
			Expect(len(cassdc.Spec.StorageConfig.AdditionalVolumes[0].PVCSpec.AccessModes)).To(Equal(1))
			Expect(cassdc.Spec.StorageConfig.AdditionalVolumes[0].PVCSpec.AccessModes[0]).To(Equal(corev1.ReadWriteMany))
			Expect(cassdc.Spec.StorageConfig.AdditionalVolumes[0].Name).To(Equal("medusa-backups"))
		})

		It("enabling only medusa with local storage but missing size and storageclass", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"medusa.enabled": "true",
					"medusa.storage": "local",
					"reaper.enabled": "false",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Not(Succeed()))
		})

		It("enabling reaper and medusa", func() {
			// Simple verification that both have properties correctly applied
			options := &helm.Options{
				SetValues: map[string]string{"medusa.enabled": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer, MedusaInitContainer)

			AssertContainerNamesMatch(cassdc, CassandraContainer, MedusaContainer)

			AssertContainerSecurityContextExists(cassdc, BaseConfigInitContainer, ConfigInitContainer,
				JmxCredentialsInitContainer)

			// TODO - this will change once mgmt-api root access needs are addressed.
			isReadOnlyRootFilesystemAllowed := false
			expectedCtx := corev1.SecurityContext{ReadOnlyRootFilesystem: &isReadOnlyRootFilesystemAllowed}
			AssertContainerSecurityContextExistsAndMatches(cassdc, CassandraContainer, expectedCtx)

			// TODO - this will change once medusa root access needs are addressed.
			AssertContainerSecurityContextExistsAndMatches(cassdc, MedusaContainer, expectedCtx)
			AssertContainerSecurityContextExistsAndMatches(cassdc, MedusaInitContainer, expectedCtx)

		})

		It("adding additionalSeeds", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.additionalSeeds[0]": "127.0.0.1",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.AdditionalSeeds).To(HaveLen(1))
		})

		It("setting allowMultipleNodesPerWorker to true", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.allowMultipleNodesPerWorker": "true",
					"cassandra.resources.limits.memory":     "2Gi",
					"cassandra.resources.limits.cpu":        "1",
					"cassandra.resources.requests.memory":   "2Gi",
					"cassandra.resources.requests.cpu":      "1"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(true))
		})

		It("setting allowMultipleNodesPerWorker to false", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.allowMultipleNodesPerWorker": "false",
					"cassandra.resources.limits.memory":     "2Gi",
					"cassandra.resources.limits.cpu":        "1",
					"cassandra.resources.requests.memory":   "2Gi",
					"cassandra.resources.requests.cpu":      "1",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(false))
			Expect(*cassdc.Spec.Resources.Limits.Memory()).To(Equal(resource.MustParse("2Gi")))
			Expect(*cassdc.Spec.Resources.Limits.Cpu()).To(Equal(resource.MustParse("1")))
			Expect(*cassdc.Spec.Resources.Requests.Memory()).To(Equal(resource.MustParse("2Gi")))
			Expect(*cassdc.Spec.Resources.Requests.Cpu()).To(Equal(resource.MustParse("1")))
		})

		It("setting allowMultipleNodesPerWorker to false without resources", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.allowMultipleNodesPerWorker": "false",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(false))
		})

		It("setting allowMultipleNodesPerWorker to true without resources", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.allowMultipleNodesPerWorker": "true",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderedErr := renderTemplate(options)

			Expect(renderedErr).ToNot(BeNil())
			Expect(renderedErr.Error()).To(ContainSubstring("set resource limits/requests when enabling allowMultipleNodesPerWorker"))

		})

		It("with the configOverride", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/config-override-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.NumTokens).To(Equal(int64(16)))
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("800m"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("800m"))
			Expect(config.JvmOptions.AdditionalJvmOptions).To(ConsistOf(
				"-Dcassandra.test=true",
				"-Dcassandra.k8ssandra=true",
			))
		})

		It("using allocateTokensForLocalRF with Cassandra 3.11", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": "3.11.11",
					"cassandra.datacenters[0].allocateTokensForLocalRF": "3",
				},
			}

			err := renderTemplate(options)

			Expect(err).ToNot(BeNil(), "Rendering should fail when using allocateTokensForLocalRF with Cassandra 3.11")
		})

		It("using allocateTokensForLocalRF with Cassandra 4.0 and default value", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": "4.0.0",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.AllocateTokensForLocalRF).To(Equal(int64(3)))
		})

		It("using allocateTokensForLocalRF with Cassandra 4.0 and non-default value", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version":                                 "4.0.0",
					"cassandra.datacenters[0].name":                     "test",
					"cassandra.datacenters[0].allocateTokensForLocalRF": "5",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.AllocateTokensForLocalRF).To(Equal(int64(5)))
		})

		It("using tolerations", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				ValuesFiles:    []string{"./testdata/tolerations-values.yaml"},
			}

			Expect(renderTemplate(options)).To(Succeed())

			tolerations := cassdc.Spec.Tolerations
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
	})

	Context("when configuring the JVM heap for Cassandra 3.11", func() {
		It("at cluster-level only", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version":             "3.11.11",
					"cassandra.heap.size":           "700M",
					"cassandra.heap.newGenSize":     "350M",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())

			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("700M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("700M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("350M"))
			Expect(config.JvmServerOptions).To(BeNil())
		})

		// Note: currently only one DC supported, to be expanded in future release.
		It("at dc-level overriding cluster level", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version":                        "3.11.11",
					"cassandra.heap.size":                      "700M",
					"cassandra.heap.newGenSize":                "350M",
					"cassandra.datacenters[0].heap.size":       "300M",
					"cassandra.datacenters[0].heap.newGenSize": "150M",
					"cassandra.datacenters[0].name":            dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("300M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("300M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("150M"))
			Expect(config.JvmServerOptions).To(BeNil())
		})

		It("at dc-level without newGenSize", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version":                  "3.11.11",
					"cassandra.datacenters[0].heap.size": "300M",
					"cassandra.datacenters[0].name":      dcName,
					// Note: not setting - "cassandra.datacenters[0].heap.newGenSize": "150M",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("300M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("300M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal(""))
			Expect(config.JvmServerOptions).To(BeNil())
		})

		It("at dc-level without size", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version": "3.11.11",
					// Note: not setting "cassandra.datacenters[0].heap.size":       "300M",
					"cassandra.datacenters[0].name":            dcName,
					"cassandra.datacenters[0].heap.newGenSize": "150M",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal(""))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal(""))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("150M"))
			Expect(config.JvmServerOptions).To(BeNil())
		})

		It("at cluster-level without newGenSize", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version":             "3.11.11",
					"cassandra.heap.size":           "300M",
					"cassandra.datacenters[0].name": dcName,
					// Note: not setting - "cassandra.heap.newGenSize": "150M",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())

			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("300M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("300M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal(""))
			Expect(config.JvmServerOptions).To(BeNil())
		})

		It("at cluster-level without size", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version": "3.11.11",
					// Note: not setting - "cassandra.heap.size": "300M",
					"cassandra.heap.newGenSize":     "150M",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())

			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal(""))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal(""))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("150M"))
			Expect(config.JvmServerOptions).To(BeNil())
		})
	})

	Context("when specifying JVM heap size formatted values at dc or cluster levels", func() {
		It("detected an invalid IEC formatted value at dc-level heap size", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":                      "700M",
					"cassandra.heap.newGenSize":                "350M",
					"cassandra.datacenters[0].heap.newGenSize": "150M",
					"cassandra.datacenters[0].heap.size":       "300MiB",
					"cassandra.datacenters[0].name":            dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Specify datacenter.heap.size using one of these suffixes"))

		})

		It("detected an invalid IEC formatted value at dc-level heap newGenSize", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":                      "700M",
					"cassandra.heap.newGenSize":                "350M",
					"cassandra.datacenters[0].heap.newGenSize": "150MiB",
					"cassandra.datacenters[0].heap.size":       "300M",
					"cassandra.datacenters[0].name":            dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Specify datacenter.heap.newGenSize using one of these suffixes"))
		})

		It("detected an invalid IEC formatted value at cluster heap size", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":           "700MiB",
					"cassandra.heap.newGenSize":     "350M",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Specify cassandra.heap.size using one of these suffixes"))
		})

		It("detected an invalid IEC formatted value at cluster heap newGenSize", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":           "700M",
					"cassandra.heap.newGenSize":     "350MiB",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Specify cassandra.heap.newGenSize using one of these suffixes"))
		})

		It("detected an invalid formatted value at cluster heap newGenSize", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size": "700M",
					// No spaces allowed in this format.
					"cassandra.heap.newGenSize":     "9000 k",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Specify cassandra.heap.newGenSize using one of these suffixes"))
		})

		It("detected an invalid formatted value at cluster heap size", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size": "700X",
					// No spaces allowed in this format.
					"cassandra.heap.newGenSize":     "9000k",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Specify cassandra.heap.size using one of these suffixes"))
		})

		It("detected a valid formatted value as decimal at cluster heap newGenSize", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":           "700M",
					"cassandra.heap.newGenSize":     "10240.5k",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).To(BeNil())
		})

		It("detected a valid formatted value as decimal at cluster heap size", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":           "4.5G",
					"cassandra.heap.newGenSize":     "1.5G",
					"cassandra.datacenters[0].name": dcName,
					// helmUtils.OPT_RENDER_TEMPLATE:   `{ "dir":"/tmp/foobar","name":"cluster_heap_size_valid.yaml"}`,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).To(BeNil())
		})

		It("detected a valid formatted value as decimal at dc heap size", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":                      "700M",
					"cassandra.heap.newGenSize":                "350M",
					"cassandra.datacenters[0].heap.newGenSize": "1T",
					"cassandra.datacenters[0].heap.size":       "4.5T",
					"cassandra.datacenters[0].name":            dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).To(BeNil())
		})

		It("detected a valid formatted value as decimal at dc heap newGenSize", func() {
			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":                      "700M",
					"cassandra.heap.newGenSize":                "350M",
					"cassandra.datacenters[0].heap.newGenSize": "1.202T",
					"cassandra.datacenters[0].heap.size":       "4.5T",
					"cassandra.datacenters[0].name":            dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}
			Expect(options).ToNot(BeNil())
			err := renderTemplate(options)

			Expect(err).To(BeNil())
		})
	})

	Context("when enabling auth", func() {
		It("with caches configured", func() {
			dcName := "test"
			clusterSize := 3
			clusterName := "auth-test"

			authCachePeriod := int64(7200000)
			cacheValidityPeriod := authCachePeriod + 1
			cacheUpdateInterval := authCachePeriod + 2

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":                    clusterName,
					"cassandra.datacenters[0].name":            dcName,
					"cassandra.datacenters[0].size":            strconv.Itoa(clusterSize),
					"cassandra.auth.enabled":                   "true",
					"cassandra.auth.cacheValidityPeriodMillis": strconv.FormatInt(cacheValidityPeriod, 10),
					"cassandra.auth.cacheUpdateIntervalMillis": strconv.FormatInt(cacheUpdateInterval, 10),
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Name).To(Equal(dcName))

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.Authenticator).To(Equal("PasswordAuthenticator"))
			Expect(config.CassandraConfig.Authorizer).To(Equal("CassandraAuthorizer"))
			Expect(config.CassandraConfig.RolesValidityMillis).To(Equal(cacheValidityPeriod))
			Expect(config.CassandraConfig.RolesUpdateMillis).To(Equal(cacheUpdateInterval))
			Expect(config.CassandraConfig.PermissionsValidityMillis).To(Equal(cacheValidityPeriod))
			Expect(config.CassandraConfig.PermissionsUpdateMillis).To(Equal(cacheUpdateInterval))
			Expect(config.CassandraConfig.CredentialsValidityMillis).To(Equal(cacheValidityPeriod))
			Expect(config.CassandraConfig.CredentialsUpdateMillis).To(Equal(cacheUpdateInterval))
			Expect(config.JvmServerOptions.AdditionalJvmOptions).To(ConsistOf(
				"-Dcassandra.system_distributed_replication_dc_names="+dcName,
				"-Dcassandra.system_distributed_replication_per_dc="+strconv.Itoa(clusterSize),
			))

			Expect(cassdc.Spec.Users).To(ConsistOf(
				cassdcv1beta1.CassandraUser{Superuser: true, SecretName: clusterName + "-reaper"},
				cassdcv1beta1.CassandraUser{Superuser: true, SecretName: clusterName + "-stargate"},
			))
		})

		It("with a superuser secret", func() {
			clusterName := "superuser-test"
			secretName := "test-secret"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":           clusterName,
					"cassandra.auth.enabled":          "true",
					"cassandra.auth.superuser.secret": secretName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.SuperuserSecretName).To(Equal(secretName))
		})

		It("with the superuser username", func() {
			clusterName := "superuser-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":             clusterName,
					"cassandra.auth.enabled":            "true",
					"cassandra.auth.superuser.username": "admin",
				},
			}

			err := renderTemplate(options)
			fmt.Println("error: ", err)
			Expect(err).To(Succeed())

			Expect(cassdc.Spec.SuperuserSecretName).To(Equal(clusterName + "-superuser"))
		})

		It("with a default secret for medusa", func() {
			clusterName := "medusa-user-test"
			secretName := clusterName + "-medusa"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.clusterName":  clusterName,
					"cassandra.auth.enabled": "true",
					"medusa.enabled":         "true",
					"reaper.enabled":         "false",
					"stargate.enabled":       "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.Users).To(ContainElement(cassdcv1beta1.CassandraUser{Superuser: true, SecretName: clusterName + "-medusa"}))

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer, JmxCredentialsInitContainer, MedusaInitContainer)

			initContainer := GetInitContainer(cassdc, "medusa-restore")
			Expect(initContainer).To(Not(BeNil()))

			cqlUsernameEnvVar := corev1.EnvVar{
				Name: "CQL_USERNAME",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: secretName,
						},
						Key: "username",
					},
				},
			}
			cqlPasswordEnvVar := corev1.EnvVar{
				Name: "CQL_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: secretName,
						},
						Key: "password",
					},
				},
			}

			Expect(initContainer.Env).To(ConsistOf([]corev1.EnvVar{
				{
					Name:  "MEDUSA_MODE",
					Value: "RESTORE",
				},
				cqlUsernameEnvVar,
				cqlPasswordEnvVar,
			}))

			AssertContainerNamesMatch(cassdc, CassandraContainer, MedusaContainer)

			cassandraContainer := GetContainer(cassdc, CassandraContainer)
			Expect(cassandraContainer).To(Not(BeNil()))

			medusaContainer := GetContainer(cassdc, MedusaContainer)
			Expect(medusaContainer).To(Not(BeNil()))

			Expect(medusaContainer.Env).To(ConsistOf([]corev1.EnvVar{
				{
					Name:  "MEDUSA_MODE",
					Value: "GRPC",
				},
				cqlUsernameEnvVar,
				cqlPasswordEnvVar,
			}))

			verifyMedusaVolumeMounts(medusaContainer)

			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Volumes)).To(Equal(6))
			AssertVolumeNamesMatch(cassdc, CassandraConfigVolumeName, CassandraMetricsCollConfigVolumeName,
				CassandraTmpVolumeName, medusaConfigVolumeName, MedusaBucketKeyVolumeName, PodInfoVolumeName)
			Expect(cassdc.Spec.Users).To(ContainElement(cassdcv1beta1.CassandraUser{SecretName: secretName, Superuser: true}))
		})

		It("with user-defined secret for medusa", func() {
			clusterName := "medusa-user-test"
			secretName := "medusa-user"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.clusterName":       clusterName,
					"cassandra.auth.enabled":      "true",
					"medusa.enabled":              "true",
					"medusa.cassandraUser.secret": secretName,
					"reaper.enabled":              "false",
					"stargate.enabled":            "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Spec.Users).To(ContainElement(cassdcv1beta1.CassandraUser{Superuser: true, SecretName: secretName}))

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer, JmxCredentialsInitContainer, MedusaInitContainer)

			initContainer := GetInitContainer(cassdc, MedusaInitContainer)
			Expect(initContainer).To(Not(BeNil()))

			cqlUsernameEnvVar := corev1.EnvVar{
				Name: "CQL_USERNAME",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: secretName,
						},
						Key: "username",
					},
				},
			}
			cqlPasswordEnvVar := corev1.EnvVar{
				Name: "CQL_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: secretName,
						},
						Key: "password",
					},
				},
			}

			Expect(initContainer.Env).To(ConsistOf([]corev1.EnvVar{
				{
					Name:  "MEDUSA_MODE",
					Value: "RESTORE",
				},
				cqlUsernameEnvVar,
				cqlPasswordEnvVar,
			}))

			AssertContainerNamesMatch(cassdc, CassandraContainer, MedusaContainer)

			cassandraContainer := GetContainer(cassdc, CassandraContainer)
			Expect(cassandraContainer).To(Not(BeNil()))

			medusaContainer := GetContainer(cassdc, MedusaContainer)
			Expect(medusaContainer).To(Not(BeNil()))

			Expect(medusaContainer.Env).To(ConsistOf([]corev1.EnvVar{
				{
					Name:  "MEDUSA_MODE",
					Value: "GRPC",
				},
				cqlUsernameEnvVar,
				cqlPasswordEnvVar,
			}))

			verifyMedusaVolumeMounts(medusaContainer)

			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Volumes)).To(Equal(6))

			AssertVolumeNamesMatch(cassdc, CassandraConfigVolumeName, CassandraTmpVolumeName,
				CassandraMetricsCollConfigVolumeName, medusaConfigVolumeName, MedusaBucketKeyVolumeName, PodInfoVolumeName)

			Expect(cassdc.Spec.Users).To(ContainElement(cassdcv1beta1.CassandraUser{SecretName: secretName, Superuser: true}))
		})

		It("with stargate enabled", func() {
			dcName := "test"
			clusterSize := 3
			clusterName := "auth-test"

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"stargate.enabled":              "true",
					"cassandra.clusterName":         clusterName,
					"medusa.enabled":                "false",
					"reaper.enabled":                "false",
					"cassandra.auth.enabled":        "true",
					"cassandra.datacenters[0].name": dcName,
					"cassandra.datacenters[0].size": strconv.Itoa(clusterSize),
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Name).To(Equal(dcName))

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.Authenticator).To(Equal("PasswordAuthenticator"))
			Expect(config.CassandraConfig.Authorizer).To(Equal("CassandraAuthorizer"))

			Expect(cassdc.Spec.Users).To(ConsistOf(cassdcv1beta1.CassandraUser{Superuser: true, SecretName: clusterName + "-stargate"}))
		})
	})

	Context("when disabling auth", func() {
		It("with reaper enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.auth.enabled": "false",
					"reaper.enabled":         "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.Authenticator).To(Equal("AllowAllAuthenticator"))
			Expect(config.CassandraConfig.Authorizer).To(Equal("AllowAllAuthorizer"))

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer, JmxCredentialsInitContainer)
		})

		It("with reaper disabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.auth.enabled": "false",
					"reaper.enabled":         "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.Authenticator).To(Equal("AllowAllAuthenticator"))
			Expect(config.CassandraConfig.Authorizer).To(Equal("AllowAllAuthorizer"))

			AssertInitContainerNamesMatch(cassdc, BaseConfigInitContainer, ConfigInitContainer)
		})
	})

	Context("when configuring the JVM heap for Cassandra 4.0", func() {
		It("at cluster-level only", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.version":             "4.0.0",
					"cassandra.heap.size":           "700M",
					"cassandra.heap.newGenSize":     "350M",
					"cassandra.datacenters[0].name": dcName,
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())

			Expect(config.JvmServerOptions).ToNot(BeNil())
			Expect(config.JvmServerOptions.InitialHeapSize).To(Equal("700M"))
			Expect(config.JvmServerOptions.MaxHeapSize).To(Equal("700M"))
			Expect(config.JvmServerOptions.YoungGenSize).To(Equal("350M"))
			Expect(config.JvmOptions).To(BeNil())
		})
	})

	Context("when configuring the Cassandra version and/or image", func() {
		cassandraVersionImageMap := map[string]string{
			"3.11.7":  "k8ssandra/cass-management-api:3.11.7-v0.1.28",
			"3.11.8":  "k8ssandra/cass-management-api:3.11.8-v0.1.28",
			"3.11.9":  "k8ssandra/cass-management-api:3.11.9-v0.1.27",
			"3.11.10": "k8ssandra/cass-management-api:3.11.10-v0.1.27",
			"3.11.11": "k8ssandra/cass-management-api:3.11.11-v0.1.28",
			"4.0.0":   "k8ssandra/cass-management-api:4.0.0-v0.1.28",
		}

		It("using the default version", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal("4.0.0"))
			Expect(cassdc.Spec.ServerImage).To(Equal("k8ssandra/cass-management-api:4.0.0-v0.1.28"))
		})

		It("using 3.11.7", func() {
			version := "3.11.7"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(cassandraVersionImageMap[version]))
		})

		It("using 3.11.8", func() {
			version := "3.11.8"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(cassandraVersionImageMap[version]))
		})

		It("using 3.11.9", func() {
			version := "3.11.9"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(cassandraVersionImageMap[version]))
		})

		It("using 3.11.10", func() {
			version := "3.11.10"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(cassandraVersionImageMap[version]))
		})

		It("using 3.11.11", func() {
			version := "3.11.11"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(cassandraVersionImageMap[version]))
		})

		It("using 4.0.0", func() {
			version := "4.0.0"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(cassandraVersionImageMap[version]))
		})

		It("using an unsupported version", func() {
			ver := "3.12.225"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": ver,
				},
			}

			renderedErr := renderTemplate(options)
			Expect(renderedErr).To(HaveOccurred())
		})

		It("using 3.11.11 and a custom image", func() {
			version := "3.11.11"
			repository := "my_cassandra"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version":          version,
					"cassandra.image.repository": repository,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(version))
			Expect(cassdc.Spec.ServerImage).To(Equal(DefaultRegistry + "/" + repository + ":latest"))
		})
	})

	It("enabling Cassandra auth with Stargate", func() {
		dcName := "test"
		clusterSize := 3
		clusterName := "auth-test"

		options := &helm.Options{
			KubectlOptions: defaultKubeCtlOptions,
			SetValues: map[string]string{
				"stargate.enabled":              "true",
				"cassandra.clusterName":         clusterName,
				"medusa.enabled":                "false",
				"reaper.enabled":                "false",
				"cassandra.auth.enabled":        "true",
				"cassandra.datacenters[0].name": dcName,
				"cassandra.datacenters[0].size": strconv.Itoa(clusterSize),
			},
		}

		Expect(renderTemplate(options)).To(Succeed())

		Expect(cassdc.Name).To(Equal(dcName))

		var config Config
		Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
		Expect(config.CassandraConfig.Authenticator).To(Equal("PasswordAuthenticator"))
		Expect(config.CassandraConfig.Authorizer).To(Equal("CassandraAuthorizer"))

		Expect(cassdc.Spec.Users).To(ConsistOf(cassdcv1beta1.CassandraUser{Superuser: true, SecretName: clusterName + "-stargate"}))
	})

	DescribeTable("check num_tokens",
		func(version, numTokens string, expectedTokens int) {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": version,
					// Need the following parameters to succeed in rendering
					"cassandra.datacenters[0].name": "dc-test",
					"cassandra.datacenters[0].size": "1",
				},
			}
			if numTokens != "" {
				options.SetValues["cassandra.datacenters[0].num_tokens"] = numTokens
			}

			Expect(renderTemplate(options)).To(Succeed())
			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.NumTokens).To(Equal(int64(expectedTokens)))
		},
		Entry("3.11.11 default", "3.11.11", "", 256),
		Entry("3.11.11 custom", "3.11.11", "16", 16),
	)
})

func verifyMedusaVolumeMounts(container *corev1.Container) {
	ExpectWithOffset(1, len(container.VolumeMounts)).To(Equal(4))
	ExpectWithOffset(1, container.VolumeMounts[0]).To(Equal(corev1.VolumeMount{Name: medusaConfigVolumeName, MountPath: "/etc/medusa"}))
}
