package main

import (
	"flag"
	"log"
	"os"

	"github.com/k8ssandra/k8ssandra/pkg/cleaner"
)

var (
	podNameSpaceEnvVar = "POD_NAMESPACE"
)

func main() {
	namespace := os.Getenv(podNameSpaceEnvVar)
	if namespace == "" {
		log.Fatalf("Failed to parse pod's namespace from env variable %s", podNameSpaceEnvVar)
		return
	}

	var releaseName string
	flag.StringVar(&releaseName, "release", "", "Defines the releaseName to be cleaned")
	flag.Parse()

	// Add flags for parsing stuff
	ca, err := cleaner.New(namespace)
	if err != nil {
		log.Fatalf("Failed to create new cleaner: %v", err)
	}

	err = ca.RemoveResources(releaseName)
	if err != nil {
		log.Fatalf("Failed to remove resources: %v", err)
	}
}
