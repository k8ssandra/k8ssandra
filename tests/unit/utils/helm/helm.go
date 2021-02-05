package helm

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// RenderAndUnmarshall renders a Helm template and invokes the provided function to unmarshal the result.
func RenderAndUnmarshall(templatePath string, options *helm.Options, helmChartPath string, HelmReleaseName string, unmarshallFunction func(string) error) error {
	renderedOutput, renderErr := helm.RenderTemplateE(
		GinkgoT(), options, helmChartPath, HelmReleaseName,
		[]string{templatePath},
	)
	if renderErr == nil {
		unmarshalErr := unmarshallFunction(renderedOutput)
		ExpectWithOffset(1, unmarshalErr).ToNot(HaveOccurred(), "Unmarshall Error. There is probably a type incompatibility issue in the test code. Make sure you are passing a pointer to UnmarshalK8SYamlE in your unmarshall function.")
		return unmarshalErr
	}
	return renderErr
}
