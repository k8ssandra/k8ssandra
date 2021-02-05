package steps

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	resty "github.com/go-resty/resty/v2"
	. "github.com/onsi/gomega"
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
