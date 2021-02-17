package unit_test

import (
	"path/filepath"

	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Verify medusa config template", func() {
	var (
		helmChartPath string
		secret        *corev1.Secret
	)

	BeforeEach(func() {
		path, err := filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
		helmChartPath = path
		secret = &corev1.Secret{}
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/medusa/medusa-config.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, secret)
			})
	}

	Context("generating medusa storage properties", func() {
		DescribeTable("render template",
			func(storageType string, expected bool) {
				options := &helm.Options{
					KubectlOptions: defaultKubeCtlOptions,
					SetValues: map[string]string{
						"backupRestore.medusa.enabled":      "true",
						"backupRestore.medusa.storage":      storageType,
						"backupRestore.medusa.bucketName":   "testbucket",
						"backupRestore.medusa.bucketSecret": "secretkey",
					},
				}
				Expect(renderTemplate(options)).To(Succeed())
			},
			Entry("supported s3", "s3", true),
			Entry("supported s3 compatible", "s3_compatible", true),
			Entry("supported gcs", "gcs", true),
			Entry("supported local", "local", true),
			Entry("unsupported azure", "azure", false),
			Entry("unsupported ibm_storage (use s3_compatible instead)", "ibm_storage", false),
			Entry("supported value", "random", false),
		)
	})
})
