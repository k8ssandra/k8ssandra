package cassdc

import (
	cassop "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

// GetInitContainer finds an initContainer with the given in the podTemplateSpec of the given CassandraDatacenter
func GetInitContainer(cassdc *cassop.CassandraDatacenter, name string) *corev1.Container {
	return GetContainerByName(cassdc.Spec.PodTemplateSpec.Spec.InitContainers, name)
}

// GetContainer finds a container with the given in the podTemplateSpec of the given CassandraDatacenter
func GetContainer(cassdc *cassop.CassandraDatacenter, name string) *corev1.Container {
	return GetContainerByName(cassdc.Spec.PodTemplateSpec.Spec.Containers, name)

}

// AssertInitContainerNamesMatch asserts that the names of the initContainers defined in the podTemplateSpec
// of the given CassandraDatacenter match the given names
func AssertInitContainerNamesMatch(cassdc *cassop.CassandraDatacenter, names ...string) {
	initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
	actualNames := GetContainerNames(initContainers)

	ExpectWithOffset(1, actualNames).To(Equal(names))
}

// AssertContainerNamesMatch asserts that the names of the containers defined in the podTemplateSpec
// of the given CassandraDatacenter match the given names
func AssertContainerNamesMatch(cassdc *cassop.CassandraDatacenter, names ...string) {
	containers := cassdc.Spec.PodTemplateSpec.Spec.Containers
	actualNames := GetContainerNames(containers)

	ExpectWithOffset(1, actualNames).To(Equal(names))
}
