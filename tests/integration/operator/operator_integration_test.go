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

const (
	SuiteName            = "k8ssandra-operator-integration-suite"
	KindClusterName      = "k8ssandra-cluster"
	K8ssandraContext     = "kind-k8ssandra-cluster"
	K8ssandraReleaseName = "k8ssandra"
	NetworkReleaseName   = "traefik"
	Namespace            = "k8ssandra"

	K8ssandraSingeNodeValues = "../preconditions/k8ssandra-single-node.yaml"
	TraefikValues            = "../preconditions/k8ssandra-traefik-values.yaml"
	KubeConfigFile           = "../preconditions/k8ssandra-config-test.yaml"
	KindConfigFile           = "../preconditions/k8ssandra-kind-config.yaml"
	KindImage                = "kindest/node:v1.20.2"

	includeDatacenterVerification = false
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, SuiteName)
}

// setup used for creating test option context and cleanup
func setup() K8ssandraOptions {

	options := CreateK8ssandraOptions(k8s.NewKubectlOptions(K8ssandraContext,
		"", Namespace), K8ssandraReleaseName, NetworkReleaseName)
	ConfigureContext(options.GetOperatorCtx(), KubeConfigFile)

	// Check if kind is available.
	if IsKindClusterExisting(KindClusterName) {
		DeleteKindCluster(KindClusterName)
	}

	// Create Cluster fresh
	CreateKindCluster(KindClusterDetail{Name: KindClusterName,
		Image: KindImage, ConfigFile: KindConfigFile})

	RepoUpdate(options.GetOperatorCtx())
	return options
}

var _ = Describe(SuiteName, func() {

	Describe("k8ssandra Operator and Cassandra integration", func() {

		var ko K8ssandraOptions
		BeforeSuite(func() {
			ko = setup()
			Expect(ko).ToNot(BeNil())

			By("Expecting Traefik to be deployed")
			InstallChartWithValues(ko.GetNetworkCtx(), "traefik/traefik", TraefikValues)
			Expect(IsReleaseDeployed(ko.GetNetworkCtx())).Should(BeTrue())
			k8s.WaitUntilAllNodesReady(GinkgoT(), ko.GetNetworkCtx().KubeOptions, 15, 3*time.Second)
		})

		It("should install utilizing charts; k8ssandra operator then cluster", func() {

			By("Expecting k8ssandra to be deployed.")
			ctx := ko.GetOperatorCtx()
			ctx.SetArgs([]string{
				// "--set", "datacenters[0].size=2",
				"--set", "k8ssandra.Namespace=" + ko.GetOperatorCtx().Namespace,
				"--set", "clusterName=" + KindClusterName})

			InstallChartWithValues(ctx, "k8ssandra/k8ssandra", K8ssandraSingeNodeValues)

			WaitFor(func() bool { return IsReleaseDeployed(ko.GetOperatorCtx()) },
				"k8ssandra to be deployed", 6, 30)
			Expect(IsReleaseDeployed(ko.GetOperatorCtx())).Should(BeTrue())
			k8s.WaitUntilAllNodesReady(GinkgoT(), ko.GetOperatorCtx().KubeOptions, 15, 3*time.Second)

			By("Expecting a running k8ssandra cass-operator with pod available.")
			var cassOperatorLabel = []PodLabel{{Key: "app.kubernetes.io/name", Value: "cass-operator"}}
			WaitFor(func() bool { return IsPodWithLabel(ko.GetOperatorCtx(), cassOperatorLabel) },
				"Cass-operator with pod availability", 10, 300)
			Expect(IsPodWithLabel(ko.GetOperatorCtx(), cassOperatorLabel)).Should(BeTrue())

			// Requires considerable time, optional verification as a post test can do this after some time.
			if includeDatacenterVerification {

				By("Expecting a datacenter with name dc1.")
				Expect(IsExisting(ko.GetOperatorCtx(), "CassandraDatacenter",
					"cassandradatacenter.cassandra.datastax.com/dc1")).Should(BeTrue())
				WaitFor(func() bool { return IsDeploymentReady(ko.GetOperatorCtx(), "cass-operator") },
					"Deployment of cass-operator", 10, 60)

				By("Expecting datacenter with node state started.")
				var labels = []PodLabel{{Key: "cassandra.datastax.com/datacenter", Value: "dc1"},
					{Key: "cassandra.datastax.com/node-state", Value: "Started"}}
				WaitFor(func() bool { return IsPodWithLabel(ko.GetOperatorCtx(), labels) },
					"Cassdc dc1 Running", 15, 600)
			}

		})
	})
})
