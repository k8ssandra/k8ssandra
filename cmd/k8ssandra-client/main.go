package main

import (
	"flag"
	"log"
	"os"

	"github.com/k8ssandra/k8ssandra/pkg/cleaner"
	"github.com/k8ssandra/k8ssandra/pkg/crds"
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
	cleanResources := flag.Bool("clean", false, "Clean resources with finalizers")

	var targetVersion string
	flag.StringVar(&targetVersion, "targetVersion", "", "Defines the targetVersion to be upgraded to")
	upgradeCRDs := flag.Bool("upgradecrds", false, "Upgrade CRDs to target version")
	flag.Parse()

	// Add flags for parsing stuff
	if *cleanResources {
		log.Printf("Cleaning resources for uninstall")

		if releaseName == "" {
			log.Fatalf("No releaseName set")
			return
		}

		ca, err := cleaner.New(namespace)
		if err != nil {
			log.Fatalf("Failed to create new cleaner: %v", err)
			return
		}

		err = ca.RemoveResources(releaseName)
		if err != nil {
			log.Fatalf("Failed to remove resources: %v", err)
			return
		}
	}

	if *upgradeCRDs {
		log.Printf("Upgrading CRDs to version %s", targetVersion)
		if targetVersion == "" {
			log.Fatal("No targetVersion set")
			return
		}

		u, err := crds.New(namespace)
		if err != nil {
			log.Fatalf("Failed to create new CRD upgrader: %v", err)
			return
		}

		err = u.Upgrade(targetVersion)
		if err != nil {
			log.Fatalf("Failed to update CRDs: %v", err)
			return
		}
	}
}
