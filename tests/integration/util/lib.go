package util

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/yaml"
)

/* Utility providing support primarily for k8ssandra integration tests */

// AnnotateResult model for release annotations applied
type AnnotateResult struct {
	IsReleaseNamespaceApplied bool
	IsReleaseNameApplied      bool
}

// CRDIdentity test model for Custom Resource Identifier
type CRDIdentity struct {
	CRDName         string
	CRDHookName     string
	CRDResourceName string
	CRDMetaName     string
}

// K8ssandraClusterIdentity test model for cluster
type K8ssandraClusterIdentity struct {
	ClusterName            string
	RenderedCassdcTemplate string
}

// CreateOperatorIdentity factory function for Operator Identity
func CreateOperatorIdentity() CRDIdentity {

	return CRDIdentity{
		CRDName:         "crd/cassandradatacenters.cassandra.datastax.com",
		CRDHookName:     "cass-operator-webhook",
		CRDResourceName: "cassandradatacenters.cassandra.datastax.com"}
}

// Annotate of release name, and release namespace, scoped to namespace
func Annotate(t *testing.T, helmOptions *helm.Options, releaseName string) *AnnotateResult {

	Log(t, "\n===== APPLY ANNOTATIONS", "", "", nil)
	namespace := helmOptions.KubectlOptions.Namespace
	result, errReleaseName := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "annotate", "-n", namespace, "--overwrite",
		"--all", "Deployment", "meta.helm.sh/release-name="+releaseName)
	Log(t, "Annotate", fmt.Sprintf("Overwrite with release name:%s for namespace:%s", releaseName, namespace), result, errReleaseName)

	resultRelNS, errReleaseNamespace := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "annotate", "-n", namespace, "--overwrite",
		"--all", "Deployment", "meta.helm.sh/release-namespace="+namespace)
	Log(t, "Annotate", fmt.Sprintf("Overwrite release-namespace with:%s", namespace), resultRelNS, errReleaseNamespace)
	return &AnnotateResult{IsReleaseNamespaceApplied: (errReleaseName == nil), IsReleaseNameApplied: (errReleaseNamespace == nil)}
}

// GenerateRandomNamespace provides a random namespace used in testing.
func GenerateRandomNamespace(namespacePrefix string) string {
	testID := strings.ToLower(random.UniqueId())
	name := fmt.Sprintf("%s-%s", namespacePrefix, testID)
	return name
}

//IsNamespaceExisting indicates if the target namespace exists.
func IsNamespaceExisting(t *testing.T, kubeOptions *k8s.KubectlOptions, namespace string) bool {

	var lookupErr error
	var namespaces []string

	namespaces, lookupErr = lookupNamespaces(t, kubeOptions)

	if lookupErr != nil {
		Log(t, "Namesace", "Lookup", "", lookupErr)
	}

	if namespaces != nil {
		for _, ns := range namespaces {
			trimmedNamespace := strings.TrimPrefix(ns, "namespace/")
			if namespace == trimmedNamespace {
				return true
			}
		}
	}
	return false
}

// DeleteRelease cleanup of release by name.
func DeleteRelease(t *testing.T, helmOptions *helm.Options, releaseName string) bool {

	Log(t, "\n===== DELETE RELEASE", "", "", nil)
	deleteErr := helm.DeleteE(t, helmOptions, releaseName, true)
	Log(t, "Release", fmt.Sprintf("Delete - %s", releaseName), "", deleteErr)
	return deleteErr == nil
}

// LookupPods provides current state of pods scoped by namespace
func LookupPods(t *testing.T, helmOptions *helm.Options, namespace string) ([]string, bool) {

	lookupResult, lookupErr := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "get", "pods", "-n", helmOptions.KubectlOptions.Namespace, "-o", "name")
	Log(t, "Get", "pods", lookupResult, lookupErr)
	return strings.Split(lookupResult, "\n"), lookupErr == nil
}

// DeletePodsByNamespace removal of pods scoped to helm option namespace
// Fails fast upon first delete error detected indicating failure
func DeletePodsByNamespace(t *testing.T, helmOptions *helm.Options) bool {

	namespace := helmOptions.KubectlOptions.Namespace
	pods, isError := LookupPods(t, helmOptions, namespace)
	if !isError {
		fmt.Println("pods found for namespace: ", namespace)
	}
	for _, pod := range pods {
		deleteResult, deleteErr := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "-n", namespace, "delete", pod)
		Log(t, "Delete", "pods", deleteResult, deleteErr)
		if deleteErr != nil {
			return false
		}
	}
	return true
}

// LookupCRDByName obtains current state of CRD by name
func LookupCRDByName(t *testing.T, helmOptions *helm.Options, name string) string {

	var lookupResult = ""
	result, err := k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "get", "crd", name, "-o", "name")
	if err != nil {
		Log(t, "Get", "crd", result, err)
	} else {
		lookupResult = result
	}
	return lookupResult
}

