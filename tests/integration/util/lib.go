package util

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ReleaseDetail model of release state details
type ReleaseDetail struct {
	Name       string
	Namespace  string
	Revision   string
	UpdatedZDT string
	Status     string
	Chart      string
	AppVersion string
}

// Options model for cluster and operator
type Options struct {
	Cluster  *helm.Options
	Operator *helm.Options
}

// InstallChart chart with specific release.
func InstallChart(t *testing.T, helmOptions *helm.Options, chartName string, releaseName string) {

	result, err := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "install", "-n",
		helmOptions.KubectlOptions.Namespace, releaseName, chartName)

	Log(t, "Helm", "Install", result, err)
	require.Nil(t, err)
	k8s.WaitUntilAllNodesReadyE(t, helmOptions.KubectlOptions, 15, 3*time.Second)
}

// UninstallRelease release target to be uninstalled.
func UninstallRelease(t *testing.T, helmOptions *helm.Options, releaseName string) {

	result, err := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "uninstall", "-n",
		helmOptions.KubectlOptions.Namespace, releaseName, "--no-hooks")
	Log(t, "Helm", "Uninstall", result, err)
	Pause(t, "Pause for Uninstall completion", 1, 6*time.Second)
}

// LookupReleases constructs results of current release details.  nil returned upon error condition.
func LookupReleases(t *testing.T, helmOptions *helm.Options) []ReleaseDetail {

	var releaseDetail []ReleaseDetail
	result, err := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "ls", "-n", helmOptions.KubectlOptions.Namespace, "-o", "json")
	if err != nil {
		return nil
	}

	json.Unmarshal([]byte(result), &releaseDetail)
	return releaseDetail
}

// Log helper for output of action details.
func Log(t *testing.T, subject string, action string, result string, errorDetail error) {

	if errorDetail != nil {
		logger.Log(t, fmt.Sprintf("[ %s ] [ %s] [ Error: %s]", subject, action, errorDetail.Error()))
	} else {
		if result != "" {
			logger.Log(t, fmt.Sprintf("[ %s ] [ %s ] [ %s]", subject, action, result))
		} else {
			logger.Log(t, fmt.Sprintf("[ %s ] [ %s ]", subject, action))
		}
	}
}

// LogResults helper for output of map-based result details
func LogResults(t *testing.T, message string, results map[string]string) {
	for k, v := range results {
		logger.Log(t, fmt.Sprintf("%s Id: %s Result: %s", message, k, v))
	}
}

// Pause helper responsible for pausing for a specified time duration and interval.
func Pause(t *testing.T, message string, interval, timeout time.Duration) {

	logger.Log(t, fmt.Sprintf("%s Timeout: %s", message, timeout))
	wait.Poll(interval, timeout, func() (bool, error) {
		if true {
			time.Sleep(1 * time.Second)
		}
		return false, nil
	})
}
