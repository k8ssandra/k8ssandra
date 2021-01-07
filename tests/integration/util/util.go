package util

import (
	"encoding/json"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/shell"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/wait"
)

type waitCondition func() bool

var (
	newlineSeparator = "\n"
	deployedStatus   = "deployed"
)

func CreateK8ssandraOptions(kubeOptions *k8s.KubectlOptions, clusterReleaseName string,
	operatorReleaseName string, networkReleaseName string) K8ssandraOptions {

	k8ssandraOptions := NewK8ssandraOptions()
	k8ssandraOptions.SetCluster(&helm.Options{KubectlOptions: kubeOptions},
		clusterReleaseName)
	k8ssandraOptions.SetOperator(&helm.Options{KubectlOptions: kubeOptions},
		operatorReleaseName)
	k8ssandraOptions.SetNetwork(&helm.Options{KubectlOptions: kubeOptions},
		networkReleaseName)
	options := k8ssandraOptions.GetK8ssandraOptions()

	Expect(options).NotTo(BeNil())
	Expect(options.Operator).NotTo(BeNil())
	Expect(options.Cluster).NotTo(BeNil())
	Expect(options.Network).NotTo(BeNil())

	Expect(options.GetClusterCtx()).NotTo(BeNil())
	Expect(options.GetNetworkCtx()).NotTo(BeNil())
	Expect(options.GetOperatorCtx()).NotTo(BeNil())

	return options
}

// InstallChart chart with specific release.
func InstallChart(ctx OptionsContext, chartName string) {

	cmdArgs := append(ctx.Args, "-n", ctx.Namespace, ctx.ReleaseName, chartName, "--create-namespace")
	result, err := helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.HelmOptions, "install", cmdArgs...)
	Log("Helm", "Install", result, err)
	Expect(err).Should(BeNil())
}

// InstallChartFromFile
func InstallChartWithValues(ctx OptionsContext, chartName string, valuesTarget string) {

	cmdArgs := append(ctx.Args, "-n", ctx.Namespace, ctx.ReleaseName, chartName, "-f", valuesTarget, "--create-namespace")
	result, err := helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.HelmOptions, "install", cmdArgs...)
	Log("Helm", "Install", result, err)
	Expect(err).Should(BeNil())
}

// LookupReleases constructs results of current release details.  nil returned upon error condition.
func LookupReleases(ctx OptionsContext) []ReleaseDetail {

	var releaseDetail []ReleaseDetail
	result, err := helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.HelmOptions,
		"ls", "-n", ctx.Namespace, "-o", "json")
	if err != nil {
		return nil
	}
	unmarshalErr := json.Unmarshal([]byte(result), &releaseDetail)
	Expect(unmarshalErr).Should(BeNil())
	return releaseDetail
}

func IsDeploymentReady(ctx OptionsContext, name string) bool {

	result, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), ctx.KubeOptions,
		"get", "deployment", "-A",
		"--field-selector", fmt.Sprintf("metadata.namespace=%s", ctx.Namespace),
		"--field-selector", fmt.Sprintf("metadata.name=%s", name),
		"-o", "name")

	Log("Cluster Deployment", "Verify", result, err)
	return err == nil && strings.TrimPrefix(result, "deployment.apps/") == ctx.Namespace
}

func IsPodWithLabel(ctx OptionsContext, podLabels []PodLabel) bool {

	var cmd = []string{"get"}
	var args = append(cmd, "pod", "-n", ctx.Namespace)
	for _, podLabel := range podLabels {
		args = append(args, "-l", labelKV(podLabel))
	}
	args = append(args, "-o", "name")
	result, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), ctx.KubeOptions, args...)
	return err == nil && strings.HasPrefix(result, "pod/")
}

func IsClusterExisting(clusterName string) bool {

	var command = shell.Command{Command: "kind", Args: []string{"get", "clusters"}}
	result, err := shell.RunCommandAndGetOutputE(GinkgoT(), command)

	if err != nil || strings.Contains(result, "No kind clusters found") {
		return false
	}

	clusters := strings.Split(result, newlineSeparator)
	for _, cluster := range clusters {
		if strings.Contains(cluster, clusterName) {
			return true
		}
	}
	return false
}

