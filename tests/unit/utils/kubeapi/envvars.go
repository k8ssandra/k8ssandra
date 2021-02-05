package kubeapi

import (
	corev1 "k8s.io/api/core/v1"
)

// FindEnvVarByName finds an EnvVar with the given name in the given Container
func FindEnvVarByName(container corev1.Container, name string) *corev1.EnvVar {
	for _, candidate := range container.Env {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}
