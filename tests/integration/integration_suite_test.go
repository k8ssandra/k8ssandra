package integration

import (
	"log"
	"strings"

	. "github.com/k8ssandra/k8ssandra/tests/integration/steps"

	"fmt"
	"os"
	"testing"
)

const (
	medusaTestTable    = "medusa_test"
	medusaTestKeyspace = "medusa"
	traefikNamespace   = "traefik"
	minioNamespace     = "minio"
)

// "" means latest stable version in the Helm repo
var upgradeStartVersions = []string{"v1.0.0", "latest"}

func TestMain(m *testing.M) {
	err := InitTestClient()
	if err != nil {
		log.Fatalf("failed to initialize test client: %s", err)
	}

	os.Exit(m.Run())
}

func initializeCluster(t *testing.T) string {
	log.Println(Step("Initializing cluster"))

	CheckK8sClusterIsReachable(t)
	InstallTraefik(t)
	namespace := CreateNamespace(t)
	CheckNamespaceWasCreated(t, namespace)
	return namespace
}

func cleanupClusterOption() string {

	if os.Getenv("CLUSTER_CLEANUP") != "" {
		return os.Getenv("CLUSTER_CLEANUP")
	} else {
		return "always"
	}
}

func shouldCleanupCluster(success bool) bool {
	if cleanupClusterOption() == "always" || (cleanupClusterOption() == "success" && success) {
		return true
	}
	return false
}

func cleanupCluster(t *testing.T, namespace string, success bool) {
	if shouldCleanupCluster(success) {
		log.Println(Step("Cleaning up cluster"))
		UninstallHelmRealeaseAndNamespace(t, "k8ssandra", namespace)
		WaitForCassandraDatacenterDeletion(t, namespace)
		UninstallHelmRealeaseAndNamespace(t, "traefik", traefikNamespace)
		UninstallHelmRealeaseAndNamespace(t, "minio", minioNamespace)
	} else {
		log.Println(Info("Not cleaning up cluster as requested"))
	}
}

// Full stack scenario:
// - Install Minio
// - Create the Minio credentials secret
// - Register a cluster with 3 nodes
// - Run the Reaper test scenario
// - Run the Medusa test scenario
// - Run the Prometheus test scenario
// - Run the Grafana test scenario
// - Run the Stargate test scenario
// - Terminate the namespace and cleanup the cluster
func TestFullStackScenario(t *testing.T) {
	const (
		medusaBackend = "Minio"
		backupName    = "backup1"
	)

	namespace := initializeCluster(t)

	success := t.Run("Full Stack Test", func(t *testing.T) {
		createMedusaSecretAndInstallDeps(t, namespace, medusaBackend)
		deployFullStackCluster(t, namespace, true)

		t.Run("Test Reaper", func(t *testing.T) {
			testReaper(t, namespace)
		})

		t.Run("Test Medusa", func(t *testing.T) {
			testMedusa(t, namespace, medusaBackend, backupName, true)
		})

		t.Run("Test Prometheus", func(t *testing.T) {
			testPrometheus(t, namespace)
		})

		t.Run("Test Grafana", func(t *testing.T) {
			testGrafana(t, namespace)
		})

		t.Run("Test Stargate", func(t *testing.T) {
			// The backup/restore test runs before this. Because it shuts down
			// the Cassandra cluster, we need to restart Stargate. See
			// https://github.com/k8ssandra/k8ssandra/issues/411 for details.
			releaseName := "k8ssandra"
			dcName := "dc1"
			RestartStargate(t, releaseName, dcName, namespace)
			testStargate(t, namespace)
		})
	})

	cleanupCluster(t, namespace, success)
}

func deployFullStackCluster(t *testing.T, namespace string, useLocalCharts bool) {
	DeployClusterWithValues(t, namespace, "minio", "cluster_full_stack.yaml", 3, false, useLocalCharts, "")
	checkResourcePresenceForReaper(t, namespace)
	waitForReaperPod(t, namespace)
	checkReaperRegistered(t, namespace)
}

