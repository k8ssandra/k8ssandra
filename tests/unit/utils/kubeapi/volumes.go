package kubeapi

import (
	corev1 "k8s.io/api/core/v1"
)

func GetVolumeMountNames(container *corev1.Container) []string {
	names := make([]string, 0)
	for _, mount := range container.VolumeMounts {
		names = append(names, mount.Name)
	}
	return names
}

func GetVolumeNames(podTemplateSpec *corev1.PodTemplateSpec) []string {
	names := make([]string, 0)
	for _, volume := range podTemplateSpec.Spec.Volumes {
		names = append(names, volume.Name)
	}
	return names
}
