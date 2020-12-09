package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

// PodDetail model
type PodDetail struct {
	Name   string
	Ready  string
	Status string
}

// K8ssandraOptionsContext options model for cluster and operator
type K8ssandraOptionsContext struct {
	Cluster  OptionsContext
	Operator OptionsContext
}

// OptionsContext model for testing options context.
type OptionsContext struct {
	helmOptions *helm.Options
	kubeOptions *k8s.KubectlOptions
	namespace   string
}

// CreateK8ssandraOptionsContext creates an options context model for k8ssandra test usage.
func CreateK8ssandraOptionsContext(
	operatorOptions *helm.Options, clusterOptions *helm.Options) K8ssandraOptionsContext {

	operatorCtx := OptionsContext{
		helmOptions: operatorOptions,
		kubeOptions: operatorOptions.KubectlOptions,
		namespace:   operatorOptions.KubectlOptions.Namespace}

	clusterCtx := OptionsContext{
		helmOptions: clusterOptions,
		kubeOptions: clusterOptions.KubectlOptions,
		namespace:   clusterOptions.KubectlOptions.Namespace}

	return K8ssandraOptionsContext{Operator: operatorCtx, Cluster: clusterCtx}
}

// InstallChart chart with specific release.
func InstallChart(ctx OptionsContext, chartName string, releaseName string, additionalArgs string) {

	var result string
	var err error

	if additionalArgs != "" {
		result, err = helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.helmOptions, "install", "-n",
			ctx.namespace, releaseName, chartName, additionalArgs)
	} else {
		result, err = helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.helmOptions, "install", "-n",
			ctx.namespace, releaseName, chartName)
	}
	Log("Helm", "Install", result, err)

	Ω(err).Should(BeNil())
	k8s.WaitUntilAllNodesReady(GinkgoT(), ctx.kubeOptions, 15, 3*time.Second)
}

// UninstallRelease release target to be uninstalled.
func UninstallRelease(ctx OptionsContext, releaseName string) {

	result, err := helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.helmOptions, "uninstall", "-n",
		ctx.namespace, releaseName, "--no-hooks")
	Log("Helm", "Uninstall", result, err)

	Pause("Pause for Uninstall completion", 1, 6*time.Second)
}

// LookupReleases constructs results of current release details.  nil returned upon error condition.
func LookupReleases(ctx OptionsContext) []ReleaseDetail {

	var releaseDetail []ReleaseDetail
	result, err := helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.helmOptions,
		"ls", "-n", ctx.namespace, "-o", "json")
	if err != nil {
		return nil
	}
	unmarshalErr := json.Unmarshal([]byte(result), &releaseDetail)
	Ω(unmarshalErr).Should(BeNil())

	return releaseDetail
}

// IsReleaseDeployed indicates status of a named release
func IsReleaseDeployed(ctx OptionsContext, name string) bool {

	for _, rel := range LookupReleases(ctx) {
		if rel.Status == "deployed" && rel.Name == name {
			return true
		}
	}
	return false
}

// IsPodRunning indicates if a running pod exists and has the prefix supplied.
func IsPodRunning(ctx OptionsContext, prefix string) bool {

	pods := LookupRunningPods(ctx)
	if pods == nil {
		return false
	}
	for _, pod := range pods {
		if prefix != "" && strings.Contains(pod, prefix) {
			return true
		}
	}
	return false
}

// IsLabeledPodExisting
func IsLabeledPodExisting(ctx OptionsContext, label string) bool {
	pods := LookupPodsByLabel(ctx, label)
	return pods != nil && len(pods) > 0
}

// LookupPodsByLabel lookup pods by label
func LookupPodsByLabel(ctx OptionsContext, label string) []string {
	result, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), ctx.kubeOptions, "get", "pods", "-n",
		ctx.namespace, "--show-labels", "-l", "name="+label, "--no-headers=true")
	Ω(err).Should(BeNil())

	if strings.Contains(result, "No resources found") {
		return nil
	}
	return strings.Split(result, "\n")
}

// LookupRunningPods provides current state of pods scoped by namespace.
func LookupRunningPods(ctx OptionsContext) []string {

	result, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), ctx.kubeOptions, "get", "pods", "-n",
		ctx.namespace, "--field-selector=status.phase=Running", "-o", "name", "--no-headers=true")
	Ω(err).Should(BeNil())

	if strings.Contains(result, "No resources found") {
		return nil
	}
	return strings.Split(result, "\n")
}

// Log helper for output of action details.
func Log(subject string, action string, result string, errorDetail error) {

	if errorDetail != nil {
		logger.Log(GinkgoT(), fmt.Sprintf("[ %s ] [ %s] [ Error: %s]", subject, action, errorDetail.Error()))
	} else {
		if result != "" {
			logger.Log(GinkgoT(), fmt.Sprintf("[ %s ] [ %s ] [ %s]", subject, action, result))
		} else {
			logger.Log(GinkgoT(), fmt.Sprintf("[ %s ] [ %s ]", subject, action))
		}
	}
}

// LogResults helper for output of map-based result details.
func LogResults(message string, results map[string]string) {
	for k, v := range results {
		logger.Log(GinkgoT(), fmt.Sprintf("%s Id: %s Result: %s", message, k, v))
	}
}

// Pause helper responsible for pausing for a specified time duration and interval.
func Pause(message string, interval, timeout time.Duration) {

	logger.Log(GinkgoT(), fmt.Sprintf("%s Timeout: %s", message, timeout))
	wait.Poll(interval, timeout, func() (bool, error) {
		if true {
			time.Sleep(1 * time.Second)
		}
		return false, nil
	})
}
