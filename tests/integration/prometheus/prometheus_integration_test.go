package prometheus

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/k8ssandra/k8ssandra/tests/integration/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	suiteName           = "k8ssandra-prometheus-integration-suite"
	namespace           = "k8ssandra"
	clusterReleaseName  = "cassdc"
	operatorReleaseName = "k8ssandra"
	networkReleaseName  = "traefik"
	operatorChart       = "../../../charts/k8ssandra"
	clusterChart        = "../../../charts/k8ssandra-cluster"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, suiteName)
}

// setup used for creating test option context and cleanup of operator and cluster.
func setup() K8ssandraOptionsContext {

	kubeOptions := k8s.NewKubectlOptions("", "", namespace)

	operatorOptions := &helm.Options{KubectlOptions: kubeOptions}
	clusterOptions := &helm.Options{KubectlOptions: kubeOptions}
	networkOptions := &helm.Options{KubectlOptions: kubeOptions}

	ctx := CreateK8ssandraOptionsContext(operatorOptions, clusterOptions, networkOptions)
	UninstallRelease(ctx.Cluster, clusterReleaseName)
	UninstallRelease(ctx.Operator, operatorReleaseName)
	UninstallRelease(ctx.Network, networkReleaseName)
	return ctx
}

var _ = Describe(suiteName, func() {

	var options K8ssandraOptionsContext

	Describe("k8ssandra route prefix for Prometheus", func() {

		BeforeEach(func() {
			options = setup()
		})

		It("should install traefik, operator, & cluster where cluster has overrides", func() {

			By("Expecting traefik to be deployed")
			InstallChart(options.Network, "traefik/traefik", networkReleaseName)
			Ω(IsReleaseDeployed(options.Network, networkReleaseName)).Should(BeTrue())

			By("Expecting cass-operator to be deployed")
			InstallChart(options.Operator, operatorChart, operatorReleaseName)
			Ω(IsReleaseDeployed(options.Operator, operatorReleaseName)).Should(BeTrue())

			By("Expecting cass-cluster to be deployed, to include custom prom-values")
			options.Cluster.Args = []string{"--set-string", "routePrefix=/prometheus"}
			InstallChart(options.Cluster, clusterChart, clusterReleaseName)
			Ω(IsReleaseDeployed(options.Cluster, clusterReleaseName)).Should(BeTrue())

			By("Expecting running cass-operator pod")
			Ω(IsPodRunning(options.Operator, "cass-operator")).Should(BeTrue())

			By("Expecting labeled cass-operator pod to exist")
			Ω(IsLabeledPodExisting(options.Operator, "cass-operator")).Should(BeTrue())

		})
	})
})
