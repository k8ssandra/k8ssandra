package helmutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	ManagedLabel      = "app.kubernetes.io/managed-by"
	ManagedLabelValue = "Helm"
	ReleaseAnnotation = "meta.helm.sh/release-name"

	RepoURL          = "https://helm.k8ssandra.io/"
	DefaultChartName = "k8ssandra"
)

func GetChartTargetDir(targetVersion string) (string, error) {
	// Extract the files
	subDir := filepath.Join("helm", targetVersion)
	extractDir, err := GetCacheDir(subDir)
	if err != nil {
		return "", err
	}

	return extractDir, err
}

// DownloadChartRelease fetches the k8ssandra target version and extracts it to a directory which path is returned
func DownloadChartRelease(chartName string, targetVersion string) (string, error) {
	// Unfortunately, the helm's chart pull command uses "internal" marked structs, so it can't be used for
	// pulling the data. Thus, we need to replicate the implementation here and use our own cache

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
		URL:  RepoURL,
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

	url, err := repo.ResolveReferenceURL(RepoURL, cv.URLs[0])
	if err != nil {
		return "", err
	}

	// Download to filesystem for extraction purposes
	dir, err := ioutil.TempDir("", "helmutil-")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(dir)

	// _ is ProvenanceVerify (TODO we might want to verify the release)
	saved, _, err := c.DownloadTo(url, targetVersion, dir)
	if err != nil {
		return "", err
	}

	// extractDir is a target directory
	extractDir, err := GetChartTargetDir(targetVersion)
	if err != nil {
		return "", nil
	}

	_, err = CreateIfNotExistsDir(extractDir)
	if err != nil {
		return "", nil
	}

	err = chartutil.ExpandFile(extractDir, saved)
	if err != nil {
		return "", err
	}

	return extractDir, nil
}
