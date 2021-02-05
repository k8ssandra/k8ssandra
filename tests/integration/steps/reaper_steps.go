package steps

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	resty "github.com/go-resty/resty/v2"
	. "github.com/onsi/gomega"
)

// Reaper related steps
func CheckClusterIsRegisteredInReaper(t *testing.T, clusterName string) {
	g(t).Eventually(func() bool {
		return clusterIsRegisteredInReaper(t, clusterName)
	}, RETRY_TIMEOUT, RETRY_INTERVAL).Should(BeTrue(), "Cluster wasn't properly registered in Reaper")
}

func clusterIsRegisteredInReaper(t *testing.T, clusterName string) bool {
	restClient := resty.New()
	response, err := restClient.R().Get("http://repair.127.0.0.1.nip.io:8080/cluster")
	if err != nil {
		log.Println(fmt.Sprintf("The HTTP request failed with error %s", err))
	} else {
		data := response.Body()
		log.Println(fmt.Sprintf("Reaper response: %s", data))
		var clusters []string
		json.Unmarshal([]byte(data), &clusters)
		if len(clusters) > 0 {
			if clusterName == clusters[0] {
				return true
			}
		}
	}
	return false
}

func CancelRepair(t *testing.T, repairId string) {
	restClient := resty.New()
	// Start the previously created repair run
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		Put(fmt.Sprintf("http://repair.127.0.0.1.nip.io:8080/repair_run/%s/state/ABORTED", repairId))

	log.Println(fmt.Sprintf("Reaper response: %s", response.Body()))
	log.Println(fmt.Sprintf("Reaper status code: %d", response.StatusCode()))

	errMessage := fmt.Sprintf("Failed aborting repair %s: %s / %s", repairId, err, response.Body())
	g(t).Expect(err).To(BeNil(), errMessage)
	g(t).Expect(response.StatusCode()).Should(Equal(200), errMessage)
}

// Starts a repair on keyspace and return the repair id
func TriggerRepair(t *testing.T, namespace, keyspace string) string {
	restClient := resty.New()

	// Create the repair run
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"clusterName":  "k8ssandra",
			"keyspace":     keyspace,
			"owner":        "k8ssandra",
			"segmentCount": "5",
		}).
		Post("http://repair.127.0.0.1.nip.io:8080/repair_run")

	data := response.Body()
	log.Println(fmt.Sprintf("Reaper response: %s", data))
	var reaperResponse interface{}
	err2 := json.Unmarshal(data, &reaperResponse)

	errMessageCreateRepair := fmt.Sprintf("The REST request or response parsing failed with error %s %s: %s", err, err2, data)
	g(t).Expect(err).To(BeNil(), errMessageCreateRepair)
	g(t).Expect(err2).To(BeNil(), errMessageCreateRepair)

	reaperResponseMap := reaperResponse.(map[string]interface{})
	repairId := fmt.Sprintf("%s", reaperResponseMap["id"])
	// Start the previously created repair run
	response, err = restClient.R().
		SetHeader("Content-Type", "application/json").
		Put(fmt.Sprintf("http://repair.127.0.0.1.nip.io:8080/repair_run/%s/state/RUNNING", repairId))

	log.Println(fmt.Sprintf("Reaper response: %s", response.Body()))
	log.Println(fmt.Sprintf("Reaper status code: %d", response.StatusCode()))

	errMessageStart := fmt.Sprintf("Failed starting repair %s: %s / %s", repairId, err, response.Body())
	g(t).Expect(err).To(BeNil(), errMessageStart)
	g(t).Expect(response.StatusCode()).Should(Equal(200), errMessageStart)

	return repairId
}

func WaitForOneSegmentToBeDone(t *testing.T, repairId string) {
	restClient := resty.New()
	g(t).Eventually(func() bool {
		return oneSegmentIsDone(t, repairId, restClient)
	}, RETRY_TIMEOUT, RETRY_INTERVAL).Should(BeTrue(), "No repair segment was fully processed within timeout")
}

func oneSegmentIsDone(t *testing.T, repairId string, restClient *resty.Client) bool {
	response, err := restClient.R().Get(fmt.Sprintf("http://repair.127.0.0.1.nip.io:8080/repair_run/%s/segments", repairId))
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("The HTTP request failed with error %s", err))

	return strings.Contains(fmt.Sprintf("%s", response.Body()), "\"state\":\"DONE\"")
}
