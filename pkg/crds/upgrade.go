package crds

import (
	"bytes"
	"context"
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
	extractDir, err := helmutil.DownloadChartRelease(targetVersion)
	if err != nil {
		return nil, err
	}

	// reaper and medusa subdirs have the required yaml files
	chartPath := filepath.Join(extractDir, helmutil.ChartName)
	defer os.RemoveAll(chartPath)

	crds := make([]unstructured.Unstructured, 0)

	// For each dir under the charts subdir, check the "crds/"
	paths, _ := findCRDDirs(chartPath)

	for _, path := range paths {
		err = parseChartCRDs(&crds, path)
		if err != nil {
			return nil, err
		}
	}

	for _, obj := range crds {
		existingCrd := obj.DeepCopy()
		log.Printf("Finding %s", obj.GetName())
		err = u.client.Get(context.TODO(), client.ObjectKey{Name: obj.GetName()}, existingCrd)
		if apierrors.IsNotFound(err) {
			log.Printf("Creating %v\n", obj.GetName())
			if err = u.client.Create(context.TODO(), &obj); err != nil {
				log.Fatalf("Failed to create %s: %v\n", obj.GetName(), err)
			}
		} else if err == nil {
			log.Printf("Updating %v\n", obj.GetName())
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

		dec := deser.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

		if len(b) == 0 {
			log.Printf("Skipping %s\n", path)
			// TODO Implement skipping this, it's a YAML with "---" starter
			return nil
		}

		crd := unstructured.Unstructured{}
		_, gvk, err := dec.Decode(b, nil, &crd)
		if err != nil {
			return err
		}

		if gvk.Kind != "CustomResourceDefinition" {
			return nil
		}

		*crds = append(*crds, crd)
		// }

		// reader := k8syaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(b)))

		_ = k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(b), 4096)

		// doc, err := reader.Read()
		// // log.Printf("Doc read: %s\n", string(doc))
		// if err != nil {
		// 	return err
		// }
		// var obj runtime.Object
		// // ext := runtime.RawExtension{}

		// err = decoder.Decode(obj)
		// if err != nil {
		// 	return err
		// }

		// log.Printf("RAW: %s\n", string(ext.Raw))

		// TODO Single crd could include multiple objects.. we need to check if we actually got anything or do we want to move forward

		// if err = yaml.Unmarshal(doc, &crd.Object); err != nil {
		// 	return err
		// }

		// log.Printf("Read input: %s\n", crd.GetName())

		return err
	})

	return errOuter
}
