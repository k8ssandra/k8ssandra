package upgrade

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	repoURL   = "https://helm.k8ssandra.io/"
	chartName = "k8ssandra"
)

// Upgrader is a utility to update the CRDs in a helm chart's pre-upgrade hook
type Upgrader struct {
	client client.Client
}

// NewWithClient returns a new Upgrader client using the given controller-runtime client.Client
func NewWithClient(c client.Client) (*Upgrader, error) {
	return &Upgrader{
		client: c,
	}, nil
}

// New returns a new Upgrader client
func New() (*Upgrader, error) {
	_ = api.AddToScheme(scheme.Scheme)
	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Upgrader{
		client: c,
	}, nil
}

// Upgrade installs the missing CRDs or updates them if they exists already
func (u *Upgrader) Upgrade(targetVersion string) error {
	extractDir, err := u.helmDownloadChartRelease(targetVersion)
	if err != nil {
		return err
	}

	// reaper and medusa subdirs have the required yaml files
	crdDir := filepath.Join(extractDir, chartName, "crds/")
	defer os.RemoveAll(crdDir)

	crds := make([]unstructured.Unstructured, 0)

	err = u.parseChartCRDs(&crds, crdDir)

	var res []client.Object
	for _, obj := range crds {
		res = append(res, &obj)
	}

	for _, obj := range res {
		if u.client.Create(context.TODO(), obj); err != nil {
			if apierrors.IsAlreadyExists(err) {
				if u.client.Update(context.TODO(), obj); err != nil {
					return err
				}
			}
		}
	}

	return err
}

func (u *Upgrader) parseChartCRDs(crds *[]unstructured.Unstructured, crdDir string) error {
	err := filepath.Walk(crdDir, func(path string, info os.FileInfo, err error) error {
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

	return err
}

func (u *Upgrader) helmDownloadChartRelease(targetVersion string) (string, error) {
	settings := cli.New()
	var out strings.Builder

	c := downloader.ChartDownloader{
		Out: &out,
		// Keyring: p.Keyring,
		Verify:  downloader.VerifyNever,
		Getters: getter.All(settings),
		Options: []getter.Option{
			// getter.WithBasicAuth(p.Username, p.Password),
			// getter.WithTLSClientConfig(p.CertFile, p.KeyFile, p.CaFile),
			// getter.WithInsecureSkipVerifyTLS(p.InsecureSkipTLSverify),
		},
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
	}

	// helm repo add k8ssandra https://helm.k8ssandra.io/
	r, err := repo.NewChartRepository(&repo.Entry{
		Name: chartName,
		URL:  repoURL,
	}, getter.All(settings))

	if err != nil {
		return "", err
	}

	// helm repo update
	index, err := r.DownloadIndexFile()
	if err != nil {
		return "", err
	}

	// Read the index file for the repository to get chart information and return chart URL
	repoIndex, err := repo.LoadIndexFile(index)
	if err != nil {
		return "", err
	}

	// chart name, chart version
	cv, err := repoIndex.Get(chartName, targetVersion)
	if err != nil {
		return "", err
	}

	url, err := repo.ResolveReferenceURL(repoURL, cv.URLs[0])
	if err != nil {
		return "", err
	}

	// Download to filesystem or otherwise to a usable format
	dir, err := ioutil.TempDir("", "upgrade-")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(dir)

	// _ is ProvenanceVerify (we'll want to verify later)
	saved, _, err := c.DownloadTo(url, targetVersion, dir)
	if err != nil {
		return "", err
	}

	// Extract the files
	extractDir, err := ioutil.TempDir("", "upgrade-extract-")
	if err != nil {
		return "", err
	}

	// extractDir is a target directory
	err = chartutil.ExpandFile(extractDir, saved)
	if err != nil {
		return "", err
	}

	return extractDir, nil
}
