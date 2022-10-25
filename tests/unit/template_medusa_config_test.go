package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
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
		return helmUtils.RenderAndUnmarshall("templates/medusa/medusa-config.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &corev1.ConfigMap{})
			}) == nil
	}

	Context("generating medusa storage properties", func() {
		DescribeTable("render template",
			func(storageType string, expected bool) {
				options := &helm.Options{
					KubectlOptions: defaultKubeCtlOptions,
					SetValues: map[string]string{
						"medusa.enabled":                 "true",
						"medusa.storage":                 storageType,
						"medusa.bucketName":              "testbucket",
						"medusa.storageSecret":           "secretkey",
						"medusa.podStorage.size":         "30Gi",
						"medusa.podStorage.storageClass": "nfs",
					},
				}
				Expect(renderTemplate(options)).To(Equal(expected))
			},
			Entry("supported s3", "s3", true),
			Entry("supported s3 compatible", "s3_compatible", true),
			Entry("supported s3 rgw", "s3_rgw", true),
			Entry("supported google_storage", "google_storage", true),
			Entry("supported azure_blobs", "azure_blobs", true),
			Entry("supported local", "local", true),
			Entry("unsupported ibm_storage (use s3_compatible instead)", "ibm_storage", false),
			Entry("unsupported value", "random", false),
		)
	})
})
