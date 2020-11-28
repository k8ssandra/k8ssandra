package tests

import (
	"fmt"
	cassdcv1beta1 "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var releaseName = "k8ssandra-test"
var defaultTestNamespace = "k8ssandra"
var defaultKubeCtlOptions = k8s.NewKubectlOptions("", "", defaultTestNamespace)

var _ = Describe("Verify CassandraDatacenter template", func() {
	var (
		helmChartPath string
		err error
		cassdc *cassdcv1beta1.CassandraDatacenter
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs("../charts/k8ssandra-cluster")
		Expect(err).To(BeNil())
		cassdc = &cassdcv1beta1.CassandraDatacenter{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) {
		renderedOutput := helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, releaseName,
			[]string{"templates/cassdc.yaml"},
		)

		helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, cassdc)
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)

			Expect(cassdc.Kind).To(Equal("CassandraDatacenter"))

			// Reaper should be enabled in default - verify
			// Verify reaper annotation is set
			Expect(cassdc.Annotations).Should(HaveKeyWithValue("reaper.cassandra-reaper.io/instance", fmt.Sprintf("%s-reaper-k8ssandra", releaseName)))
			// Initcontainer should only have one (reaper, not medusa)
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.InitContainers)).To(Equal(1))
			// Verify initContainers includes JMX credentials
			Expect(cassdc.Spec.PodTemplateSpec.Spec.InitContainers[0].Name).To(Equal("jmx-credentials"))
			// Verify LOCAL_JMX value
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(1))
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Name).To(Equal("LOCAL_JMX"))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Value).To(Equal("no"))
		})

		It("disabling reaper", func() {
			options := &helm.Options{
				SetValues:   map[string]string{"repair.reaper.enabled": "false"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
			Expect(cassdc.Annotations).ShouldNot(HaveKeyWithValue("reaper.cassandra-reaper.io/instance", fmt.Sprintf("%s-reaper-k8ssandra", releaseName)))
		})

		It("enabling only medusa", func() {
			options := &helm.Options{
				SetValues:   map[string]string{"backupRestore.medusa.enabled": "true", "repair.reaper.enabled": "false"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)

			// Verify medusa is present
			// Initcontainer should only have one (medusa)
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.InitContainers)).To(Equal(2))
			// Verify initContainers includes jolokia which medusa needs
			Expect(cassdc.Spec.PodTemplateSpec.Spec.InitContainers[0].Name).To(Equal("get-jolokia"))
			// Verify initContainers includes medusa-restore
			Expect(cassdc.Spec.PodTemplateSpec.Spec.InitContainers[1].Name).To(Equal("medusa-restore"))
			// Two containers, medusa and cassandra
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(2))
			// Cassandra container should have JVM_EXTRA_OPTS for jolokia
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[0].Env[0].Name).To(Equal("JVM_EXTRA_OPTS"))
			// Second container should be medusa
			Expect(cassdc.Spec.PodTemplateSpec.Spec.Containers[1].Name).To(Equal("medusa"))
		})

		It("enabling reaper and medusa", func() {
			// Simple verification that both have properties correctly applied
			options := &helm.Options{
				SetValues:   map[string]string{"backupRestore.medusa.enabled": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)

			// Verify both are present
			// Initcontainer should only have jmx and jolokia
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.InitContainers)).To(Equal(3))
			// Two containers, medusa and cassandra
			Expect(len(cassdc.Spec.PodTemplateSpec.Spec.Containers)).To(Equal(2))
		})
	})
})
