package util

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"k8s.io/apimachinery/pkg/util/wait"
)

// GenerateNamespaceName provides namespace name for test inclusion.
func GenerateNamespaceName(namespacePrefix string) string {
	testID := strings.ToLower(random.UniqueId())
	name := fmt.Sprintf("%s-%s", namespacePrefix, testID)
	return name
}

// GetNamespaces provides a string array list of namspaces found during test.
func GetNamespaces(t *testing.T, kubeOptions *k8s.KubectlOptions) ([]string, error) {
	lookupResult, err := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "namespace", "-o", "name")
	namespaces := strings.Split(lookupResult, "\n")
	return namespaces, err
}

// CleanupRelease performs cleanup of specific release used for test.
func CleanupRelease(t *testing.T, helmOptions *helm.Options, releaseName string) {
	fmt.Println("Performing a Delete of the release name: ", releaseName)
	helmDeleteReleaseError := helm.DeleteE(t, helmOptions, releaseName, true)
	if helmDeleteReleaseError != nil {
		fmt.Println("Helm delete release: ", helmDeleteReleaseError)
	}
}

// CreateTestNamespace creates provided namespace while providing visiblity into existing namspace state.
// Returns error if unable to get the list of namespaces.
func CreateTestNamespace(t *testing.T, kubeOptions *k8s.KubectlOptions, releaseName string, namespace string) error {
	k8s.CreateNamespace(t, kubeOptions, namespace)
	namespaces, err := GetNamespaces(t, kubeOptions)
	for _, ns := range namespaces {
		logger.Log(t, "Namespace: ", ns)
	}
	return err
}

// DeleteResources performs cleanup of target resources provided unique to namespace and/or resourceName.
// Returns a map of delete resource results and/or errors for the target list provided.
func DeleteResources(t *testing.T, kubeOptions *k8s.KubectlOptions, targets []string, namespace string, resourceName string) map[string]string {

	resultMap := map[string]string{}
	for _, target := range targets {

		var lookup string = ""
		var err error = nil
		if strings.ToLower(target) == "secret" {

			tokens, secretErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "Secret", "-A", "--field-selector",
				fmt.Sprintf("metadata.namespace=%s", namespace), "-o", "name")

			tokenList := strings.Split(tokens, "\n")
			if secretErr == nil && tokenList != nil {
				for _, token := range tokenList {
					Log(t, "Secret", "Found token: ", token, secretErr)
					deleteResult, deleteErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "--namespace", namespace, "delete", token)
					Log(t, "Secret", "Delete", deleteResult, deleteErr)
				}
			}
			if secretErr != nil {
				resultMap[target] = fmt.Sprintf(" 'get' returned error: %s", secretErr.Error())
				continue
			}

			// Identify any of the Secret tokens removed.
			resultMap[target] = tokens

		} else {
			lookup, err = k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", target, "-A", "--field-selector",
				fmt.Sprintf("metadata.name=%s", resourceName), "--namespace", namespace, "-o", "name")
		}

		if err != nil {
			resultMap[target] = fmt.Sprintf(" 'get' returned error: %s", err.Error())
			continue
		}

		if lookup != "" && lookup != "No resources found" {

			noNamespaceOptions := k8s.NewKubectlOptions("", "", "")
			deleteResult, err := k8s.RunKubectlAndGetOutputE(t, noNamespaceOptions, "delete", lookup)
			Pause(t, "Pausing for Resource to be deleted", 1*time.Second, 2*time.Second)

			if err != nil {
				resultMap[target] = fmt.Sprintf(" 'delete' returned error: %s", err.Error())
				continue
			}
			resultMap[target] = deleteResult
		}
	}
	return resultMap
}