// Reaper scenario:
// - Install Traefik
// - Create a namespace
// - Register a cluster with 3 Cassandra nodes
// - Verify that Reaper correctly initializes
// - Start a repair on the reaper_db keyspace
// - Wait for at least one segment to be processed
// - Cancel the repair
// - Terminate the namespace and delete the cluster
func TestReaperDeploymentScenario(t *testing.T) {
	namespace := initializeCluster(t)
	success := t.Run("Test Reaper", func(t *testing.T) {
		deployClusterForReaper(t, namespace, true)
		testReaper(t, namespace)
	})
	cleanupCluster(t, namespace, success)
}

func testReaper(t *testing.T, namespace string) {
	log.Println(Step("Testing Reaper..."))
	repairId := triggerRepair(t, namespace)
	waitForSegmentDoneAndCancel(t, repairId)
}

func deployClusterForReaper(t *testing.T, namespace string, useLocalCharts bool) {
	log.Println(Info("Deploying K8ssandra and waiting for Reaper to be ready"))
	DeployClusterWithValues(t, namespace, "default", "cluster_with_reaper.yaml", 3, false, useLocalCharts, "")
	checkResourcePresenceForReaper(t, namespace)
	waitForReaperPod(t, namespace)
	checkReaperRegistered(t, namespace)
}

func checkResourcePresenceForReaper(t *testing.T, namespace string) {
	CheckResourceWithLabelsIsPresent(t, namespace, "service", map[string]string{"app.kubernetes.io/managed-by": "reaper-operator"})
	CheckClusterExpectedResources(t, namespace)
}

func waitForReaperPod(t *testing.T, namespace string) {
	WaitForReaperPod(t, namespace)
}

func checkReaperRegistered(t *testing.T, namespace string) {
	CheckKeyspaceExists(t, namespace, "reaper_db")
	CheckClusterIsRegisteredInReaper(t, "k8ssandra")
}

func triggerRepair(t *testing.T, namespace string) string {
	log.Println(Info("Starting a repair"))
	return TriggerRepair(t, namespace, "reaper_db")
}

func waitForSegmentDoneAndCancel(t *testing.T, repairId string) {
	log.Println(Info("Waiting for one segment to be repaired and canceling run"))
	WaitForOneSegmentToBeDone(t, repairId)
	CancelRepair(t, repairId)
}

// Medusa scenario (invoked with a specific backend name):
// - Register a cluster with 1 node
// - Potentially install backend specific dependencies (such as Minio)
// - Create the backend credentials secret
// - Create a keyspace and a table
// - Load 10 rows and check that we can read 10 rows
// - Perform a backup with Medusa
// - Load 10 rows and check that we can read 20 rows
// - Restore the backup
// - Verify that we can read 10 rows
// - Terminate the namespace and delete the cluster
func TestMedusaDeploymentScenario(t *testing.T) {
	const backupName = "backup1"
	backends := []string{"Minio", "S3", "local"}
	for _, backend := range backends {
		t.Run(fmt.Sprintf("Medusa on %s", backend), func(t *testing.T) {
			namespace := initializeCluster(t)
			medusaSuccess := t.Run("Test backup and restore", func(t *testing.T) {
				createMedusaSecretAndInstallDeps(t, namespace, backend)
				deployClusterForMedusa(t, namespace, backend, 1, true, "")
				testMedusa(t, namespace, backend, backupName, true)
				scaleUpCassandra(t, namespace, backend, 2)
			})
			cleanupCluster(t, namespace, medusaSuccess)
		})
	}
}

