package main

import (
	"flag"
	"log"
	"os"

	"github.com/k8ssandra/k8ssandra/pkg/cleaner"
	"github.com/k8ssandra/k8ssandra/pkg/upgrade"
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
	var releaseVersion string
	flag.StringVar(&releaseName, "release", "", "Defines the releaseName to be cleaned")
	cleanResources := flag.Bool("clean", false, "Clean resources with finalizers")
	upgradeCRDs := flag.Bool("update-crds", false, "Upgrade CRDs during helm upgrade")
	flag.StringVar(&releaseVersion, "version", "", "target version to upgrade to")
	flag.Parse()

	// Add flags for parsing stuff
	if *cleanResources {
		log.Printf("Cleaning resources for uninstall")

		ca, err := cleaner.New(namespace)
		if err != nil {
			log.Fatalf("Failed to create new cleaner: %v", err)
		}

		err = ca.RemoveResources(releaseName)
		if err != nil {
			log.Fatalf("Failed to remove resources: %v", err)
		}
	}

	if *upgradeCRDs {
		u, err := upgrade.New()
		if err != nil {
			log.Fatalf("Failed to create new CRD upgrader: %v", err)
		}

		err = u.Upgrade(releaseVersion)
		if err != nil {
			log.Fatalf("Failed to upgrade CRDs: %v", err)
		}
	}
}
