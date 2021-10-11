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

var (
	reaperOperatorNonRootFlag            = true
	reaperOperatorRunAsUser              = int64(65534)
	reaperOperatorRunAsGroup             = int64(65534)
	reaperOperatorReadOnlyRootFilesystem = true
	reaperOperatorDefaultSecurityContext = corev1.SecurityContext{RunAsNonRoot: &reaperOperatorNonRootFlag,
		RunAsUser: &reaperOperatorRunAsUser, RunAsGroup: &reaperOperatorRunAsGroup,
		ReadOnlyRootFilesystem: &reaperOperatorReadOnlyRootFilesystem}
)

var _ = Describe("Verify Reaper operator deployment template", func() {
	var (
		helmChartPath string
		err           error
		deployment    *appsv1.Deployment
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ReaperOperatorChartsPath)
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

			// Specific container security context checks
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).To(BeEquivalentTo(&reaperOperatorDefaultSecurityContext))

		})

		It("using customized securityContext values", func() {
			options := &helm.Options{
				ValuesFiles:    []string{"./testdata/reaper-operator-security-context-custom-values.yaml"},
				KubectlOptions: defaultKubeCtlOptions,
			}

			_ = renderTemplate(options)

			Expect(deployment.Kind).To(Equal("Deployment"))

			// Specific pod security context values expected
			Expect(*deployment.Spec.Template.Spec.SecurityContext.FSGroup).To(BeIdenticalTo(int64(1)))

			// Specific container security context checks
			Expect(deployment.Spec.Template.Spec.Containers[0].SecurityContext).ToNot(BeNil())
			Expect(*deployment.Spec.Template.Spec.Containers[0].SecurityContext.RunAsNonRoot).To(BeTrue())
		})
	})
})
