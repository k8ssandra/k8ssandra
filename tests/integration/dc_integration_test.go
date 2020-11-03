package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

// Deploy the rendered cassdc template on to a real Kubernetes cluster.
func TestCassdc(t *testing.T) {

	t.Parallel()

	helmChartPath, err := filepath.Abs("../../charts/k8ssandra-cluster")
	require.NoError(t, err)

	uniqueID := strings.ToLower(random.UniqueId())
	outputFileName := fmt.Sprintf("test-rendered-output-%s", uniqueID)
	releaseName := fmt.Sprintf("cassdc-release-%s", uniqueID)
	options := &helm.Options{
		SetStrValues: map[string]string{
			"name":        fmt.Sprintf("test-meta-name-%s", uniqueID),
			"clusterName": fmt.Sprintf("test-cluster-name-%s", uniqueID),
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", "k8ssandra-test-ns"),
	}

	defer helm.Delete(t, options, releaseName, true)
	defer cleanupRenderedOutput(t, outputFileName)

	kubeResourcePath, err := filepath.Abs(renderTemplate(t, options, helmChartPath, outputFileName))
	fmt.Println("Resource path: ", kubeResourcePath)

	require.NoError(t, err)

	helm.Install(t, options, helmChartPath, releaseName)
	k8s.WaitUntilAllNodesReady(t, options.KubectlOptions, 30, 2*time.Second)

}

// Render using template.
// Returns the output file created containing the rendered content.
func renderTemplate(t *testing.T, options *helm.Options, helmChartPath string, renderedOutputName string) string {

	renderedOutput := helm.RenderTemplate(
		t, options, helmChartPath, "k8ssandra-test-release",
		[]string{"templates/cassdc.yaml"},
	)

	require.NotEmpty(t, renderedOutput)

	outputFile := writeRenderedOutput(t, renderedOutputName, []byte(renderedOutput))
	require.NotEmpty(t, outputFile)
	return outputFile
}

/*
	cassandraOperatorProgress field in the CassandraDatacenter status has a value of Ready. The operator will set cassandraOperatorProgress to Ready when the desired state and the actual state match.
*/

// defer k8s.KubectlDelete(t, options.KubectlOptions, kubeResourcePath)
// k8s.KubectlApply(t, options.KubectlOptions, kubeResourcePath)

//k8s.KubectlApply(t, options.KubectlOptions, kubeResourcePath)
//k8s.WaitUntilAllNodesReady(t, options.KubectlOptions, 30, 2*time.Second)

// TODO: verify something here with assertion
// 	k8s.WaitUntilServiceAvailable(t, options, "hello-world-service", 10, 1*time.Second)
// service := k8s.GetService(t, options, "hello-world-service")
// url := fmt.Sprintf("http://%s", k8s.GetServiceEndpoint(t, options, service, 5000))
// Make an HTTP request to the URL and make sure it returns a 200 OK with the body "Hello, World".
// http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello world!", 30, 3*time.Second)

// Local test helper function for creating resource output used by tests.
// Returns location of the test resource.
// Validates for error condition and will fail the test if error occurs.
func writeRenderedOutput(t *testing.T, fileName string, data []byte) string {

	outputDir, err := filepath.Abs("./.output")
	err = os.MkdirAll(outputDir, 0755)
	destFile := filepath.Join(outputDir, fileName)
	if destFile != "" {
		err = ioutil.WriteFile(destFile, data, 0755)
	}

	require.NoError(t, err, "Failure during write of rendered output: %s", destFile)
	println(fmt.Sprintf("File: %s Created", destFile))
	return destFile
}

// Local test helper function for removal/cleanup of resource output used by tests.
// Validates for error condition and will fail the test if error occurs.
func cleanupRenderedOutput(t *testing.T, fileName string) {

	outputDir, err := filepath.Abs("./.output")
	destFile := filepath.Join(outputDir, fileName)
	err = os.Remove(destFile)

	require.NoError(t, err, "Failure during remove of rendered output: %s.", destFile)
	println(fmt.Sprintf("File: %s Deleted", destFile))
}
