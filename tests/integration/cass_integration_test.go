package tests

import (
	"strings"
	"sync"
	"testing"
	"time"

	util "./util"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"
)

var (
	chartName = "cass-operator"
	chartPath = "../../charts/k8ssandra/charts/cass-operator"
	// TODO: put your own config path here ...
	configPath        = "C:\\Users\\Jeff Banks\\.kube\\config"
	crdDefinitionYaml = "../../charts/k8ssandra/charts/cass-operator/crds/customresourcedefinition.yaml"
	crdResourceName   = "cassandradatacenters.cassandra.datastax.com"

	dc1Name          = "dc1"
	dc1Yaml          = "../../cassandra/example-cassdc-minimal.yaml"
	defaultNamespace = "default"

	namespacePrefix = "test-minimal-cassdc"
	releaseName     = "test-release-minimal-cassdc"
	testName        = "Helm Chart cass-cluster"
)

// Setup performs namespace lookup having test related resources to be cleaned.
func setup(t *testing.T) (string, *k8s.KubectlOptions, *helm.Options) {

	t.Parallel()
	namespace := util.GenerateNamespaceName(namespacePrefix)
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

	result := util.CreateTestNamespace(t, kubeOptions, releaseName, namespace)
	assert.NotNil(t, result)

	util.Install(t, helmOptions, chartPath, namespace, releaseName)
	k8s.WaitUntilPodAvailable(t, kubeOptions, dc1Name, 5, 2*time.Second)

	k8s.KubectlApply(t, kubeOptions, dc1Yaml)

	k8s.WaitUntilServiceAvailable(t, kubeOptions, "cass-operator", 3, time.Second*3)

	// teardown(t, helmOptions)
}

// Test for cleaning up manually as needed.
func TestCleanup(t *testing.T) {

	var cleanupWaiter sync.WaitGroup

	t.Parallel()
	namespace := "test-cleanup-1"
	kubeOptions := k8s.NewKubectlOptions("", "", namespace)
	helmOptions := &helm.Options{KubectlOptions: kubeOptions}

	cleanup(t, helmOptions, &cleanupWaiter)
	cleanupWaiter.Wait()
}

func cleanup(t *testing.T, helmOptions *helm.Options, waiter *sync.WaitGroup) {

	namespaces, err := util.GetNamespaces(t, helmOptions.KubectlOptions)
	for _, ns := range namespaces {

		// Must match our namespace pattern for it to be cleaned up.
		if strings.Contains(ns, namespacePrefix) || strings.Contains(ns, defaultNamespace) {

			namespace := strings.TrimPrefix(ns, "namespace/")
			kubeOptions := k8s.NewKubectlOptions("", "", namespace)

			waiter.Add(1)
			go func() {
				defer waiter.Done()
				util.CleanupDeployment(t, kubeOptions, chartName)
			}()

			waiter.Add(1)
			go func() {
				defer waiter.Done()
				util.CleanupCRD(t, kubeOptions, crdResourceName, crdDefinitionYaml)
			}()

			waiter.Add(1)
			go func() {
				defer waiter.Done()
				util.CleanupRelease(t, helmOptions, releaseName)
			}()

			if namespace != defaultNamespace {
				waiter.Add(1)
				go func() {
					defer waiter.Done()
					k8s.DeleteNamespace(t, kubeOptions, namespace)
				}()
			}
		}
	}
	util.Log(t, "Namespace", "lookup", strings.Join(namespaces, ","), err)
}
