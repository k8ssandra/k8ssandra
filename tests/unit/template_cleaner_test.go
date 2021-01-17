package unit_test

import (
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
		helmChartPath, err = filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		cleanerJob = &v1batch.Job{}
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) {
		renderedOutput := helm.RenderTemplate(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/cleaner/batch_job.yaml"},
		)

		helm.UnmarshalK8SYaml(GinkgoT(), renderedOutput, cleanerJob)
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					// Note that the test fails without this. I am not sure why because
					// this is the default in k8ssandra/values.yaml.
					"cass-operator.enabled": "true",
				},
			}

			renderTemplate(options)

			By("checking that correct hook annotations are present")
			Expect(cleanerJob.Annotations).Should(HaveKeyWithValue(helmHookAnnotation, "pre-delete"))
			Expect(cleanerJob.Annotations).Should(HaveKeyWithValue(helmHookPreDeleteAnnotation, "hook-succeeded,before-hook-creation"))

			Expect(len(cleanerJob.Spec.Template.Spec.Containers)).To(Equal(1))
			Expect(len(cleanerJob.Spec.Template.Spec.Containers[0].Env)).To(Equal(1))
			Expect(cleanerJob.Spec.Template.Spec.Containers[0].Env[0].Name).To(Equal("POD_NAMESPACE"))
		})
	})
})
