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
	clusterName         = "dc1"
)

// TestCassOperator performs basic installation of k8ssandra operator and c* cluster using setup & teardown.
func TestCassOperatorAndClusterInstall(t *testing.T) {

	kubeOptions := k8s.NewKubectlOptions("", "", k8ssandraNamespace)
	helmOptions := &helm.Options{KubectlOptions: kubeOptions}

	fmt.Println("helmOptions namespace is", helmOptions.KubectlOptions.Namespace)
	require.NotNil(t, helmOptions)

	operatorInstallPreconditions(t, helmOptions, operatorReleaseName)

	util.Annotate(t, helmOptions, operatorReleaseName)

	// cass-operator install
	isOperatorInstalled := install(t, helmOptions, operatorReleaseName, "k8ssandra/k8ssandra")
	require.True(t, isOperatorInstalled)
	assert.True(t, repoUpdate(t, helmOptions))

	// verification of crd existence
	lookupResult := util.LookupCRDByName(t, helmOptions, "cassandradatacenters.cassandra.datastax.com")
	require.Equal(t, "customresourcedefinition.apiextensions.k8s.io/cassandradatacenters.cassandra.datastax.com",
		lookupResult)

	// cluster install
	clusterInstallPreconditions(t, helmOptions, clusterReleaseName)
	util.Annotate(t, helmOptions, clusterReleaseName)
	isClusterInstalled := install(t, helmOptions, clusterReleaseName, "k8ssandra/k8ssandra-cluster")
	require.True(t, isClusterInstalled)
	assert.True(t, repoUpdate(t, helmOptions))

}

// clusterInstallPreconditions provides test cleanup and preconditions prior to test function execution
func clusterInstallPreconditions(t *testing.T, helmOptions *helm.Options, releaseName string) {

	assert.NotNil(t, helmOptions.KubectlOptions)
	assert.NotNil(t, helmOptions)

	util.CreateNamespace(t, helmOptions)
	clusterIdentity := util.CreateClusterIdentity(t, helmOptions.KubectlOptions, clusterName, releaseName)

	cleanupCluster(t, helmOptions, releaseName, clusterIdentity)

	repoAdd(t, helmOptions, "k8ssandra", "https://helm.k8ssandra.io/")
	repoUpdate(t, helmOptions)
}

// operatorInstallPreconditions provides test cleanup and preconditions prior to test function execution
func operatorInstallPreconditions(t *testing.T, helmOptions *helm.Options, releaseName string) {

	assert.NotNil(t, helmOptions.KubectlOptions)
	assert.NotNil(t, helmOptions)

	util.CreateNamespace(t, helmOptions)

	cleanupOperator(t, helmOptions, releaseName, util.CreateOperatorIdentity())

	repoAdd(t, helmOptions, "k8ssandra", "https://helm.k8ssandra.io/")
	repoUpdate(t, helmOptions)
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

// cleanupOperator cleanup of operator specific resources.
func cleanupOperator(t *testing.T, helmOptions *helm.Options, releaseName string, identity util.CRDIdentity) {

	if util.IsNamespaceExisting(t, helmOptions.KubectlOptions, k8ssandraNamespace) {

		kubeOptions := k8s.NewKubectlOptions("", "", k8ssandraNamespace)
		helmOptions := &helm.Options{KubectlOptions: kubeOptions}

		util.Annotate(t, helmOptions, releaseName)
		util.DeleteDeployment(t, helmOptions, chartName)
		util.DeleteOperatorCRD(t, kubeOptions, identity, deleteCRDTimeoutSeconds)
		util.DeleteRelease(t, helmOptions, releaseName)
		util.DeletePodsByNamespace(t, helmOptions)
		util.DeleteNamespace(t, helmOptions)
	}
	uninstall(t, helmOptions, releaseName)
}

// cleanupCluster cleanup of cluster specific resources.
func cleanupCluster(t *testing.T, helmOptions *helm.Options, releaseName string, identity util.ClusterIdentity) {

	if util.IsNamespaceExisting(t, helmOptions.KubectlOptions, k8ssandraNamespace) {

		kubeOptions := k8s.NewKubectlOptions("", "", k8ssandraNamespace)
		helmOptions := &helm.Options{KubectlOptions: kubeOptions}

		util.Annotate(t, helmOptions, releaseName)
		util.DeleteDeployment(t, helmOptions, chartName)
		util.DeleteClusterCRD(t, kubeOptions, identity, deleteCRDTimeoutSeconds)
		util.DeleteRelease(t, helmOptions, releaseName)
		util.DeletePodsByNamespace(t, helmOptions)
		util.DeleteNamespace(t, helmOptions)
	}
	uninstall(t, helmOptions, releaseName)
}
