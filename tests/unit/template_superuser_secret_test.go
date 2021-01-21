package unit_test

import (
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
		path, err := filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		helmChartPath = path
		secret = &corev1.Secret{}
	})

	renderTemplate := func(options *helm.Options) error {
		renderedOutput, err := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/cassandra/superuser-secret.yaml"},
		)

		if err == nil {
			err = helm.UnmarshalK8SYamlE(GinkgoT(), renderedOutput, secret)
		}

		return err
	}

	Context("generating superuser secret", func() {
		It("specifying superuser username", func() {
			username := "admin"
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"k8ssandra.clusterName":                           clusterName,
					"k8ssandra.configuration.auth.superuser.username": username,
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
					"k8ssandra.configuration.auth.superuser.secret": secretName,
				},
			}

			err := renderTemplate(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})
	})
})
