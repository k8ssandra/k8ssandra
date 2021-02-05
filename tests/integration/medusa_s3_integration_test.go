package integration

import (
	. "github.com/k8ssandra/k8ssandra/tests/integration/steps"

	. "github.com/onsi/ginkgo"
)

var (
	medusaS3TestTable    = "medusa_test"
	medusaS3TestKeyspace = "medusa"
)

var _ = Describe("Medusa on S3", func() {
	Context("Deploy K8ssandra with Medusa", func() {
		It("Create a kind cluster and a namespace", func() {
			AKindClusterIsRunningAndReachableStep("one worker")
			IInstallTraefikStep()
			ICreateTheNamespaceStep()
			ICanSeeTheNamespaceInTheListOfNamespacesStep()
		})
		It("Create the Medusa secret ", func() {
			ICreateTheMedusaSecretInTheNamespaceApplyingTheFileStep("~/medusa_secret.yaml")
			ICanSeeTheSecretInTheListOfSecretsInTheNamespaceStep("medusa-bucket-key")
		})
		It("Install K8ssandra with Medusa on S3", func() {
			IDeployAClusterWithOptionsInTheNamespaceUsingTheValuesStep("no Traefik", "one_node_cluster_with_medusa_s3.yaml")
		})
		It("Check the presence of expected resources", func() {
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-dc1-all-pods-service")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-dc1-service")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-seed-service")
		})
		It("Create a keyspace and a table", func() {
			ICreateTheTableInTheKeyspaceStep(medusaS3TestTable, medusaS3TestKeyspace)
		})
		It("Load 10 rows and check that we can read that exact number of rows", func() {
			ILoadRowsInTheTableInTheKeyspaceStep(10, medusaS3TestTable, medusaS3TestKeyspace)
			ICanReadRowsInTheTableInTheKeyspaceStep(10, medusaS3TestTable, medusaS3TestKeyspace)
		})
		It("Perform a backup using Medusa", func() {
			IPerformABackupWithMedusaNamedStep("backup1")
		})
		It("Load 10 additional rows and check that we can read 20 rows now", func() {
			ILoadRowsInTheTableInTheKeyspaceStep(10, medusaS3TestTable, medusaS3TestKeyspace)
			ICanReadRowsInTheTableInTheKeyspaceStep(20, medusaS3TestTable, medusaS3TestKeyspace)
		})
		It("Restore the backup", func() {
			IRestoreTheBackupNamedUsingMedusaStep("backup1")
		})
		It("Check that we can read 10 rows after the restore", func() {
			ICanReadRowsInTheTableInTheKeyspaceStep(10, medusaS3TestTable, medusaS3TestKeyspace)
		})
		It("Delete the namespace and the kind cluster", func() {
			IDeleteTheNamespaceStep()
			ICannotSeeTheNamespaceInTheListOfNamespacesStep()
			ICanDeleteTheKindClusterStep()
		})
	})
})
