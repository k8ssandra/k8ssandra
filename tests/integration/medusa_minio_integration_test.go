package integration

import (
	. "github.com/k8ssandra/k8ssandra/tests/integration/steps"

	. "github.com/onsi/ginkgo"
)

var (
	medusaMinioTestTable    = "medusa_test"
	medusaMinioTestKeyspace = "medusa"
)

var _ = Describe("Medusa on MinIO", func() {
	Context("Deploy K8ssandra with Medusa", func() {
		It("Create a kind cluster and a namespace", func() {
			AKindClusterIsRunningAndReachableStep("one worker")
			IInstallTraefikStep()
			ICreateTheNamespaceStep()
			ICanSeeTheNamespaceInTheListOfNamespacesStep()
		})
		It("Create the Medusa secret ", func() {
			IDeployMinIOUsingHelmAndCreateTheBucketStep("k8ssandra-medusa")
			ICreateTheMedusaSecretInTheNamespaceApplyingTheFileStep("secret/medusa_minio_secret.yaml")
			ICanSeeTheSecretInTheListOfSecretsInTheNamespaceStep("medusa-bucket-key")
		})
		It("Install K8ssandra with Medusa on MinIO", func() {
			IDeployAClusterWithOptionsInTheNamespaceUsingTheValuesStep("minio", "one_node_cluster_with_medusa_minio.yaml")
		})
		It("Check the presence of expected resources", func() {
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-dc1-all-pods-service")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-dc1-service")
			ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep("service", "k8ssandra-seed-service")
		})
		It("Create a keyspace and a table", func() {
			ICreateTheTableInTheKeyspaceStep(medusaMinioTestTable, medusaMinioTestKeyspace)
		})
		It("Load 10 rows and check that we can read that exact number of rows", func() {
			ILoadRowsInTheTableInTheKeyspaceStep(10, medusaMinioTestTable, medusaMinioTestKeyspace)
			ICanReadRowsInTheTableInTheKeyspaceStep(10, medusaMinioTestTable, medusaMinioTestKeyspace)
		})
		It("Perform a backup using Medusa", func() {
			IPerformABackupWithMedusaNamedStep("backup1")
		})
		It("Load 10 additional rows and check that we can read 20 rows now", func() {
			ILoadRowsInTheTableInTheKeyspaceStep(10, medusaMinioTestTable, medusaMinioTestKeyspace)
			ICanReadRowsInTheTableInTheKeyspaceStep(20, medusaMinioTestTable, medusaMinioTestKeyspace)
		})
		It("Restore the backup", func() {
			IRestoreTheBackupNamedUsingMedusaStep("backup1")
		})
		It("Check that we can read 10 rows after the restore", func() {
			ICanReadRowsInTheTableInTheKeyspaceStep(10, medusaMinioTestTable, medusaMinioTestKeyspace)
		})
		It("Delete the namespace and the kind cluster", func() {
			IDeleteTheNamespaceStep()
			ICannotSeeTheNamespaceInTheListOfNamespacesStep()
			ICanDeleteTheKindClusterStep()
		})
	})
})