// CleanupCRD performs lookup of resource targets specific to CRD metatdata name.
// Removes the dependency resources followed by the CRD itself.
func CleanupCRD(t *testing.T, kubeOptions *k8s.KubectlOptions, metadataName string, crdDefinitionYaml string) {

	lookupResult, lookupError := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "crd", "-A", "--field-selector",
		fmt.Sprintf("metadata.name=%s", metadataName), "--namespace", kubeOptions.Namespace, "-o", "name")

	Pause(t, "Get CRD", 1*time.Second, 2*time.Second)
	Log(t, "CRD", "Lookup", lookupResult, lookupError)

	if lookupError == nil && lookupResult != "" {
		resources := []string{
			"ClusterRoleBinding",
			"ClusterRole",
			"ValidatingWebhookConfiguration",
			"MutatingWebhookConfiguration",
			"Secret",
			"Service",
		}

		results := DeleteResources(t, kubeOptions, resources, kubeOptions.Namespace, "cass-operator-webhook")
		LogResults(t, results)

		patchResult, patchError := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "patch", "crd/cassandradatacenters.cassandra.datastax.com",
			"-p", "{\"metadata\":{\"finalizers\":[]}}", "--type=merge")

		Log(t, "CRD", "Patch", patchResult, patchError)

		deleteResult, deleteError := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "delete", "-f", crdDefinitionYaml)
		Log(t, "CRD", "Delete", deleteResult, deleteError)
	}
}

// Install performing installation using Helm when annotation overwrite as required.
// Returns success or failure.
func Install(t *testing.T, options *helm.Options, chartPath string, namespace string, releaseName string) bool {

	k8s.RunKubectlAndGetOutputE(t, options.KubectlOptions, "-n", namespace, "annotate", "--namespace", namespace, "--overwrite", "--all", "Deployment",
		"meta.helm.sh/release-name="+releaseName)

	err := helm.InstallE(t, options, chartPath, releaseName)
	if err != nil {
		logger.Log(t, fmt.Sprintf("Failed Helm install of release: %s using chartPath: %s", releaseName, chartPath))
		logger.Log(t, "Verify that local Kubernetes is running.")
		logger.Log(t, "Error condition: ", err)
	} else {
		k8s.WaitUntilAllNodesReady(t, options.KubectlOptions, 30, 2*time.Second)
		return true
	}
	return false
}

// Log helper for output of action details.
func Log(t *testing.T, subject string, action string, result string, errorDetail error) {

	if errorDetail != nil {
		logger.Log(t, fmt.Sprintf("[Subject: %s] [Action: %s] [Error: %s]", subject, action, errorDetail.Error()))
	} else {
		logger.Log(t, fmt.Sprintf("[Subject: %s] [Action: %s] [Result: %s]", subject, action, result))
	}
}

// LogResults helper for output of map-based result details
func LogResults(t *testing.T, results map[string]string) {
	for k, v := range results {
		logger.Log(t, fmt.Sprintf("Key: %s Value: %s", k, v))
	}
}

// Pause helper responsible for pausing for a specified time duration and interval.
func Pause(t *testing.T, message string, interval, timeout time.Duration) {

	wait.Poll(interval, timeout, func() (bool, error) {
		if true {
			logger.Log(t, fmt.Sprintf("Pausing: %s Timeout: %s", message, timeout))
		}
		return false, nil
	})
}

// CleanupDeployment performs lookup of deployment based on the namespace
// followed by delete of deployment specific to namespace.
func CleanupDeployment(t *testing.T, kubeOptions *k8s.KubectlOptions, chartName string) {

	lookupResult, lookupErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "deployment", "-A", "--field-selector",
		fmt.Sprintf("metadata.namespace=%s", kubeOptions.Namespace), "-o", "name")
	Log(t, "Deployment", "Delete", lookupResult, lookupErr)

	if lookupErr == nil && lookupResult == chartName {
		deleteResult, deleteErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "delete", "deployment", chartName, "--namespace", kubeOptions.Namespace)
		Log(t, "Deployment", "Delete", deleteResult, deleteErr)
	}

}
