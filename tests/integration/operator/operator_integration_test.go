package operator

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/k8ssandra/k8ssandra/tests/integration/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	suiteName           = "k8ssandra-integration-suite"
	operatorReleaseName = "k8ssandra"
	namespace           = "k8ssandra"
	clusterReleaseName  = "cassdc"
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
	return ctx
}

var _ = Describe(suiteName, func() {

	var options K8ssandraOptionsContext

	Describe("baseline required preconditions", func() {

		It("should perform setup followed by operator install", func() {

			options = setup()

			By("Expecting operator and cluster options to be setup")
			Ω(options).ShouldNot(BeNil())
			Ω(options.Cluster).ShouldNot(BeNil())
			Ω(options.Operator).ShouldNot(BeNil())
		})
	})

	Describe("k8ssandra cluster integration", func() {

		BeforeEach(func() {
			options = setup()
		})

		It("should install utilizing charts; k8ssandra operator then cluster", func() {

			InstallChart(options.Operator, operatorChart, operatorReleaseName)
			InstallChart(options.Cluster, clusterChart, clusterReleaseName)

			By("Expecting to have the operator and cluster in deployed status")
			Ω(IsReleaseDeployed(options.Operator, "k8ssandra")).Should(BeTrue())
			Ω(IsReleaseDeployed(options.Cluster, "cassdc")).Should(BeTrue())

			By("Expecting to have running cass-operator pod")
			Ω(IsPodRunning(options.Operator, "cass-operator")).Should(BeTrue())

			By("Expecting to have labeled cass-operator pod existing")
			Ω(IsLabeledPodExisting(options.Operator, "cass-operator")).Should(BeTrue())

			Pause("Waiting before fully cleaning up and verifying", 1, 5)

			By("First the cluster uninstall")
			UninstallRelease(options.Cluster, clusterReleaseName)

			By("Followed by the operator uninstall")
			UninstallRelease(options.Operator, operatorReleaseName)

			Pause("Waiting to verify the uninstall", 1, 10)
			Ω(IsLabeledPodExisting(options.Operator, "cass-operator")).Should(BeFalse())
			Ω(IsReleaseDeployed(options.Operator, "k8ssandra")).Should(BeFalse())
			Ω(IsReleaseDeployed(options.Cluster, "cassdc")).Should(BeFalse())
		})
	})
})