func testMedusa(t *testing.T, namespace, backend, backupName string, useLocalChartForBackup bool) {
	log.Println(Step("Testing Medusa..."))
	log.Println("Creating test keyspace and table")
	CreateCassandraTable(t, namespace, medusaTestTable, medusaTestKeyspace)

	loadRowsAndCheckCount(t, namespace, 10, 10)

	log.Println(Info("Backing up Cassandra"))
	PerformBackup(t, namespace, backupName, useLocalChartForBackup)

	if !useLocalChartForBackup {
		// This will upgrade the cluster to the local version if the stable chart was used to perform the backup
		scaleUpCassandra(t, namespace, backend, 1)
	}

	loadRowsAndCheckCount(t, namespace, 10, 20)

	log.Println(Info("Restoring backup and checking row count"))
	RestoreBackup(t, namespace, backupName)
	CheckRowCountInTable(t, 10, namespace, medusaTestTable, medusaTestKeyspace)
}

func deployClusterForMedusa(t *testing.T, namespace, backend string, nodes int, useLocalCharts bool, version string) {
	log.Println(Info(fmt.Sprintf("Deploying K8ssandra with Medusa using %s", backend)))
	valuesFile := fmt.Sprintf("cluster_with_medusa_%s.yaml", strings.ToLower(backend))
	DeployClusterWithValues(t, namespace, strings.ToLower(backend), valuesFile, nodes, false, useLocalCharts, version)
	CheckClusterExpectedResources(t, namespace)
}

func loadRowsAndCheckCount(t *testing.T, namespace string, rowsToLoad, rowsExpected int) {
	log.Println(Info(fmt.Sprintf("Loading %d rows and checking we have %d after", rowsToLoad, rowsExpected)))
	LoadRowsInTable(t, rowsToLoad, namespace, medusaTestTable, medusaTestKeyspace)
	CheckRowCountInTable(t, rowsExpected, namespace, medusaTestTable, medusaTestKeyspace)
}

func createMedusaSecretAndInstallDeps(t *testing.T, namespace, backend string) {
	log.Println(Info("Creating medusa secret to access the backend"))
	if backend == "local" {
		return
	} else if backend == "Minio" {
		DeployMinioAndCreateBucket(t, "k8ssandra-medusa")
		CreateMedusaSecretWithFile(t, namespace, "secret/medusa_minio_secret.yaml")
	} else {
		// S3
		CreateMedusaSecretWithFile(t, namespace, "~/medusa_secret.yaml")
	}
}

func scaleUpCassandra(t *testing.T, namespace, backend string, nodes int) {
	log.Println(Info("Scaling up Cassandra"))
	valuesFile := fmt.Sprintf("cluster_with_medusa_%s.yaml", strings.ToLower(backend))
	DeployClusterWithValues(t, namespace, strings.ToLower(backend), valuesFile, nodes, true, true, "")
}

// Monitoring scenario:
// - Install Traefik
// - Create a namespace
// - Register a cluster with three Cassandra nodes and one Stargate node
// - Check that Prometheus is reachable through its REST API
// - Check the number of active Prometheus targets
// - Check that Grafana is reachable through http
// - Terminate the namespace and delete the cluster
func TestMonitoringDeploymentScenario(t *testing.T) {
	namespace := initializeCluster(t)

	success := t.Run("Test Monitoring", func(t *testing.T) {
		deployClusterForMonitoring(t, namespace)

		t.Run("Test Prometheus", func(t *testing.T) {
			testPrometheus(t, namespace)
		})

		t.Run("Test Grafana", func(t *testing.T) {
			testGrafana(t, namespace)
		})
	})

	cleanupCluster(t, namespace, success)
}

func deployClusterForMonitoring(t *testing.T, namespace string) {
	DeployClusterWithValues(t, namespace, "default", "cluster_with_stargate_and_monitoring.yaml", 3, false, true, "")
	CheckClusterExpectedResources(t, namespace)
	WaitForStargatePodReady(t, namespace)
}

// Prometheus tests
func testPrometheus(t *testing.T, namespace string) {
	log.Println(Step("Testing Prometheus..."))
	PodWithLabelsIsReady(t, namespace, map[string]string{"app": "prometheus"})
	CheckPrometheusMetricExtraction(t)
	expectedActiveTargets := CountMonitoredItems(t, namespace)
	CheckPrometheusActiveTargets(t, expectedActiveTargets) // We're monitoring 3 Cassandra nodes and 1 Stargate instance
}

