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

	renderTemplate := func(options *helm.Options) error {
		renderedOutput, err := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/cassdc.yaml"},
		)

		if err == nil {
			err = helm.UnmarshalK8SYamlE(GinkgoT(), renderedOutput, cassdc)
		}

		return err
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
		})

		It("override clusterName", func() {
			clusterName := "test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.clusterName": clusterName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ClusterName).To(Equal(clusterName))
		})

		It("override datacenterName", func() {
			dcName := "test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.datacenterName": dcName,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Name).To(Equal(dcName))
		})

		It("override size", func() {
			size := "3"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.size": size,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.Size, 3)
		})

		It("use cassandraVersion 3.11.7", func() {
			cassandraVersion := "3.11.7"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.cassandraVersion": cassandraVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(cassandraVersion))
			Expect(cassdc.Spec.ServerImage).To(Equal("datastax/cassandra-mgmtapi-3_11_7:v0.1.17"))
		})

		It("use cassandraVersion 3.11.8", func() {
			cassandraVersion := "3.11.8"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.cassandraVersion": cassandraVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(cassandraVersion))
			Expect(cassdc.Spec.ServerImage).To(Equal("datastax/cassandra-mgmtapi-3_11_8:v0.1.17"))
		})

		It("use cassandraVersion 3.11.9", func() {
			cassandraVersion := "3.11.9"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.cassandraVersion": cassandraVersion,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.ServerVersion).To(Equal(cassandraVersion))
			Expect(cassdc.Spec.ServerImage).To(Equal("datastax/cassandra-mgmtapi-3_11_9:v0.1.17"))
		})

		It("use cassandraVersion with unsupported value", func() {
			cassandraVersion := "3.12.225"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.cassandraVersion": cassandraVersion,
				},
			}

			err := renderTemplate(options)

			Expect(err).To(HaveOccurred())
		})

		It("disabling Cassandra auth", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.configuration.auth.enabled": "false",
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
					"k8ssandra.datacenterName":                               dcName,
					"k8ssandra.size":                                         strconv.Itoa(clusterSize),
					"k8ssandra.configuration.auth.enabled":                   "true",
					"k8ssandra.configuration.auth.cacheValidityPeriodMillis": strconv.FormatInt(cacheValidityPeriod, 10),
					"k8ssandra.configuration.auth.cacheUpdateIntervalMillis": strconv.FormatInt(cacheUpdateInterval, 10),
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
					"k8ssandra.allowMultipleNodesPerWorker": "true",
					"k8ssandra.resources.limits.memory":     "2Gi",
					"k8ssandra.resources.limits.cpu":        "1",
					"k8ssandra.resources.requests.memory":   "2Gi",
					"k8ssandra.resources.requests.cpu":      "1"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(true))
		})

		It("setting allowMultipleNodesPerWorker to false", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"k8ssandra.allowMultipleNodesPerWorker": "false",
					"k8ssandra.resources.limits.memory":     "2Gi",
					"k8ssandra.resources.limits.cpu":        "1",
					"k8ssandra.resources.requests.memory":   "2Gi",
					"k8ssandra.resources.requests.cpu":      "1",
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
					"k8ssandra.allowMultipleNodesPerWorker": "false",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())

			Expect(cassdc.Spec.AllowMultipleNodesPerWorker).To(Equal(false))
		})

		It("setting allowMultipleNodesPerWorker to true without resources", func() {
			options := &helm.Options{
				SetValues: map[string]string{
					"k8ssandra.allowMultipleNodesPerWorker": "true",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			err := renderTemplate(options)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("set resource limits/requests when enabling allowMultipleNodesPerWorker"))

		})
	})
})
