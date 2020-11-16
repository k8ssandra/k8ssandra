package integration

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/k8ssandra/k8ssandra/tests/integration/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	chartName         = "cass-operator"
	chartPath         = "../../charts/k8ssandra/charts/cass-operator"
	configPath        = "C:\\Users\\Jeff Banks\\.kube\\config"
	crdDefinitionYaml = "../../charts/k8ssandra/charts/cass-operator/crds/customresourcedefinition.yaml"
	crdResourceName   = "cassandradatacenters.cassandra.datastax.com"

	dc1Name          = "dc1"
	dc1Yaml          = "../../cassandra/example-cassdc-minimal.yaml"
	defaultNamespace = "default"

	namespacePrefix = "test-ns-cassop"
	releaseName     = "test-release-cassop"
	testName        = "Helm Chart cass-cluster"
)

// Setup performs namespace lookup having test related resources to be cleaned.
func setup(t *testing.T) (string, *k8s.KubectlOptions, *helm.Options) {

	t.Parallel()
	var setupWaiter sync.WaitGroup
	var namespace string

	setupWaiter.Add(1)
	go func() {
		defer setupWaiter.Done()
		namespace = util.GenerateNamespaceName(namespacePrefix)
	}()
	setupWaiter.Wait()

	kubeOptions := k8s.NewKubectlOptions("", "", namespace)
	helmOptions := &helm.Options{KubectlOptions: kubeOptions}

	var cleanupWaiter sync.WaitGroup
	cleanup(t, helmOptions, &cleanupWaiter)
	cleanupWaiter.Wait()

	return namespace, kubeOptions, helmOptions
}

//
func teardown(t *testing.T, helmOptions *helm.Options) {
	var cleanupWaiter sync.WaitGroup
	cleanup(t, helmOptions, &cleanupWaiter)
	cleanupWaiter.Wait()
}

// TestCassOperator performs basic installation of cass-operator.
func TestCassOperator(t *testing.T) {

	// Setup
	namespace, kubeOptions, helmOptions := setup(t)

	assert.NotNil(t, namespace)
	assert.NotNil(t, kubeOptions)
	assert.NotNil(t, helmOptions)

	util.CreateNamespace(t, helmOptions)
	util.ApplyAnnotation(t, helmOptions, releaseName)
	
	isInstalled := installCassOperator(t, helmOptions)
	require.True(t, isInstalled)

	lookupResult := util.LookupCRDByName(t, helmOptions, "cassandradatacenters.cassandra.datastax.com")
	require.Equal(t, "customresourcedefinition.apiextensions.k8s.io/cassandradatacenters.cassandra.datastax.com", lookupResult)
}

// installCassOperator provides apply and returned result of cass operator
// based on selector provided by caller.
func installCassOperator(t *testing.T, helmOptions *helm.Options) bool {

	var applyWaiter sync.WaitGroup
	var installErr error
	namespace := helmOptions.KubectlOptions.Namespace

	applyWaiter.Add(1)

	go func() {
		defer applyWaiter.Done()
		installErr = k8s.KubectlApplyE(t, helmOptions.KubectlOptions, crdDefinitionYaml)
		util.Log(t, "K8s", fmt.Sprintf("Install cass-operator namespace:%s release:%s", namespace, releaseName), "", installErr)
	}()
	applyWaiter.Wait()
	return (installErr == nil)
}

// Cleanup testing artifacts scoped to default and test namespaces.
// Utilizes a wait group for caller coordination with cleanup activities.
func cleanup(t *testing.T, helmOptions *helm.Options, waiter *sync.WaitGroup) {

	namespaces, err := util.GetNamespaces(t, helmOptions.KubectlOptions)
	util.Log(t, "Namespace", "lookup", strings.Join(namespaces, ","), err)

	for _, ns := range namespaces {

		// Must match our namespace pattern for it to be cleaned up.
		if strings.Contains(ns, namespacePrefix) {

			namespace := strings.TrimPrefix(ns, "namespace/")
			kubeOptions := k8s.NewKubectlOptions("", "", namespace)
			helmOptions := &helm.Options{KubectlOptions: kubeOptions}

			util.ApplyAnnotation(t, helmOptions, releaseName)
			util.Log(t, "Processing", fmt.Sprintf("Cleanup Deployment ns:%s", namespace), "", nil)
			waiter.Add(1)
			go func() {
				defer waiter.Done()
				util.CleanupDeployment(t, kubeOptions, chartName)
			}()

			util.Log(t, "Processing", fmt.Sprintf("Cleanup CRD ns:%s", namespace), "", nil)
			waiter.Add(1)
			go func() {
				defer waiter.Done()
				util.CleanupCRD(t, kubeOptions, crdResourceName, crdDefinitionYaml)
			}()

			util.Log(t, "Processing", fmt.Sprintf("Cleanup Release ns:%s", namespace), "", nil)
			waiter.Add(1)
			go func() {
				defer waiter.Done()
				util.CleanupRelease(t, helmOptions, releaseName)
			}()

			if namespace != defaultNamespace {
				util.Log(t, "Processing", fmt.Sprintf("Delete Namespace ns:%s", namespace), "", nil)
				waiter.Add(1)
				go func() {
					defer waiter.Done()
					k8s.DeleteNamespaceE(t, kubeOptions, namespace)
				}()
			}
		}
	}
}
