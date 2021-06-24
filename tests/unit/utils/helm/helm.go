package helm

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"path/filepath"
)

const (
	OPT_RENDER_TEMPLATE = "--rendertmpl"
)

type OutputDirective struct {
	Name string `json:"name"`
	Dir  string `json:"dir"`
}

// RenderAndUnmarshall renders a Helm template and invokes the provided function to unmarshal the result.
func RenderAndUnmarshall(templatePath string, options *helm.Options, helmChartPath string, HelmReleaseName string,
	unmarshallFunction func(string) error) error {

	renderedOutput, renderErr := helm.RenderTemplateE(
		GinkgoT(), options, helmChartPath, HelmReleaseName,
		[]string{templatePath},
	)

	// Use CLI option: --rendertmpl={"dir":"/tmp/foo", "name":"my-test.yaml"}
	opt := options.SetValues
	if opt != nil && opt[OPT_RENDER_TEMPLATE] != "" {

		fmt.Println("opt[OPT_RENDER_TEMPLATE]: ", opt[OPT_RENDER_TEMPLATE])
		var directive = OutputDirective{}
		json.Unmarshal([]byte(opt[OPT_RENDER_TEMPLATE]), &directive)

		fmt.Println("Render output directive found: ", directive.Name)
		if renderedOutput != "" {
			output(directive, []byte(renderedOutput))
		}

		if renderErr != nil {
			output(directive, []byte(renderErr.Error()))
		}
	}

	if renderErr == nil {
		unmarshalErr := unmarshallFunction(renderedOutput)
		ExpectWithOffset(1, unmarshalErr).ToNot(HaveOccurred(), "Unmarshall Error. "+
			"There is probably a type incompatibility issue in the test code. Make sure you are passing a pointer to "+
			"UnmarshalK8SYamlE in your unmarshall function.")
		return unmarshalErr
	}
	return renderErr
}

// Helper for test template diagnostics
func output(directive OutputDirective, data []byte) {

	outputDir, err := filepath.Abs(directive.Dir)
	if err != nil {
		Fail(fmt.Sprintf("Rendered output requested - invalid path: %s due to err: %s", outputDir, err.Error()))
	}

	err = os.MkdirAll(outputDir, 0700)
	if err != nil {
		Fail(fmt.Sprintf("Rendered output requested - unable to create: %s due to err: %s", outputDir, err.Error()))
	}

	outputFile := filepath.Join(directive.Dir, directive.Name)
	err = ioutil.WriteFile(outputFile, data, 0700)
	if err != nil {
		Fail(fmt.Sprintf("Rendered output requested - unable to write output to "+
			"file: %s due to err: %s", outputFile, err.Error()))
	}
	println(fmt.Sprintf("Rendered output requested - written to: %s ", directive))
}
