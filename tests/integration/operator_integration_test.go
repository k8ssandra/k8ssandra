package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/k8ssandra/k8ssandra/tests/integration/util"
)

var (
	testSuite          = "k8ssandra-integration"
	releaseName        = "k8ssandra"
	namespace          = "k8ssandra"
	operatorChart      = "../../charts/k8ssandra"
	clusterChart       = "../../charts/k8ssandra-cluster"
	clusterReleaseName = "cassdc"
	deployedStatus     = "deployed"
)

func setup(t *testing.T) Options {

	kubeOptions := k8s.NewKubectlOptions("", "", namespace)
	operatorOptions := &helm.Options{KubectlOptions: kubeOptions}
	clusterOptions := &helm.Options{KubectlOptions: kubeOptions}

	UninstallRelease(t, clusterOptions, clusterReleaseName)
	UninstallRelease(t, operatorOptions, releaseName)
	return Options{Cluster: clusterOptions, Operator: operatorOptions}

}

func TestOperatorAndClusterInstall(t *testing.T) {

	options := setup(t)

	InstallChart(t, options.Operator, operatorChart, releaseName)
	var releases = LookupReleases(t, options.Operator)
	require.Len(t, releases, 1, "Expected a single k8ssandra release entry.")

	var release = releases[0]
	require.Equal(t, release.Name, releaseName)
	require.Equal(t, release.Status, deployedStatus)
	require.Equal(t, release.Namespace, namespace)

	// k8ssandra-cluster requires that the previous operator is installed.
	InstallChart(t, options.Cluster, clusterChart, clusterReleaseName)
	releases = LookupReleases(t, options.Cluster)
	require.Len(t, releases, 2, "Expected two k8ssandra namespaced release entries.")

	for _, rel := range releases {
		if rel.Name == clusterReleaseName {
			require.True(t, rel.Status == deployedStatus && rel.Namespace == namespace)
		}
	}
}

func TestK8ssandraSetup(t *testing.T) {

	options := setup(t)
	require.NotNil(t, options)
	require.NotNil(t, options.Cluster)
	require.NotNil(t, options.Operator)

	var releases = LookupReleases(t, options.Operator)
	require.Len(t, releases, 0, "Expected there is not a k8ssandra release entry.")
}
