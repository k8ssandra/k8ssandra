package integration

import (
	. "github.com/k8ssandra/k8ssandra/tests/integration/steps"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Reaper", func() {
	Context("Deploy K8ssandra with Reaper", func() {
		It("Create a kind cluster and a namespace", func() {
			AKindClusterIsRunningAndReachableStep("three workers")
			IInstallTraefikStep()
			ICreateTheNamespaceStep()
			ICanSeeTheNamespaceInTheListOfNamespacesStep()
		})
		It("Install K8ssandra with Reaper", func() {
			IDeployAClusterWithOptionsInTheNamespaceUsingTheValuesStep("default", "three_nodes_cluster_with_reaper.yaml")
		})
		It("Check the presence of expected resources", func() {
			ICanCheckThatResourceOfTypeWithLabelIsPresentInNamespaceStep("service", "app.kubernetes.io/managed-by=reaper-operator")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-dc1-all-pods-service")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-dc1-service")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-seed-service")
		})
		It("Wait for Reaper to be ready", func() {
			IWaitForTheReaperPodToBeReadyInNamespaceStep()
		})
		It("Check that Reaper has registered the required elements", func() {
			ICanSeeThatTheKeyspaceExistsInCassandraInNamespaceStep("reaper_db")
			ICanCheckThatAClusterNamedWasRegisteredInReaperInNamespaceStep("k8ssandra")
		})
		It("Start a repair on the reaper_db keyspace", func() {
			ITriggerARepairOnTheKeyspaceStep("reaper_db")
		})
		It("Wait for at least one segment to be processed and cancel the repair", func() {
			IWaitForAtLeastOneSegmentToBeProcessedStep()
			ICanCancelTheRunningRepairStep()
		})
		It("Delete the namespace and the kind cluster", func() {
			IDeleteTheNamespaceStep()
			ICannotSeeTheNamespaceInTheListOfNamespacesStep()
			ICanDeleteTheKindClusterStep()
		})
	})
})
