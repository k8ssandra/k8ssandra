package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify medusa config template", func() {
	var (
		helmChartPath string
	)

	BeforeEach(func() {
		path, err := filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		helmChartPath = path
	})

	renderTemplate := func(options *helm.Options) bool {
		_, renderErr := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, HelmReleaseName,
			[]string{"templates/medusa/medusa-config.yaml"})

		return renderErr == nil
	}

	Context("generating medusa storage properties", func() {
		DescribeTable("render template",
			func(storageType string, expected bool) {
				options := &helm.Options{
					KubectlOptions: defaultKubeCtlOptions,
					SetValues: map[string]string{
						"medusa.enabled":      "true",
						"medusa.storage":      storageType,
						"medusa.bucketName":   "testbucket",
						"medusa.bucketSecret": "secretkey",
					},
				}
				Expect(renderTemplate(options)).To(Equal(expected))
			},
			Entry("supported s3", "s3", true),
			Entry("supported s3 compatible", "s3_compatible", true),
			Entry("supported gcs", "google_storage", true),
			Entry("supported local", "local", true),
			Entry("unsupported azure", "azure", false),
			Entry("unsupported ibm_storage (use s3_compatible instead)", "ibm_storage", false),
			Entry("supported value", "random", false),
		)
	})
})
