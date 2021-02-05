package kubeapi

import (
	corev1 "k8s.io/api/core/v1"
)

func GetContainerByName(containers []corev1.Container, name string) *corev1.Container {
	for _, container := range containers {
		if container.Name == name {
			return &container
		}
	}
	return nil
}

func GetContainerNames(containers []corev1.Container) []string {
	names := make([]string, 0)
	for _, container := range containers {
		names = append(names, container.Name)
	}
	return names
}
