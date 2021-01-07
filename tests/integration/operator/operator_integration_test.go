package operator

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/k8ssandra/k8ssandra/tests/integration/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "k8s.io/api/apps/v1"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/types"
	"testing"
	"time"
)

var (
	suiteName            = "k8ssandra-operator-integration-suite"
	k8ssandraClusterName = "k8ssandra-cluster"
	k8ssandraContext     = "kind-k8ssandra-cluster"
	operatorReleaseName  = "k8ssandra"
	clusterReleaseName   = "cassdc"
	networkReleaseName   = "traefik"
	namespace            = "k8ssandra"

	operatorChart  = "../../../charts/k8ssandra"
	clusterChart   = "../../../charts/k8ssandra-cluster"
	clusterValues  = "../preconditions/k8ssandra-cluster-values.yaml"
	traefikValues  = "../preconditions/k8ssandra-traefik-values.yaml"
	kubeConfigFile = "../preconditions/k8ssandra-config-test.yaml"
	kindConfigFile = "../preconditions/k8ssandra-kind-config.yaml"
	kindImage      = "kindest/node:v1.18.2"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, suiteName)
}

// setup used for creating test option context and cleanup of operator and cluster.
func setup() K8ssandraOptions {

	options := CreateK8ssandraOptions(k8s.NewKubectlOptions(k8ssandraContext,
		"", namespace), clusterReleaseName, operatorReleaseName, networkReleaseName)
	ConfigureContext(options.GetClusterCtx(), kubeConfigFile)

	// Check if kind is available.
	if IsClusterExisting(k8ssandraClusterName) {
		DeleteKindCluster(k8ssandraClusterName)
	}
	CreateKindCluster(KindClusterDetail{Name: k8ssandraClusterName,
		Image: kindImage, ConfigFile: kindConfigFile})
	RepoUpdate(options.GetOperatorCtx())
	ConfigureContext(options.GetClusterCtx(), kubeConfigFile)
	return options
}

var _ = Describe(suiteName, func() {

	Describe("Operator and cluster integration", func() {

		var ko K8ssandraOptions
		BeforeEach(func() {
			ko = setup()
			By("Expecting Traefik to be deployed")
			InstallChartWithValues(ko.GetNetworkCtx(), "traefik/traefik", traefikValues)
			Expect(IsReleaseDeployed(ko.GetNetworkCtx())).Should(BeTrue())
			k8s.WaitUntilAllNodesReady(GinkgoT(), ko.GetNetworkCtx().KubeOptions, 15, 3*time.Second)
		})

		AfterEach(func() {
			Log("Post-Test", "Timestamp",
				fmt.Sprintf("Test run end: %s", time.Now().String()), nil)
		})

		It("should install utilizing charts; k8ssandra operator then cluster", func() {

			By("Expecting cass-operator release deployed.")
			InstallChart(ko.GetOperatorCtx(), operatorChart)
			WaitFor(func() bool { return IsReleaseDeployed(ko.GetOperatorCtx()) },
				"Cass-operator release to be deployed", 6, 30)
			Expect(IsReleaseDeployed(ko.GetOperatorCtx())).Should(BeTrue())

			By("Expecting a cass-cluster to be deployed.")
			ctx := ko.GetClusterCtx()
			ctx.SetArgs([]string{"--set", "size=1", "--set", "k8ssandra.namespace=" + ko.GetClusterCtx().Namespace,
				"--set", "clusterName=" + k8ssandraClusterName})
			InstallChartWithValues(ctx, clusterChart, clusterValues)

			WaitFor(func() bool { return IsReleaseDeployed(ko.GetClusterCtx()) },
				"Cass-cluster release to be deployed", 6, 30)
			Expect(IsReleaseDeployed(ko.GetClusterCtx())).Should(BeTrue())
			k8s.WaitUntilAllNodesReady(GinkgoT(), ko.GetClusterCtx().KubeOptions, 15, 3*time.Second)

			By("Expecting a running cass-operator with pod available.")
			var cassOperatorLabel = []PodLabel{{Key: "name", Value: "cass-operator"}}
			WaitFor(func() bool { return IsPodWithLabel(ko.GetOperatorCtx(), cassOperatorLabel) },
				"Cass-operator with pod availability", 10, 60)
			Expect(IsPodWithLabel(ko.GetOperatorCtx(), cassOperatorLabel)).Should(BeTrue())

			By("Expecting a Cass-cluster datacenter with name dc1.")
			Expect(IsExisting(ko.GetClusterCtx(), "CassandraDatacenter",
				"cassandradatacenter.cassandra.datastax.com/dc1")).Should(BeTrue())
			WaitFor(func() bool { return IsDeploymentReady(ko.GetClusterCtx(), "cass-operator") },
				"Deployment of cass-operator", 10, 60)

			By("Expecting datacenter with node state started.")
			var labels = []PodLabel{{Key: "cassandra.datastax.com/datacenter", Value: "dc1"},
				{Key: "cassandra.datastax.com/node-state", Value: "Started"}}
			WaitFor(func() bool { return IsPodWithLabel(ko.GetClusterCtx(), labels) },
				"Cassdc dc1 Running", 20, 120)

		})
	})
})
