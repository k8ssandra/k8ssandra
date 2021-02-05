package kubeapi

import (
	coreV1 "k8s.io/api/core/v1"
)

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
