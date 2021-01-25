package unit_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	cassdcv1beta1 "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
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
}

type JvmOptions struct {
	AdditionalJvmOptions []string `json:"additional-jvm-opts"`
	InitialHeapSize      string   `json:"initial_heap_size"`
	MaxHeapSize          string   `json:"max_heap_size"`
	YoungGenSize         string   `json:"heap_size_young_generation"`
}

type Config struct {
	CassandraConfig CassandraConfig `json:"cassandra-yaml"`
	JvmOptions      JvmOptions      `json:"jvm-options"`
}

var (
	reaperInstanceValue    = fmt.Sprintf("%s-reaper-k8ssandra", helmReleaseName)
	medusaConfigVolumeName = fmt.Sprintf("%s-medusa-config-k8ssandra", helmReleaseName)
	defaultKubeCtlOptions  = k8s.NewKubectlOptions("", "", defaultTestNamespace)
)

var _ = Describe("Verify CassandraDatacenter template", func() {
	var (
		helmChartPath string
		cassdc        *cassdcv1beta1.CassandraDatacenter
	)

	BeforeEach(func() {
		path, err := filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		helmChartPath = path
		cassdc = &cassdcv1beta1.CassandraDatacenter{}
	})

	renderTemplate := func(options *helm.Options) (error, error) {

		templatePath := filepath.Join(".", "/templates/cassandra/cassdc.yaml")
		renderedOutput, renderErr := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{templatePath},
		)
		unmarshalErr := helm.UnmarshalK8SYamlE(GinkgoT(), renderedOutput, cassdc)
		return renderErr, unmarshalErr
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
			Expect(cassdc.Annotations).Should(HaveKeyWithValue(reaperInstanceAnnotation, reaperInstanceValue))

			initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
			Expect(len(initContainers)).To(Equal(2))
			Expect(initContainers[0].Name).To(Equal("server-config-init"))

			// Verify initContainers includes JMX credentials
			Expect(initContainers[1].Name).To(Equal("jmx-credentials"))
			// Verify LOCAL_JMX value
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(1))
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Name).To(Equal("LOCAL_JMX"))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Value).To(Equal("no"))
			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(false))
			Expect(*cassdc.Spec.DockerImageRunsAsCassandra).To(BeFalse())

			// Server version and mgmt-api image specified
			Expect(cassdc.Spec.ServerVersion).ToNot(BeEmpty())
			Expect(cassdc.Spec.ServerImage).ToNot(BeEmpty())

			// JVM heap options -- default to settings as defined in cassdc.yaml
			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("800M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("800M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("1.6Gi"))
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

		It("using cassandra 3.11.7", func() {
			cassandraVersion := "3.11.7"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": cassandraVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(cassandraVersion))
			Expect(cassdc.Spec.ServerImage).To(Equal("datastax/cassandra-mgmtapi-3_11_7:v0.1.17"))
		})

		It("using cassandra 3.11.8", func() {
			cassandraVersion := "3.11.8"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": cassandraVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(cassandraVersion))
			Expect(cassdc.Spec.ServerImage).To(Equal("datastax/cassandra-mgmtapi-3_11_8:v0.1.17"))
		})

		It("using cassandra 3.11.9", func() {
			cassandraVersion := "3.11.9"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": cassandraVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(cassandraVersion))
			Expect(cassdc.Spec.ServerImage).To(Equal("datastax/cassandra-mgmtapi-3_11_9:v0.1.17"))
		})

		It("using cassandra with unsupported version", func() {
			cassandraVersion := "3.12.225"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": cassandraVersion,
				},
			}

			renderedErr, _ := renderTemplate(options)
			Expect(renderedErr).To(HaveOccurred())
		})

		It("disabling Cassandra auth", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.auth.enabled": "false",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.CassandraConfig.Authenticator).To(Equal("AllowAllAuthenticator"))
			Expect(config.CassandraConfig.Authorizer).To(Equal("AllowAllAuthorizer"))
		})

		It("enabling Cassandra auth", func() {
			dcName := "test"
			clusterSize := 3

			authCachePeriod := int64(7200000)
			cacheValidityPeriod := authCachePeriod + 1
			cacheUpdateInterval := authCachePeriod + 2

			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
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
			Expect(config.JvmOptions.AdditionalJvmOptions).To(ConsistOf(
				"-Dcassandra.system_distributed_replication_dc_names="+dcName,
				"-Dcassandra.system_distributed_replication_per_dc="+strconv.Itoa(clusterSize),
			))

		})

		It("providing superuser secret", func() {
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

		It("providing superuser username", func() {
			clusterName := "superuser-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":             clusterName,
					"cassandra.auth.enabled":            "true",
					"cassandra.auth.superuser.username": "admin",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.SuperuserSecretName).To(Equal(clusterName + "-superuser"))
		})

		It("disabling reaper", func() {
			options := &helm.Options{
				SetValues:      map[string]string{"repair.reaper.enabled": "false"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(cassdc.Annotations).ShouldNot(HaveKeyWithValue(reaperInstanceAnnotation, reaperInstanceValue))

			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(1))
			// No env slice should be present
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env).To(BeNil())
			// No initcontainers slice should be present
			Expect(cassdc.Spec.PodTemplateSpec.Spec.InitContainers).To(BeNil())
		})

		It("enabling only medusa", func() {
			options := &helm.Options{
				SetValues:      map[string]string{"backupRestore.medusa.enabled": "true", "repair.reaper.enabled": "false"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			// InitContainers should only server-config-init and medusa-restore
			initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
			Expect(len(initContainers)).To(Equal(3))
			Expect(initContainers[0].Name).To(Equal("server-config-init"))
			// Verify initContainers includes jolokia which medusa needs
			Expect(cassdc.Spec.PodTemplateSpec.Spec.InitContainers[1].Name).To(Equal("get-jolokia"))
			// Verify initContainers includes medusa-restore
			Expect(cassdc.Spec.PodTemplateSpec.Spec.InitContainers[2].Name).To(Equal("medusa-restore"))
			// Two containers, medusa and cassandra
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(2))
			// Cassandra container should have JVM_EXTRA_OPTS for jolokia
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Name).To(Equal("JVM_EXTRA_OPTS"))
			// Second container should be medusa
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].Name).To(Equal("medusa"))

			// Verify volumeMounts and volumes
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].VolumeMounts)).To(Equal(4))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].VolumeMounts[0].Name).To(Equal(medusaConfigVolumeName))

			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Volumes)).To(Equal(3))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Volumes[0].Name).To(Equal(medusaConfigVolumeName))
		})

		It("enabling reaper and medusa", func() {
			// Simple verification that both have properties correctly applied
			options := &helm.Options{
				SetValues:      map[string]string{"backupRestore.medusa.enabled": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			// Verify initContainers
			initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
			Expect(len(initContainers)).To(Equal(4))
			Expect(initContainers[0].Name).To(Equal("server-config-init"))
			Expect(initContainers[1].Name).To(Equal("jmx-credentials"))
			Expect(initContainers[2].Name).To(Equal("get-jolokia"))
			Expect(initContainers[3].Name).To(Equal("medusa-restore"))

			// Verify containers
			containers := cassdc.Spec.PodTemplateSpec.Spec.Containers
			Expect(len(containers)).To(Equal(2))
			Expect(containers[0].Name).To(Equal("cassandra"))
			Expect(containers[1].Name).To(Equal("medusa"))
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

			renderedErr, _ := renderTemplate(options)

			Expect(renderedErr).ToNot(BeNil())
			Expect(renderedErr.Error()).To(ContainSubstring("set resource limits/requests when enabling allowMultipleNodesPerWorker"))

		})

		It("setting JVM heap settings at cluster-level only", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.datacenters[0].name": dcName,
					"cassandra.datacenters[0].heap": "",
					"cassandra.heap.size":           "100M",
					"cassandra.heap.newGenSize":     "200M",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("100M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("100M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("200M"))
		})

		// Note: currently only one DC supported, to be expanded in future release.
		It("setting JVM heap settings at dc level overriding cluster level", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap.size":                      "100M",
					"cassandra.heap.newGenSize":                "200M",
					"cassandra.datacenters[0].heap.size":       "300M",
					"cassandra.datacenters[0].heap.newGenSize": "600M",
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
			Expect(config.JvmOptions.YoungGenSize).To(Equal("600M"))
		})

		// Note: currently only one DC supported, to be expanded in future release.
		It("setting JVM heap settings at datacenter WITHOUT cluster-level settings", func() {

			dcName := "dc1"
			options := &helm.Options{
				SetValues: map[string]string{
					"cassandra.heap":                           "",
					"cassandra.datacenters[0].name":            dcName,
					"cassandra.datacenters[0].heap.size":       "600M",
					"cassandra.datacenters[0].heap.newGenSize": "1.2Gi",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			var config Config
			Expect(json.Unmarshal(cassdc.Spec.Config, &config)).To(Succeed())
			Expect(config.JvmOptions).ToNot(BeNil())
			Expect(config.JvmOptions.InitialHeapSize).To(Equal("600M"))
			Expect(config.JvmOptions.MaxHeapSize).To(Equal("600M"))
			Expect(config.JvmOptions.YoungGenSize).To(Equal("1.2Gi"))
		})
	})
})
