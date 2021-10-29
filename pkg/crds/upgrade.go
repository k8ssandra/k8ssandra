package crds

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/k8ssandra/k8ssandra/pkg/helmutil"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	deser "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Upgrader is a utility to update the CRDs in a helm chart's pre-upgrade hook
type Upgrader struct {
	client    client.Client
	namespace string
}

// NewWithClient returns a new Upgrader client using the given controller-runtime client.Client
func NewWithClient(c client.Client) (*Upgrader, error) {
	return &Upgrader{
		client: c,
	}, nil
}

// New returns a new Upgrader client
func New(namespace string) (*Upgrader, error) {
	_ = api.AddToScheme(scheme.Scheme)
	_ = apiextv1.AddToScheme(scheme.Scheme)
	_ = apiextv1beta1.AddToScheme(scheme.Scheme)
	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Upgrader{
		client:    c,
		namespace: namespace,
	}, nil
}

// Upgrade installs the missing CRDs or updates them if they exists already
func (u *Upgrader) Upgrade(targetVersion string) ([]unstructured.Unstructured, error) {
	extractDir, err := helmutil.GetChartTargetDir(targetVersion)
	if err != nil {
		return nil, err
	}

	// If the targetCacheDirectory does not exist, download the chart
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		log.Printf("Downloading release %s from Helm repository", targetVersion)
		extractDir, err = helmutil.DownloadChartRelease(targetVersion)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	crds := make([]unstructured.Unstructured, 0)

	// For each dir under the charts subdir, check the "crds/"
	paths, _ := findCRDDirs(extractDir)

	for _, path := range paths {
		err = parseChartCRDs(&crds, path)
		if err != nil {
			return nil, err
		}
	}

	for _, obj := range crds {
		existingCrd := obj.DeepCopy()
		err = u.client.Get(context.TODO(), client.ObjectKey{Name: obj.GetName()}, existingCrd)
		if apierrors.IsNotFound(err) {
			log.Printf("Creating %v\n", obj.GetName())
			if err = u.client.Create(context.TODO(), &obj); err != nil {
				log.Fatalf("Failed to create %s: %v\n", obj.GetName(), err)
			}
		} else if err == nil {
			log.Printf("Updating %v\n", obj.GetName())
			if obj.GetName() == "CassandraDatacenter" {
				// We might need to patch the CRD before upgrade from the server.. only if it's v1beta1!
				if existingCrd.GetAPIVersion() == "v1beta1" {
					var crd apiextv1beta1.CustomResourceDefinition
					if err = runtime.DefaultUnstructuredConverter.FromUnstructured(existingCrd.UnstructuredContent(), crd); err != nil {
						log.Fatalf("Failed to cast CRD to CustomResourceDefinition (v1beta1): %v", err)
						return nil, err
					}

					// We don't care if err != nil, we shouldn't try to update v1 CRDs
					*crd.Spec.PreserveUnknownFields = false
					if err = u.client.Update(context.TODO(), existingCrd); err != nil {
						log.Fatalf("Failed to set preserveUnknownFields to false: %v", err)
						return nil, err
					}
				}
			}

			obj.SetResourceVersion(existingCrd.GetResourceVersion())
			if err = u.client.Update(context.TODO(), &obj); err != nil {
				log.Fatalf("Failed to update %s: %v\n", obj.GetName(), err)
				return nil, err
			}
		} else {
			log.Fatalf("Failed to Get the object %s: %v\n", obj.GetName(), err)
			return nil, err
		}
	}

	return crds, err
}

func findCRDDirs(chartDir string) ([]string, error) {
	dirs := make([]string, 0)
	err := filepath.Walk(chartDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasSuffix(path, "crds") {
				dirs = append(dirs, path)
			}
			return nil
		}
		return nil
	})
	return dirs, err
}

func parseChartCRDs(crds *[]unstructured.Unstructured, crdDir string) error {
	errOuter := filepath.Walk(crdDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Add to CRDs ..
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if len(b) == 0 {
			return nil
		}

		docs, err := parseCRDYamls(b)
		if err != nil {
			return err
		}
		dec := deser.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

		for _, b := range docs {
			crd := unstructured.Unstructured{}

			_, gvk, err := dec.Decode(b, nil, &crd)
			if err != nil {
				continue
			}

			if gvk.Kind != "CustomResourceDefinition" {
				continue
			}

			*crds = append(*crds, crd)
		}

		return err
	})

	return errOuter
}

func parseCRDYamls(b []byte) ([][]byte, error) {
	docs := [][]byte{}
	reader := k8syaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(b)))
	for {
		// Read document
		doc, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}
