package unit_test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"path/filepath"
)

var _ = Describe("Verify Medusa operator deployment template", func() {
	var (
		helmChartPath string
		err           error
		deployment    *appsv1.Deployment

		medusaOperatorNonRootFlag            = true
		medusaOperatorRunAsUser              = int64(65534)
		medusaOperatorRunAsGroup             = int64(65534)
		medusaOperatorReadOnlyRootFilesystem = true
		medusaOperatorDefaultSecurityContext = corev1.SecurityContext{RunAsNonRoot: &medusaOperatorNonRootFlag, RunAsUser: &medusaOperatorRunAsUser,
			RunAsGroup: &medusaOperatorRunAsGroup, ReadOnlyRootFilesystem: &medusaOperatorReadOnlyRootFilesystem}
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(MedusaOperatorChartsPath)
		deployment = &appsv1.Deployment{}
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/deployment.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, deployment)
			})
	}

	Context("by rendering deployment with options", func() {
		It("using only default securityContext values", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)

			Expect(deployment.Kind).To(Equal("Deployment"))
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeEquivalentTo(&medusaOperatorDefaultSecurityContext))
		})

		It("using customized securityContext values", func() {
			options := &helm.Options{
				ValuesFiles:    []string{"./testdata/medusa-operator-security-context-custom-values.yaml"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)

			Expect(deployment.Kind).To(Equal("Deployment"))
			Expect(*deployment.Spec.Template.Spec.SecurityContext.FSGroup).To(BeIdenticalTo(int64(1)))
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
			Expect(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.ReadOnlyRootFilesystem).To(BeTrue())
			Expect(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).To(BeTrue())
		})
	})
})
