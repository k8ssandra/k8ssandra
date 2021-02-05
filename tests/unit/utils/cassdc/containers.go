package cassdc

import (
	cassop "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/k8ssandra/k8ssandra/tests/unit/utils/kubeapi"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func GetInitContainer(cassdc *cassop.CassandraDatacenter, name string) *corev1.Container {
	return GetContainerByName(cassdc.Spec.PodTemplateSpec.Spec.InitContainers, name)

}

func GetContainer(cassdc *cassop.CassandraDatacenter, name string) *corev1.Container {
	return GetContainerByName(cassdc.Spec.PodTemplateSpec.Spec.Containers, name)

}

func AssertInitContainerNamesMatch(cassdc *cassop.CassandraDatacenter, names ...string) {
	initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
	actualNames := GetContainerNames(initContainers)

	ExpectWithOffset(1, actualNames).To(Equal(names))
}

func AssertContainerNamesMatch(cassdc *cassop.CassandraDatacenter, names ...string) {
	containers := cassdc.Spec.PodTemplateSpec.Spec.Containers
	actualNames := GetContainerNames(containers)

	ExpectWithOffset(1, actualNames).To(Equal(names))
}
