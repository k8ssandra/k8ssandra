package steps

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	cassdcapi "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/onsi/gomega"
)

// Medusa related functions
func CreateMedusaSecretWithFile(t *testing.T, namespace, secretFile string) {
	home, _ := os.UserHomeDir()
	medusaSecretPath, _ := filepath.Abs(strings.Replace(secretFile, "~", home, 1))
	k8s.KubectlApply(t, getKubectlOptions(namespace), medusaSecretPath)
}

func PerformBackup(t *testing.T, namespace, backupName string, useLocalCharts bool) {
	backupChartPath, err := filepath.Abs("../../charts/backup")
	g(t).Expect(err).To(BeNil())

	if !useLocalCharts {
		backupChartPath = "k8ssandra/backup"
	}

	helmOptions := &helm.Options{
		SetValues: map[string]string{
			"name":                     backupName,
			"cassandraDatacenter.name": datacenterName,
		},
		KubectlOptions: getKubectlOptions(namespace),
	}
	// Start backup
	helm.Install(t, helmOptions, backupChartPath, backupName)

	// Wait for the backup to be finished
	g(t).Eventually(func() bool {
		return backupIsFinished(t, namespace, backupName)
	}, 2*time.Minute, 5*time.Second).Should(BeTrue())
}

func backupIsFinished(t *testing.T, namespace, backupName string) bool {
	log.Printf("Checking if backup %s is finished...", backupName)
	output, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions(namespace), "get", "cassandrabackup", backupName, "-o", "jsonpath={.status.finishTime}")
	if err == nil && len(output) > 0 {
		return true
	}
	return false
}

func RestoreBackup(t *testing.T, namespace, backupName string) {
	restoreChartPath, err := filepath.Abs("../../charts/restore")
	g(t).Expect(err).To(BeNil(), "Couldn't find the absolute path for restore charts")

	startTime := time.Now()
	dcUpdates := make(map[time.Time]bool)

	helmOptions := &helm.Options{
		SetValues: map[string]string{
			"backup.name":              backupName,
			"cassandraDatacenter.name": datacenterName,
			"name":                     "restore-test2",
		},
		KubectlOptions: getKubectlOptions(namespace),
	}
	helm.Install(t, helmOptions, restoreChartPath, "restore-test")

	// Wait for restore to be completed and Cassandra to be available
	g(t).Eventually(func() bool {
		checkDatacenterUpdates(t, namespace, startTime, dcUpdates)

		return restoreIsFinished(t, namespace, backupName)
	}, retryTimeout, retryInterval).Should(BeTrue())

	g(t).Expect(len(dcUpdates)).To(Equal(1), "There should have only been 1 datacenter update during the restore operation")
}

func checkDatacenterUpdates(t *testing.T, namespace string, start time.Time, updates map[time.Time]bool) {
	dc := &cassdcapi.CassandraDatacenter{}
	if err := testClient.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: datacenterName}, dc); err == nil {
		t.Logf("Failed to get CassandraDatacenter while waiting for restore to finish: %s", err)
		return
	}

	if updating, found := dc.GetCondition(cassdcapi.DatacenterUpdating); found {
		if updating.LastTransitionTime.After(start) {
			updates[updating.LastTransitionTime.Time] = true
		}
	}
}

func restoreIsFinished(t *testing.T, namespace, backupName string) bool {
	log.Printf("Checking if restore %s is finished...", backupName)
	output, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions(namespace), "get", "cassandrarestore", "restore-test2", "-o", "jsonpath={.status.finishTime}")
	if err == nil && len(output) > 0 {
		return true
	}
	return false
}
