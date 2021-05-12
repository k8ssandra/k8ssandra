package crds

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/k8ssandra/k8ssandra/pkg/helmutil"
	"gopkg.in/yaml.v3"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"

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
func (u *Upgrader) Upgrade(targetVersion string) error {
	extractDir, err := helmutil.DownloadChartRelease(targetVersion)
	if err != nil {
		return err
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
			return err
		}
	}

	var res []runtime.Object
	for _, obj := range crds {
		res = append(res, &obj)
	}

	for _, obj := range res {
		if u.client.Create(context.TODO(), obj); err != nil {
			if apierrors.IsAlreadyExists(err) {
				if u.client.Update(context.TODO(), obj); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return err
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

		reader := k8syaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(b)))
		doc, err := reader.Read()
		if err != nil {
			return err
		}

		crd := unstructured.Unstructured{}

		if err = yaml.Unmarshal(doc, &crd); err != nil {
			return err
		}

		*crds = append(*crds, crd)
		return nil
	})

	return errOuter
}