// DeleteNamespace removes specified namespace returning success or failure.
func DeleteNamespace(t *testing.T, helmOptions *helm.Options) bool {

	namespace := helmOptions.KubectlOptions.Namespace
	deleteErr := k8s.RunKubectlE(t, helmOptions.KubectlOptions, "delete", "namespace", namespace)
	Log(t, "Delete", "Namespace", fmt.Sprintf("Using ns: %s", namespace), deleteErr)
	return deleteErr == nil
}

// DeleteCRD performs lookup of resource targets specific to CRD resourc name.
// Removes the dependency resources and finalizers
func DeleteCRD(t *testing.T, kubeOptions *k8s.KubectlOptions, crdIdentity CRDIdentity, timeoutSeconds int) bool {

	Log(t, "\n===== DELETE CRD", "", "", nil)
	Log(t, "DeleteCRD", "lookup CRD by metaName: ", crdIdentity.CRDResourceName, nil)
	if lookupCRDByMetaName(t, kubeOptions, crdIdentity.CRDResourceName) {

		errResults := deleteChildResources(t, kubeOptions, getChildResources(), crdIdentity.CRDHookName)
		if len(errResults) != 0 {
			LogResults(t, "Delete CRD child resource removal failures.", errResults)
			return false
		}

		isFinalizersRemoved := removeFinalizers(t, kubeOptions, crdIdentity.CRDName)
		if !isFinalizersRemoved {
			Log(t, "Finalizers NOT removed.", "removeFinalizers", "", nil)
			return false
		}

		isCRDDeleted := deleteCustomResourceDefinition(t, kubeOptions, crdIdentity.CRDName, timeoutSeconds)
		if !isCRDDeleted {
			Log(t, "CRD NOT deleted.", "deleteCustomResourceDefinition", "", nil)
			return false
		}
	}
	return true
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
func LogResults(t *testing.T, message string, results map[string]string) {
	for k, v := range results {
		logger.Log(t, fmt.Sprintf("%s Id: %s Result: %s", message, k, v))
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

// DeleteDeployment performs lookup of deployment based on the namespace
// followed by delete of deployment specific to namespace.
func DeleteDeployment(t *testing.T, helmOptions *helm.Options, chartName string) {

	Log(t, "\n=====  DELETE DEPLOYMENT", "", "", nil)
	var lookupErr error
	var lookupResult string
	namespace := helmOptions.KubectlOptions.Namespace

	lookupResult, lookupErr = k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "get", "deployment", "-A", "--field-selector",
		fmt.Sprintf("metadata.namespace=%s", namespace), "-o", "name")
	Log(t, "Deployment", "Delete", lookupResult, lookupErr)

	if lookupErr == nil && lookupResult == chartName {

		var deleteErr error
		var deleteResult string

		deleteResult, deleteErr = k8s.RunKubectlAndGetOutputE(t, helmOptions.KubectlOptions, "delete", "deployment",
			chartName, "--namespace", namespace)
		Log(t, "Deployment", "Delete", deleteResult, deleteErr)
	}
}

// CreateNamespace creates specific test based namespace
// as defined by namespace provided in options.
func CreateNamespace(t *testing.T, helmOptions *helm.Options) {
	namespace := helmOptions.KubectlOptions.Namespace
	k8s.CreateNamespaceE(t, helmOptions.KubectlOptions, namespace)
}

// createK8ssandraClusterIdentity performs construction of cluster identity for k8ssandra
// includes checks for required chart used along with associated rendered template.
func createK8ssandraClusterIdentity(t *testing.T, kubeOptions *k8s.KubectlOptions,
	clusterName string, releaseName string) K8ssandraClusterIdentity {

	var renderedMap map[string]interface{}
	helmChartPath, err := filepath.Abs("../../../charts/k8ssandra-cluster")
	require.NoError(t, err)

	options := &helm.Options{
		SetStrValues:   map[string]string{"name": kubeOptions.Namespace, "clusterName": clusterName},
		KubectlOptions: kubeOptions,
	}

	renderedOutput := helm.RenderTemplate(
		t, options, helmChartPath, releaseName,
		[]string{"templates/cassdc.yaml"},
	)
	jsonRendered, err := yaml.YAMLToJSON([]byte(renderedOutput))
	require.NoError(t, json.Unmarshal(jsonRendered, &renderedMap))

	return K8ssandraClusterIdentity{
		ClusterName:            clusterName,
		RenderedCassdcTemplate: renderedOutput,
	}
}

// lookupNamespaces provides a string array list of namspaces found during test.
func lookupNamespaces(t *testing.T, kubeOptions *k8s.KubectlOptions) ([]string, error) {

	var lookupErr error
	var lookupResult string

	lookupResult, lookupErr = k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "namespace", "-o", "name")
	return strings.Split(lookupResult, "\n"), lookupErr
}

// deleteTokens for a given set of tokens, remove each and short-circuit upon failure.
func deleteTokens(t *testing.T, kubeOptions *k8s.KubectlOptions, tokens []string) error {

	for _, token := range tokens {
		deleteResult, deleteErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "--namespace", kubeOptions.Namespace, "delete", token)
		Log(t, "Secret", "Delete", deleteResult, deleteErr)
		if deleteErr != nil {
			return deleteErr
		}
	}
	return nil
}

