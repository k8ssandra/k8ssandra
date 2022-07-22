package steps

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/k8ssandra/cass-operator/tests/util/kubectl"
	. "github.com/onsi/gomega"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func DeployK8ssandraOperatorCluster(t *testing.T, namespace string, k8ssandraCluster string, upgrade bool, useLocalCharts bool, version string) {
	operatorChartPath, err := filepath.Abs("../../charts/k8ssandra-operator")
	g(t).Expect(err).To(BeNil())
	clusterPath, err := filepath.Abs("../../tests/integration/k8ssandra-clusters/" + k8ssandraCluster)
	g(t).Expect(err).To(BeNil())

	if !useLocalCharts {
		installK8ssandraHelmRepo(t)
		operatorChartPath = "k8ssandra/k8ssandra-operator"
	}

	helmOptions := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
	}

	if version != "" && version != "latest" {
		g(t).Expect(useLocalCharts).To(BeFalse(), "K8ssandra version can only be passed when using Helm repo based installs, not local charts.")
		helmOptions.Version = version
	}

	defer timeTrack(time.Now(), "Installing and starting k8ssandra-operator")
	if upgrade {
		initialResourceVersion := cassDcResourceVersion(t, namespace)
		err = helm.UpgradeE(t, helmOptions, operatorChartPath, releaseName)
		g(t).Expect(err).To(BeNil(), "Failed installing k8ssandra with Helm: %v", err)
		waitForCassDcUpgrade(t, namespace, initialResourceVersion)
	} else {
		err = helm.InstallE(t, helmOptions, operatorChartPath, releaseName)
		g(t).Expect(err).To(BeNil(), "Failed installing k8ssandra with Helm: %v", err)
	}
	g(t).Eventually(func() bool {
		stdout, stderr, _ := kubectl.ApplyFiles(clusterPath).InNamespace(namespace).ExecVCapture()
		println(stdout)
		println(stderr)
		return stderr == ""
	}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 5).Should(BeTrue())

	// Wait for CassandraDatacenter to be ready..
	WaitForCassDcToBeReady(t, namespace)
}

func DeleteK8ssandraCluster(t *testing.T, namespace string, k8ssandraCluster string) {
	stdout, stderr, err := kubectl.
		DeleteFromFiles("../../tests/integration/k8ssandra-clusters/" + k8ssandraCluster).
		InNamespace(namespace).
		ExecVCapture()
	println(stdout)
	println(stderr)
	if err != nil {
		println(err)
		t.FailNow()
	}
	WaitForCassandraDatacenterDeletion(t, namespace)
}
