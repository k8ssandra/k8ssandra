package unit_test

import (
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1batch "k8s.io/api/batch/v1"
)

var _ = Describe("Verify Cleaner job template", func() {
	var (
		helmChartPath string
		err           error
		cleanerJob    *v1batch.Job
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		cleanerJob = &v1batch.Job{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/cleaner/batch_job.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, cleanerJob)
			})
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			renderTemplate(options)

			By("checking that correct hook annotations are present")
			Expect(cleanerJob.Annotations).Should(HaveKeyWithValue(HelmHookAnnotation, "pre-delete"))
			Expect(cleanerJob.Annotations).Should(HaveKeyWithValue(HelmHookPreDeleteAnnotation, "hook-succeeded,before-hook-creation"))

			Expect(len(cleanerJob.Spec.Template.Spec.Containers)).To(Equal(1))
			Expect(len(cleanerJob.Spec.Template.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cleanerJob.Spec.Template.Spec.Containers[0].Env[0].Name).To(Equal("POD_NAMESPACE"))
		})
	})
})
