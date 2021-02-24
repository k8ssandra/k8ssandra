package unit_test

import (
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Verify medusa user secret template", func() {
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
		return helmUtils.RenderAndUnmarshall("templates/medusa/medusa-user-secret.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, secret)
			})
	}

	Context("generating medusa user secret", func() {
		It("specifying medusa user username", func() {
			username := "medusa_admin"
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":         clusterName,
					"cassandra.auth.enabled":        "true",
					"medusa.enabled":                "true",
					"medusa.cassandraUser.username": username,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(secret.Name).To(Equal(clusterName + "-medusa"))
			Expect(string(secret.Data["username"])).To(Equal(username))
			Expect(len(secret.Data["password"])).To(Equal(20))
		})

		It("using default username for medusa user secret", func() {
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":  clusterName,
					"cassandra.auth.enabled": "true",
					"medusa.enabled":         "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(secret.Name).To(Equal(clusterName + "-medusa"))
			Expect(string(secret.Data["username"])).To(Equal("medusa"))
			Expect(len(secret.Data["password"])).To(Equal(20))
		})

		It("not generating medusa user secret", func() {
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":       clusterName,
					"cassandra.auth.enabled":      "true",
					"medusa.cassandraUser.secret": "medusa-secret",
				},
			}

			err := renderTemplate(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})
	})
})
