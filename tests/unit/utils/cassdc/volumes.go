package cassdc

import (
	cassop "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/gomega"
)

// AssertVolumeNamesMatch asserts that the names of the volumes defined in the podTemplateSpec
// of the given CassandraDatacenter match the given names
func AssertVolumeNamesMatch(cassdc *cassop.CassandraDatacenter, names ...string) {
	podTemplateSpec := cassdc.Spec.PodTemplateSpec
	actualNames := GetVolumeNames(podTemplateSpec)

	ExpectWithOffset(1, actualNames).To(ConsistOf(names))
}
