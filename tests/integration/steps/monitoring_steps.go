package steps

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	resty "github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
)

// Monitoring related functions
func CheckPrometheusActiveTargets(t *testing.T, expectedTargets int) {
	restClient := resty.New()
	response, err := restClient.R().Get("http://127.0.0.1:8080/prometheus/api/v1/targets")
	g(t).Expect(err).To(BeNil(), "Failed connecting to Prometheus")
	prometheusResponse := response.Body()
	var genericJson map[string]interface{}
	json.Unmarshal(prometheusResponse, &genericJson)
	g(t).Expect(genericJson["status"].(string)).Should(Equal("success"))
	g(t).Expect(len(genericJson["data"].(map[string]interface{})["activeTargets"].([]interface{}))).Should(Equal(expectedTargets))
}

func CheckPrometheusMetricExtraction(t *testing.T) {
	const metric = "scrape_duration_seconds"
	prometheusResponse := queryPrometheusMetric(t, metric)
	g(t).Expect(prometheusResponse["status"].(string)).Should(Equal("success"))
	log.Println(Info("Prometheus could be reached through HTTP"))
}

func queryPrometheusMetric(t *testing.T, metric string) map[string]interface{} {
	restClient := resty.New()
	response, err := restClient.R().Get(fmt.Sprintf("http://127.0.0.1:8080/prometheus/api/v1/query?query=%s", metric))
	g(t).Expect(err).To(BeNil(), "Failed connecting to Prometheus")
	prometheusResponse := response.Body()
	var genericJson map[string]interface{}
	json.Unmarshal(prometheusResponse, &genericJson)
	return genericJson
}

func CheckGrafanaIsReachable(t *testing.T) {
	restClient := resty.New()
	_, err := restClient.R().Get("http://localhost:8080/grafana/login")
	g(t).Expect(err).To(BeNil(), "Failed connecting to Grafana")
	log.Println(Info("Grafana could be reached through HTTP"))
}

func CountMonitoredItems(t *testing.T, namespace string) int {
	return CountPodsWithLabels(t, namespace, map[string]string{"app.kubernetes.io/managed-by": "cass-operator"}) +
		CountPodsWithLabels(t, namespace, map[string]string{"app": releaseName + "-" + datacenterName + "-stargate"})
}

func CheckNoOutOfOrderMetrics(t *testing.T, namespace string) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	prometheusPods := GetPodsWithLabels(t, namespace, map[string]string{"app.kubernetes.io/name": "prometheus"})
	g(t).Expect(len(prometheusPods.Items)).To(Equal(1), fmt.Sprintf("Expected one Prometheus pod but found %d", len(prometheusPods.Items)))
	prometheusLog, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "logs", prometheusPods.Items[0].Name, "-c", "prometheus")
	g(t).Expect(err).To(BeNil())
	g(t).Expect(prometheusLog).NotTo(ContainSubstring("Error on ingesting out-of-order samples"))
}

func CheckTableLevelMetricsArePresent(t *testing.T) {
	metrics := queryPrometheusMetric(t, "mcac_table_live_disk_space_used_total")
	require.IsType(t, map[string]interface{}{}, metrics["data"], "Expected field data to be map[string]interface{}")
	data := metrics["data"].(map[string]interface{})
	require.IsType(t, []interface{}{}, data["result"], "Expected field result to be []interface{}")
	result := data["result"].([]interface{})
	g(t).Expect(len(result)).To(BeNumerically(">", 0), "No table level metric was returned")
}
