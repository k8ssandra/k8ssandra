package unit_test

import (
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Verify reaper user secret template", func() {
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
			[]string{"templates/reaper/reaper-user-secret.yaml"},
		)

		if err == nil {
			err = helm.UnmarshalK8SYamlE(GinkgoT(), renderedOutput, secret)
		}

		return err
	}

	Context("generating reaper user secret", func() {
		It("specifying reaper user username", func() {
			username := "reaper_admin"
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":                clusterName,
					"cassandra.auth.enabled":               "true",
					"repair.reaper.cassandraUser.username": username,
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(secret.Name).To(Equal(clusterName + "-reaper"))
			Expect(string(secret.Data["username"])).To(Equal(username))
			Expect(len(secret.Data["password"])).To(Equal(20))
		})

		It("using default username for reaper user secret", func() {
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":  clusterName,
					"cassandra.auth.enabled": "true",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())
			Expect(secret.Name).To(Equal(clusterName + "-reaper"))
			Expect(string(secret.Data["username"])).To(Equal("reaper"))
			Expect(len(secret.Data["password"])).To(Equal(20))
		})

		It("not generating reaper user secret", func() {
			clusterName := "secret-test"
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.clusterName":              clusterName,
					"cassandra.auth.enabled":             "true",
					"repair.reaper.cassandraUser.secret": "reaper-secret",
				},
			}

			err := renderTemplate(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("could not find template"))
		})
	})
})