// WaitFor tests wait condition up to limit of timeout seconds provided.
// Condition checked at least once at interval provided.
func WaitFor(condition waitCondition, message string, intervalSecs int, timeoutSecs int) {

	timeout := time.Duration(timeoutSecs) * time.Second
	interval := time.Duration(intervalSecs) * time.Second

	logger.Log(GinkgoT(), fmt.Sprintf("%s Timeout: %s Interval: %s ", message, timeout.String(), interval.String()))
	err := wait.Poll(interval, timeout, func() (bool, error) {
		if condition() {
			return true, nil
		}
		time.Sleep(1 * time.Second)
		return false, nil
	})
	if err != nil {
		Log("WaitFor", "Timeout", "Error", err)
	}
}

// IsReleaseDeployed indicates status of a named release
func IsReleaseDeployed(ctx OptionsContext) bool {

	releaseName := ctx.ReleaseName
	for _, rel := range LookupReleases(ctx) {
		if rel.Status == deployedStatus && rel.Name == releaseName {
			return true
		}
	}
	return false
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

// Pause helper responsible for pausing for a specified time duration and interval.
func Pause(message string, interval, timeout time.Duration) {

	logger.Log(GinkgoT(), fmt.Sprintf("%s Timeout: %s", message, timeout))
	_ = wait.Poll(interval, timeout, func() (bool, error) {
		if true {
			time.Sleep(1 * time.Second)
		}
		return false, nil
	})
}

// RepoUpdate provides a helm repo update as sometimes needed to pickup changes in repo.
func RepoUpdate(ctx OptionsContext) bool {

	updateResult, updateErr := helm.RunHelmCommandAndGetOutputE(GinkgoT(), ctx.HelmOptions,
		"repo", "update")
	Log("Repo", "Update", updateResult, updateErr)
	Pause("Pausing after repo update", 1, 5)
	return updateErr == nil
}

func ConfigureContext(ctx OptionsContext, configFile string) {
	err := k8s.RunKubectlE(GinkgoT(), ctx.KubeOptions,
		"config", "--kubeconfig="+configFile, "use-context", ctx.KubeOptions.ContextName)
	Log("Config", "Use-Context: ", ctx.KubeOptions.ContextName, err)
	Expect(err).To(BeNil())
}

func IsExisting(ctx OptionsContext, kind string, expectedName string) bool {

	result, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), ctx.KubeOptions,
		"get", kind, "-o", "name", "--no-headers=true")
	Expect(err).Should(BeNil())

	if result != "" {
		if expectedName == "" || result == expectedName {
			return true
		}
	}
	return false
}

func CreateKindCluster(detail KindClusterDetail) {

	var command = shell.Command{Command: "kind", Args: []string{"create", "cluster",
		"--name", detail.Name, "--config", detail.ConfigFile, "--image", detail.Image}}
	result, err := shell.RunCommandAndGetOutputE(GinkgoT(), command)
	Log("Cluster", "create", result, err)
	WaitFor(func() bool { return IsClusterExisting(detail.Name) },
		"Wait for cluster create", 3, 120)
}

func DeleteKindCluster(name string) {

	var command = shell.Command{Command: "kind", Args: []string{"delete", "cluster", "--name", name}}
	result, err := shell.RunCommandAndGetOutputE(GinkgoT(), command)
	Log("delete", "cluster", result, err)
	WaitFor(func() bool { return !IsClusterExisting(name) },
		"Wait for cluster cleanup", 3, 30)
}

// labelKV extraction of k/v for pod label
func labelKV(podLabel PodLabel) string {
	Expect(podLabel).ShouldNot(BeNil())
	Expect(podLabel.Key).ShouldNot(BeEmpty())
	Expect(podLabel.Value).ShouldNot(BeNil())
	return podLabel.Key + "=" + podLabel.Value
}
