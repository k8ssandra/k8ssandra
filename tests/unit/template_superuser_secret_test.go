package unit_test

import (
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Verify superuser secret template", func() {

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
		return helmUtils.RenderAndUnmarshall("templates/cassandra/superuser-secret.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, secret)
			})
	}

	Context("generating superuser secret", func() {
		It("specifying superuser username", func() {
			username := "admin"
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":             clusterName,
					"cassandra.auth.superuser.username": username,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(secret.Name).To(Equal(clusterName + "-superuser"))
			Expect(string(secret.Data["username"])).To(Equal(username))
			Expect(len(secret.Data["password"])).To(Equal(20))
		})
	})

	Context("not generating superuser secret", func() {
		It("specifying superuser secret", func() {
			secretName := "test-secret"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.auth.superuser.secret": secretName,
				},
			}

			err := renderTemplate(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})

		It("disabling auth", func() {
			username := "admin"
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":             clusterName,
					"cassandra.auth.superuser.username": username,
					"cassandra.auth.enabled":            "false",
				},
			}

			err := renderTemplate(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})

	})
})
