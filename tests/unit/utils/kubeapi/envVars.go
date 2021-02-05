package kubeapi

import (
	coreV1 "k8s.io/api/core/v1"
)

func FindEnvVarByName(envVars []coreV1.EnvVar, name string) *coreV1.EnvVar {
	for _, candidate := range envVars {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}
