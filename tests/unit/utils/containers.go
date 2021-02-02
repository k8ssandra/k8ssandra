package utils

import (
	cassdcV1beta1 "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
)

func GetInitContainer(cassdc *cassdcV1beta1.CassandraDatacenter, name string) *coreV1.Container {
	return GetContainerByName(cassdc.Spec.PodTemplateSpec.Spec.InitContainers, name)

}

func GetContainer(cassdc *cassdcV1beta1.CassandraDatacenter, name string) *coreV1.Container {
	return GetContainerByName(cassdc.Spec.PodTemplateSpec.Spec.Containers, name)

}

func AssertInitContainerNamesMatch(cassdc *cassdcV1beta1.CassandraDatacenter, names ...string) {
	initContainers := cassdc.Spec.PodTemplateSpec.Spec.InitContainers
	actualNames := GetContainerNames(initContainers)

	ExpectWithOffset(1, actualNames).To(Equal(names))
}

func AssertContainerNamesMatch(cassdc *cassdcV1beta1.CassandraDatacenter, names ...string) {
	containers := cassdc.Spec.PodTemplateSpec.Spec.Containers
	actualNames := GetContainerNames(containers)

	ExpectWithOffset(1, actualNames).To(Equal(names))
}

func GetContainerByName(containers []coreV1.Container, name string) *coreV1.Container {
	for _, container := range containers {
		if container.Name == name {
			return &container
		}
	}
	return nil
}

func GetContainerNames(containers []coreV1.Container) []string {
	names := make([]string, 0)
	for _, container := range containers {
		names = append(names, container.Name)
	}
	return names
}