func removeFinalizers(t *testing.T, kubeOptions *k8s.KubectlOptions, customResource string) bool {

	patchResult, patchErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions, "patch",
		customResource, "-p", "{\"metadata\":{\"finalizers\":[]}}", "--type=merge")
	Log(t, "Finalizers", "Remove", patchResult, patchErr)
	return patchErr == nil
}

func lookupSecretNamespaceScoped(t *testing.T, kubeOptions *k8s.KubectlOptions) (string, error) {

	return k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", "Secret",
		"-A", "--field-selector", fmt.Sprintf("metadata.namespace=%s", kubeOptions.Namespace), "-o", "name")
}

// DeleteResources performs cleanup of target resources provided unique to namespace and/or resourceName.
// Returns a map of any errors the target list provided.
func deleteChildResources(t *testing.T, kubeOptions *k8s.KubectlOptions,
	targets []string, metaName string) map[string]string {

	errorResults := map[string]string{}

	for _, target := range targets {
		// Special handling for secret tokens
		if strings.ToLower(target) == "secret" {

			lookupSecret, lookupSecretErr := lookupSecretNamespaceScoped(t, kubeOptions)
			if lookupSecretErr != nil {
				errorResults[target] = fmt.Sprintf("Get Secret returned error: %s", lookupSecretErr.Error())
				continue
			}

			// Delete tokens for the Secret
			deleteTokenErr := deleteTokens(t, kubeOptions, strings.Split(lookupSecret, "\n"))
			if deleteTokenErr != nil {
				errorResults[target] = fmt.Sprintf("Delete Secret returned error: %s", deleteTokenErr.Error())
			}

		} else {

			var lookupResult string
			var lookupErr error

			if strings.ToLower(target) == "validatingwebhookconfiguration" {
				lookupResult, lookupErr = lookupResource(t, kubeOptions, target, "cassandradatacenter-webhook-registration")
			} else {
				lookupResult, lookupErr = lookupResource(t, kubeOptions, target, metaName)
			}

			if lookupErr != nil {
				errorResults[target] = fmt.Sprintf("Get of %s returned error: %s", target, lookupErr.Error())
				continue
			}

			// No error detected, and have a resource identified, officially delete it.
			if lookupResult != "" && lookupResult != "No resources found" {

				deleteErr := k8s.RunKubectlE(t, kubeOptions, "delete", "-n", kubeOptions.Namespace, lookupResult)
				if deleteErr != nil {
					errorResults[target] = fmt.Sprintf("Delete returned error: %s", lookupErr.Error())
					continue
				}
			}
		}
	}
	return errorResults
}

// lookupResource performs lookup of a target resource scoped by metadata.name
func lookupResource(t *testing.T, kubeOptions *k8s.KubectlOptions, target string, resourceName string) (string, error) {

	return k8s.RunKubectlAndGetOutputE(t, kubeOptions, "get", target, "-A", "--field-selector",
		fmt.Sprintf("metadata.name=%s", resourceName), "--namespace", kubeOptions.Namespace, "-o", "name")
}

// lookupCRDByMetaName for metadata scoped name.
func lookupCRDByMetaName(t *testing.T, kubeOptions *k8s.KubectlOptions, metadataName string) bool {

	Log(t, "CRD", "Lookup", fmt.Sprintf(" using namespace: %s and metadataName: %s", kubeOptions.Namespace, metadataName), nil)
	lookupResult, lookupErr := k8s.RunKubectlAndGetOutputE(t, kubeOptions,
		"get", "crd", "-A", "--field-selector",
		fmt.Sprintf("metadata.name=%s", metadataName),
		"--namespace", kubeOptions.Namespace, "-o", "name")
	Log(t, "CRD", "Lookup", lookupResult, lookupErr)
	return lookupResult != "" && lookupErr == nil
}

// getChildResources provides a known set of resources.
func getChildResources() []string {
	return []string{
		"ClusterRoleBinding",
		"ClusterRole",
		"ValidatingWebhookConfiguration",
		"MutatingWebhookConfiguration",
		"Secret",
		"Service",
		"ValidatingWebhookConfiguration",
	}
}

//deleteCustomResourceDefinition removes CRD for resource def provided.
func deleteCustomResourceDefinition(t *testing.T, kubeOptions *k8s.KubectlOptions,
	crdResourceName string, timeoutSeconds int) bool {

	deleteErr := k8s.RunKubectlE(t, kubeOptions,
		"delete", crdResourceName, fmt.Sprintf("--timeout=%ds", timeoutSeconds))
	return deleteErr == nil
}
