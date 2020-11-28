package tests

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	api "github.com/k8ssandra/reaper-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("Verify Reaper template", func() {
	var (
		helmChartPath string
		err error
		reaper *api.Reaper
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs("../charts/k8ssandra-cluster")
		Expect(err).To(BeNil())
		reaper = &api.Reaper{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) {
		renderedOutput := helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, "k8ssandra-test",
			[]string{"templates/reaper.yaml"},
		)

		helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, reaper)
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: k8s.NewKubectlOptions("", "", "k8ssandra"),
			}

			renderTemplate(options)
			Expect(string(reaper.Spec.ServerConfig.StorageType)).To(Equal("cassandra"))
			Expect(reaper.Kind).To(Equal("Reaper"))
		})

		It("changing datacenter name", func() {
			targetDcName := "reaper-dc"
			options := &helm.Options{
				SetStrValues:   map[string]string{"datacenterName": targetDcName},
				KubectlOptions: k8s.NewKubectlOptions("", "", "k8ssandra"),
			}

			renderTemplate(options)
			// Requires new version of reaper-operator and k8ssandra PR merged
			//Expect(reaper.Spec.ServerConfig.CassandraBackend.CassandraDatacenter.Name).To(Equal(targetDcName))
		})

		It("modifying autoscheduling option", func() {
			options := &helm.Options{
				SetStrValues:   map[string]string{"autoscheduling": "true"},
				KubectlOptions: k8s.NewKubectlOptions("", "", "k8ssandra"),
			}

			renderTemplate(options)
			// Requires new version of reaper-operator and k8ssandra PR merged
			//Expect(reaper.Spec.ServerConfig.Autoscheduler).ToNot(BeNil())
		})
	})
})
