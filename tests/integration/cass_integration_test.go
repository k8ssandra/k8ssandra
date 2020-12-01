package integration

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/k8ssandra/k8ssandra/tests/integration/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	chartName  = "cass-operator"
	chartPath  = "../../charts/k8ssandra/charts/cass-operator"
	configPath = "C:\\Users\\Jeff Banks\\.kube\\config"

	crdClusterYaml = "../../charts/k8ssandra-cluster/Chart.yaml"

	defaultNamespace        = "default"
	deleteCRDTimeoutSeconds = 15

	k8ssandraNamespace  = "k8ssandra-ns"
	operatorReleaseName = "k8ssandra-1"
	clusterReleaseName  = "cluster-1"
)

// TestCassOperator performs basic installation of k8ssandra operator and c* cluster using setup & teardown.
func TestCassOperatorAndClusterInstall(t *testing.T) {

	helmOptions := setup(t, operatorReleaseName)

	fmt.Println("helmOptions namespace is", helmOptions.KubectlOptions.Namespace)
	require.NotNil(t, helmOptions)

	// cass-operator install
	util.Annotate(t, helmOptions, operatorReleaseName)
	isOperatorInstalled := install(t, helmOptions, operatorReleaseName, "k8ssandra/k8ssandra")

	require.True(t, isOperatorInstalled)
	assert.True(t, repoUpdate(t, helmOptions))

	// cluster install
	isClusterInstalled := install(t, helmOptions, clusterReleaseName, "k8ssandra/k8ssandra-cluster")
	require.True(t, isClusterInstalled)
	assert.True(t, repoUpdate(t, helmOptions))

	// verification of crd existence
	lookupResult := util.LookupCRDByName(t, helmOptions, "cassandradatacenters.cassandra.datastax.com")
	require.Equal(t, "customresourcedefinition.apiextensions.k8s.io/cassandradatacenters.cassandra.datastax.com",
		lookupResult)

}

// Setup performs namespace setup w/ applied annotations for release.
// Verifies repository addition and updates are in-place prior to running test ops.
func setup(t *testing.T, releaseName string) *helm.Options {

	kubeOptions := k8s.NewKubectlOptions("", "", k8ssandraNamespace)
	helmOptions := &helm.Options{KubectlOptions: kubeOptions}
	util.CreateNamespace(t, helmOptions)

	cleanup(t, helmOptions, releaseName)
	assert.NotNil(t, kubeOptions)
	assert.NotNil(t, helmOptions)

	assert.True(t, repoAdd(t, helmOptions, "k8ssandra", "https://helm.k8ssandra.io/"))
	assert.True(t, repoUpdate(t, helmOptions))

	return helmOptions
}

func repoAdd(t *testing.T, helmOptions *helm.Options, name string, url string) bool {

	var addErr error
	addResult, addErr := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "repo", "add", name, url)
	util.Log(t, "Repo", "Add", addResult, addErr)
	return (addErr == nil)
}

func repoUpdate(t *testing.T, helmOptions *helm.Options) bool {

	var updateErr error
	updateResult, updateErr := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "repo", "update")
	util.Log(t, "Repo", "Update", updateResult, updateErr)
	return (updateErr == nil)
}

func uninstall(t *testing.T, helmOptions *helm.Options, releaseName string) bool {

	var uninstallErr error
	uninstallResult, uninstallErr := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "uninstall", releaseName, "--no-hooks")
	util.Log(t, "Helm", "Uninstall", uninstallResult, uninstallErr)
	return (uninstallErr == nil)
}

func install(t *testing.T, helmOptions *helm.Options, releaseName string, chartName string) bool {

	var installErr error
	installResult, installErr := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "install", "-n", k8ssandraNamespace,
		releaseName, chartName, "--create-namespace", "--insecure-skip-tls-verify")

	util.Log(t, "Helm", "Install", installResult, installErr)
	return (installErr == nil)
}

// Cleanup testing artifacts scoped to default and test namespaces.
func cleanup(t *testing.T, helmOptions *helm.Options, releaseName string) {

	if util.IsNamespaceExisting(t, helmOptions.KubectlOptions, k8ssandraNamespace) {

		kubeOptions := k8s.NewKubectlOptions("", "", k8ssandraNamespace)
		helmOptions := &helm.Options{KubectlOptions: kubeOptions}

		util.Annotate(t, helmOptions, releaseName)
		util.DeleteDeployment(t, helmOptions, chartName)
		util.DeleteCRD(t, kubeOptions, util.CreateOperatorIdentity(), deleteCRDTimeoutSeconds)
		util.DeleteRelease(t, helmOptions, releaseName)
		util.DeletePodsByNamespace(t, helmOptions)
		util.DeleteNamespace(t, helmOptions)
	}

	uninstall(t, helmOptions, releaseName)
}
