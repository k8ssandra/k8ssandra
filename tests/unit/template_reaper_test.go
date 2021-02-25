package unit_test

import (
	"path/filepath"

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

			renderTemplate(options)

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

			renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.CassandraBackend.CassandraDatacenter.Name).To(Equal(targetDcName))
		})

		It("modifying autoscheduling option", func() {
			options := &helm.Options{
				SetStrValues:   map[string]string{"reaper.autoschedule": "true"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.AutoScheduling).ToNot(BeNil())
		})

		It("modifying secret options", func() {
			options := &helm.Options{
				SetStrValues: map[string]string{
					"reaper.jmx.secret":           "somethingelse",
					"reaper.cassandraUser.secret": "cassandraSpecial",
				},
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)
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

			renderTemplate(options)
			Expect(reaper.Spec.ServerConfig.JmxUserSecretName).To(HavePrefix("nowyouseeme"))
		})
	})
})
