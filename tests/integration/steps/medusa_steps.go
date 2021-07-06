package steps

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/gomega"
)

// Medusa-related functions

func CreateMedusaSecretWithFile(t *testing.T, namespace, secretFile string) {
	home, _ := os.UserHomeDir()
	medusaSecretPath, _ := filepath.Abs(strings.Replace(secretFile, "~", home, 1))
	if _, err := os.Stat(medusaSecretPath); err != nil {
		t.Fatalf("Medusa secret file %s does not exist or is not readable: %v", medusaSecretPath, err)
	} else {
		k8s.KubectlApply(t, getKubectlOptions(namespace), medusaSecretPath)
	}
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
	}, retryTimeout, retryInterval).Should(BeTrue())
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
		return restoreIsFinished(t, namespace, backupName)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func restoreIsFinished(t *testing.T, namespace, backupName string) bool {
	log.Printf("Checking if restore %s is finished...", backupName)
	output, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions(namespace), "get", "cassandrarestore", "restore-test2", "-o", "jsonpath={.status.finishTime}")
	if err == nil && len(output) > 0 {
		return true
	}
	return false
}
