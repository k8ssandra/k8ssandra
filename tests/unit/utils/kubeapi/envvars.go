package kubeapi

import (
	corev1 "k8s.io/api/core/v1"
)

func FindEnvVarByName(envVars []corev1.EnvVar, name string) *corev1.EnvVar {
	for _, candidate := range envVars {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}
