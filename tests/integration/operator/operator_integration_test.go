package operator

import (
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
	kindClusterName      = "k8ssandra-cluster"
	k8ssandraContext     = "kind-k8ssandra-cluster"
	k8ssandraReleaseName = "k8ssandra"
	networkReleaseName   = "traefik"
	namespace            = "k8ssandra"

	operatorChart  = "../../../charts/k8ssandra"
	operatorValues = "../preconditions/k8ssandra-values.yaml"
	traefikValues  = "../preconditions/k8ssandra-traefik-values.yaml"
	kubeConfigFile = "../preconditions/k8ssandra-config-test.yaml"
	kindConfigFile = "../preconditions/k8ssandra-kind-config.yaml"
	kindImage      = "kindest/node:v1.18.2"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, suiteName)
}

// setup used for creating test option context and cleanup
func setup() K8ssandraOptions {

	options := CreateK8ssandraOptions(k8s.NewKubectlOptions(k8ssandraContext,
		"", namespace), k8ssandraReleaseName, networkReleaseName)
	ConfigureContext(options.GetOperatorCtx(), kubeConfigFile)

	// Check if kind is available.
	if IsKindClusterExisting(kindClusterName) {
		DeleteKindCluster(kindClusterName)
	}

	CreateKindCluster(KindClusterDetail{Name: kindClusterName,
		Image: kindImage, ConfigFile: kindConfigFile})
	RepoUpdate(options.GetOperatorCtx())
	ConfigureContext(options.GetOperatorCtx(), kubeConfigFile)
	return options
}

var _ = Describe(suiteName, func() {

	Describe("k8ssandra Operator and Cassandra integration", func() {

		var ko K8ssandraOptions
		BeforeEach(func() {
			ko = setup()
			Expect(ko).ToNot(BeNil())

			By("Expecting Traefik to be deployed")
			InstallChartWithValues(ko.GetNetworkCtx(), "traefik/traefik", traefikValues)
			Expect(IsReleaseDeployed(ko.GetNetworkCtx())).Should(BeTrue())
			k8s.WaitUntilAllNodesReady(GinkgoT(), ko.GetNetworkCtx().KubeOptions, 15, 3*time.Second)
		})

		It("should install utilizing charts; k8ssandra operator then cluster", func() {

			By("Expecting k8ssandra to be deployed.")
			ctx := ko.GetOperatorCtx()
			ctx.SetArgs([]string{"--set", "size=1", "--set", "k8ssandra.namespace=" + ko.GetOperatorCtx().Namespace,
				"--set", "clusterName=" + kindClusterName})
			InstallChartWithValues(ctx, operatorChart, operatorValues)

			WaitFor(func() bool { return IsReleaseDeployed(ko.GetOperatorCtx()) },
				"k8ssandra to be deployed", 6, 30)
			Expect(IsReleaseDeployed(ko.GetOperatorCtx())).Should(BeTrue())
			k8s.WaitUntilAllNodesReady(GinkgoT(), ko.GetOperatorCtx().KubeOptions, 15, 3*time.Second)

			By("Expecting a running k8ssandra cass-operator with pod available.")
			var cassOperatorLabel = []PodLabel{{Key: "name", Value: "cass-operator"}}
			WaitFor(func() bool { return IsPodWithLabel(ko.GetOperatorCtx(), cassOperatorLabel) },
				"Cass-operator with pod availability", 10, 60)
			Expect(IsPodWithLabel(ko.GetOperatorCtx(), cassOperatorLabel)).Should(BeTrue())

			By("Expecting a datacenter with name dc1.")
			Expect(IsExisting(ko.GetOperatorCtx(), "CassandraDatacenter",
				"cassandradatacenter.cassandra.datastax.com/dc1")).Should(BeTrue())
			WaitFor(func() bool { return IsDeploymentReady(ko.GetOperatorCtx(), "cass-operator") },
				"Deployment of cass-operator", 10, 60)

			By("Expecting datacenter with node state started.")
			var labels = []PodLabel{{Key: "cassandra.datastax.com/datacenter", Value: "dc1"},
				{Key: "cassandra.datastax.com/node-state", Value: "Started"}}
			WaitFor(func() bool { return IsPodWithLabel(ko.GetOperatorCtx(), labels) },
				"Cassdc dc1 Running", 20, 260)

		})
	})
})
