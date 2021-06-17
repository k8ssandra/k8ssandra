package steps

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/k8ssandra/reaper-client-go/reaper"
	"log"
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
)

// Reaper related steps

var reaperURL, _ = url.Parse("http://repair.127.0.0.1.nip.io:8080")
var reaperClient = reaper.NewClient(reaperURL)

func CheckClusterIsRegisteredInReaper(t *testing.T, clusterName string) {
	g(t).Eventually(func() bool {
		return clusterIsRegisteredInReaper(clusterName)
	}, retryTimeout, retryInterval).Should(BeTrue(), "Cluster wasn't properly registered in Reaper")
}

func clusterIsRegisteredInReaper(clusterName string) bool {
	if _, err := reaperClient.GetCluster(context.Background(), clusterName); err != nil {
		log.Println(fmt.Errorf("cluster %s is not registered: %w", clusterName, err))
		return false
	}
	return true
}

func CancelRepair(t *testing.T, repairId uuid.UUID) {
	err := reaperClient.AbortRepairRun(context.Background(), repairId)
	g(t).Expect(err).To(BeNil(), "Failed to abort repair run %s: %s", repairId, err)
}

// TriggerRepair starts a repair on keyspace and return the repair id
func TriggerRepair(t *testing.T, clusterName, keyspace, owner string) uuid.UUID {
	options := &reaper.RepairRunCreateOptions{SegmentCountPerNode: 5}
	repairId, err := reaperClient.CreateRepairRun(context.Background(), clusterName, keyspace, owner, options)
	g(t).Expect(err).To(BeNil(), "Failed to create repair run: %s", err)
	// Start the previously created repair run
	err = reaperClient.StartRepairRun(context.Background(), repairId)
	g(t).Expect(err).To(BeNil(), "Failed to start repair run %s: %s", repairId, err)
	return repairId
}

func WaitForOneSegmentToBeDone(t *testing.T, repairId uuid.UUID) {
	g(t).Eventually(func() bool {
		return oneSegmentIsDone(t, repairId)
	}, retryTimeout, retryInterval).Should(BeTrue(), "No repair segment was fully processed within timeout")
}

func oneSegmentIsDone(t *testing.T, repairId uuid.UUID) bool {
	segments, err := reaperClient.RepairRunSegments(context.Background(), repairId)
	g(t).Expect(err).To(BeNil(), "Failed to get segments of repair run %s: %s", repairId, err)
	for _, segment := range segments {
		if segment.State == reaper.RepairSegmentStateDone {
			return true
		}
	}
	return false
}