// Grafana tests
func testGrafana(t *testing.T, namespace string) {
	log.Println(Step("Testing Grafana..."))
	PodWithLabelsIsReady(t, namespace, map[string]string{"app.kubernetes.io/name": "grafana"})
	CheckGrafanaIsReachable(t)
}

// Stargate scenario:
// - Install Traefik
// - Create a namespace
// - Register a cluster with three Cassandra nodes and one Stargate node
// - Check Stargate rollout
// - Create a document and read it back through the Stargate document API
// - Terminate the namespace and delete the cluster
func TestStargateDeploymentScenario(t *testing.T) {
	namespace := initializeCluster(t)

	success := t.Run("Test Stargate", func(t *testing.T) {
		deployClusterForStargate(t, namespace)
		testStargate(t, namespace)
	})
	cleanupCluster(t, namespace, success)
}

func deployClusterForStargate(t *testing.T, namespace string) {
	DeployClusterWithValues(t, namespace, "default", "cluster_with_stargate.yaml", 3, false, true, "")
	CheckClusterExpectedResources(t, namespace)
	WaitForStargatePodReady(t, namespace)
}

func testStargate(t *testing.T, namespace string) {
	WaitForAuthEndpoint(t) // Wait for the auth endpoint to be reachable, this takes a little time after the Stargate rollout is complete
	log.Println(Step("Writing data to the Stargate document API"))
	token := GenerateStargateAuthToken(t, namespace)
	docNamespace := CreateStargateDocumentNamespace(t, token)
	log.Println(fmt.Sprintf("Created Stargate document namespace: %s", docNamespace))
	documentId := WriteStargateDocument(t, token, docNamespace)
	log.Println(fmt.Sprintf("Created document with id: %s", documentId))
	CheckStargateDocumentExists(t, token, docNamespace, documentId)
}

func TestUpgradeScenario(t *testing.T) {
	for _, startVersion := range upgradeStartVersions {
		namespace := initializeCluster(t)
		success := t.Run(fmt.Sprintf("Upgrade from %s", startVersion), func(t *testing.T) {
			// Install first production version
			DeployClusterWithValues(t, namespace, "default", "cluster_with_reaper.yaml", 1, false, false, startVersion)
			checkResourcePresenceForReaper(t, namespace)
			waitForReaperPod(t, namespace)

			// Upgrade to current version
			DeployClusterWithValues(t, namespace, "default", "cluster_with_reaper.yaml", 1, true, true, "")
			checkResourcePresenceForReaper(t, namespace)
			waitForReaperPod(t, namespace)
			checkReaperRegistered(t, namespace)
		})

		cleanupCluster(t, namespace, success)
		if !success {
			t.FailNow()
		}
	}
}

// Upgrade scenario:
// - Install Traefik
// - Create a namespace
// - Register a cluster with one Cassandra nodes and one Stargate node using the latest stable
// - Load data and take a backup using the stable version chart
// - Upgrade the cluster to the local version
// - Restore the backup using the local version chart and check that data is here
func TestRestoreAfterUpgrade(t *testing.T) {
	const (
		medusaBackend = "Minio"
		backupName    = "backup1"
	)

	for _, startVersion := range upgradeStartVersions {
		if startVersion != "v1.0.0" {
			// K8ssandra 1.0.0 didn't support Minio as Medusa backend
			namespace := initializeCluster(t)

			success := t.Run(fmt.Sprintf("Medusa upgrade from %s", startVersion), func(t *testing.T) {
				createMedusaSecretAndInstallDeps(t, namespace, medusaBackend)
				deployClusterForMedusa(t, namespace, medusaBackend, 1, false, startVersion)
				testMedusa(t, namespace, medusaBackend, backupName, false)

			})

			cleanupCluster(t, namespace, success)
		}
	}
}
