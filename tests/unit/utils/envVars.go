package utils

import (
	coreV1 "k8s.io/api/core/v1"
)

func FindEnvVarByName(haystack []coreV1.EnvVar, needle string) *coreV1.EnvVar {
	for _, candidate := range haystack {
		if candidate.Name == needle {
			return &candidate
		}
	}
	return nil
}
