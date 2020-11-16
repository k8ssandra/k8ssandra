package util

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

/* Utility providing support primarily for k8ssandra integration tests */

// AnnotateResult model for release annotations applied
type AnnotateResult struct {
	IsReleaseNamespaceApplied bool
	IsReleaseNameApplied      bool
}

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
	helmDeleteReleaseError := helm.DeleteE(t, helmOptions, releaseName, true)
	Log(t, "Release Name", fmt.Sprintf("Delete - %s", releaseName), "", helmDeleteReleaseError)
}

// ApplyAnnotation applys annotation specifc overwrites wrapped by waiter
func ApplyAnnotation(t *testing.T, helmOptions *helm.Options, releaseName string) {

	var applyAnnotationWaiter sync.WaitGroup
	applyAnnotationWaiter.Add(1)
	go func() {
		defer applyAnnotationWaiter.Done()
		annotate(t, helmOptions, releaseName)
	}()
	applyAnnotationWaiter.Wait()
}

// CreateNamespace creates specific test based namespace
// as defined by namespace provided in options.
func CreateNamespace(t *testing.T, helmOptions *helm.Options) {

	var createNamespaceWaiter sync.WaitGroup
	namespace := helmOptions.KubectlOptions.Namespace

	createNamespaceWaiter.Add(1)
	go func() {
		defer createNamespaceWaiter.Done()
		assert.Nil(t, k8s.CreateNamespaceE(t, helmOptions.KubectlOptions, namespace))
	}()
	createNamespaceWaiter.Wait()
}

// LookupCRDByName obtains current state of CRD by name
func LookupCRDByName(t *testing.T, helmOptions *helm.Options, name string) string {

	var crdLookupWaiter sync.WaitGroup
	var lookupResult = ""
	crdLookupWaiter.Add(1)

	go func() {

		defer crdLookupWaiter.Done()
		result, err := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "get", "crd", name, "-o", "name")

		if err != nil {
			Log(t, "GET", "CRD", "", err)
		} else {
			lookupResult = result
		}
	}()
	crdLookupWaiter.Wait()
	return lookupResult
}

// CleanupCRD performs lookup of resource targets specific to CRD metatdata name.
// Removes the dependency resources followed by the CRD itself.
func CleanupCRD(t *testing.T, kubeOptions *k8s.KubectlOptions, metadataName string, crdDefinitionYaml string) {

	lookupResult, lookupError := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "crd", "-A", "--field-selector",
		fmt.Sprintf("metadata.name=%s", metadataName), "--namespace", kubeOptions.Namespace, "-o", "name")
	Pause(t, "Looking up existing CRD", 1*time.Second, 2*time.Second)

	if lookupError == nil && lookupResult != "" {
		resources := []string{
			"ClusterRoleBinding",
			"ClusterRole",
			"ValidatingWebhookConfiguration",
			"MutatingWebhookConfiguration",
			"Secret",
			"Service",
		}

		Log(t, "CRD", "Lookup", lookupResult, lookupError)
		results := deleteResources(t, kubeOptions, resources, kubeOptions.Namespace, "cass-operator-webhook")
		Pause(t, "Deleting resources before clean up of CRD.", 1*time.Second, 3*time.Second)
		LogResults(t, results)

		patchResult, patchError := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "patch",
			"crd/cassandradatacenters.cassandra.datastax.com",
			"-p", "{\"metadata\":{\"finalizers\":[]}}", "--type=merge")
		Log(t, "CRD", "Patch", patchResult, patchError)
		Pause(t, "Pause (paws) for patch", 1*time.Second, 3*time.Second)

		deleteResult, deleteError := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "delete", "-f", crdDefinitionYaml, "--timeout=10s")
		Log(t, "CRD", "Delete", deleteResult, deleteError)
		Pause(t, "Pause for Delete CRD", 7*time.Second, 8*time.Second)

	} else {
		Log(t, "CRD", "Lookup of CRD not located.", lookupResult, lookupError)
	}
}

// Log helper for output of action details.
func Log(t *testing.T, subject string, action string, result string, errorDetail error) {

	if errorDetail != nil {
		logger.Log(t, fmt.Sprintf("[ %s ] [ %s] [ Error: %s]", subject, action, errorDetail.Error()))
	} else {
		if result != "" {
			logger.Log(t, fmt.Sprintf("[ %s ] [ %s ] [ %s]", subject, action, result))
		} else {
			logger.Log(t, fmt.Sprintf("[ %s ] [ %s ]", subject, action))
		}
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
	Pause(t, "Getting deployment -A", 1*time.Second, 3*time.Second)
	Log(t, "Deployment", "Delete", lookupResult, lookupErr)

	if lookupErr == nil && lookupResult == chartName {
		deleteResult, deleteErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "delete", "deployment", chartName, "--namespace", kubeOptions.Namespace)
		Log(t, "Deployment", "Delete", deleteResult, deleteErr)
	}
}

// Annotate for release name, and release namespace, scoped to namespace
func annotate(t *testing.T, helmOptions *helm.Options, releaseName string) *AnnotateResult {

	namespace := helmOptions.KubectlOptions.Namespace
	result, errReleaseName := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "annotate", "-n", namespace, "--overwrite",
		"--all", "Deployment", "meta.helm.sh/release-name="+releaseName)
	Log(t, "Annotate", fmt.Sprintf("Overwrite with release name:%s for namespace:%s", releaseName, namespace), result, errReleaseName)

	resultRelNS, errReleaseNamespace := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "annotate", "-n", namespace, "--overwrite",
		"--all", "Deployment", "meta.helm.sh/release-namespace="+namespace)
	Log(t, "Annotate", fmt.Sprintf("Overwrite with release namespace:%s for namespace:%s", namespace, namespace), resultRelNS, errReleaseNamespace)

	return &AnnotateResult{IsReleaseNamespaceApplied: (errReleaseName == nil), IsReleaseNameApplied: (errReleaseNamespace == nil)}
}

// DeleteResources performs cleanup of target resources provided unique to namespace and/or resourceName.
// Returns a map of delete resource results and/or errors for the target list provided.
func deleteResources(t *testing.T, kubeOptions *k8s.KubectlOptions, targets []string, namespace string, resourceName string) map[string]string {

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

		// Once no error detected, and have a resource identified, officially delete it.
		if lookup != "" && lookup != "No resources found" {
			noNamespaceOptions := k8s.NewKubectlOptions("", "", "")
			deleteResult, deleteErr := k8s.RunKubectlAndGetOutputE(t, noNamespaceOptions, "delete", lookup)
			Pause(t, "Pausing for Resource to be deleted", 1*time.Second, 2*time.Second)

			if deleteErr != nil {
				resultMap[target] = fmt.Sprintf(" 'delete' returned error: %s", err.Error())
				continue
			}
			resultMap[target] = deleteResult
		}
	}
	return resultMap
}
