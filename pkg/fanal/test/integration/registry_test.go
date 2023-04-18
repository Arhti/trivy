package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/aquasecurity/trivy/pkg/fanal/analyzer/all"
	"github.com/aquasecurity/trivy/pkg/fanal/applier"
	"github.com/aquasecurity/trivy/pkg/fanal/artifact"
	aimage "github.com/aquasecurity/trivy/pkg/fanal/artifact/image"
	"github.com/aquasecurity/trivy/pkg/fanal/cache"
	"github.com/aquasecurity/trivy/pkg/fanal/image"
	testdocker "github.com/aquasecurity/trivy/pkg/fanal/test/integration/docker"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/types"
)

const (
	registryImage    = "registry:2"
	registryPort     = "5443/tcp"
	registryUsername = "testuser"
	registryPassword = "testpassword"
)

func TestTLSRegistry(t *testing.T) {
	ctx := context.Background()

	baseDir, err := filepath.Abs(".")
	require.NoError(t, err)

	req := testcontainers.ContainerRequest{
		Name:         "registry",
		Image:        registryImage,
		ExposedPorts: []string{registryPort},
		Env: map[string]string{
			"REGISTRY_HTTP_ADDR":            "0.0.0.0:5443",
			"REGISTRY_HTTP_TLS_CERTIFICATE": "/certs/cert.pem",
			"REGISTRY_HTTP_TLS_KEY":         "/certs/key.pem",
			"REGISTRY_AUTH":                 "htpasswd",
			"REGISTRY_AUTH_HTPASSWD_PATH":   "/auth/htpasswd",
			"REGISTRY_AUTH_HTPASSWD_REALM":  "Registry Realm",
		},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(filepath.Join(baseDir, "data", "registry", "certs"), "/certs"),
			testcontainers.BindMount(filepath.Join(baseDir, "data", "registry", "auth"), "/auth"),
		),
		WaitingFor: wait.ForLog("listening on [::]:5443"),
	}

	registryC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer registryC.Terminate(ctx)

	registryURL, err := getRegistryURL(ctx, registryC, registryPort)
	require.NoError(t, err)

	config := testdocker.RegistryConfig{
		URL:      registryURL,
		Username: registryUsername,
		Password: registryPassword,
	}

	testCases := []struct {
		name         string
		imageName    string
		imageFile    string
		option       ftypes.RemoteOptions
		login        bool
		expectedOS   ftypes.OS
		expectedRepo ftypes.Repository
		wantErr      bool
	}{
		{
			name:      "happy path",
			imageName: "ghcr.io/aquasecurity/trivy-test-images:alpine-310",
			imageFile: "../../../../integration/testdata/fixtures/images/alpine-310.tar.gz",
			option: ftypes.RemoteOptions{
				Credentials: []ftypes.Credential{
					{
						Username: registryUsername,
						Password: registryPassword,
					},
				},
				Insecure: true,
			},
			expectedOS: ftypes.OS{
				Name:   "3.10.2",
				Family: "alpine",
			},
			expectedRepo: ftypes.Repository{
				Family:  "alpine",
				Release: "3.10",
			},
			wantErr: false,
		},
		{
			name:      "happy path with docker login",
			imageName: "ghcr.io/aquasecurity/trivy-test-images:alpine-310",
			imageFile: "../../../../integration/testdata/fixtures/images/alpine-310.tar.gz",
			option: ftypes.RemoteOptions{
				Insecure: true,
			},
			login: true,
			expectedOS: ftypes.OS{
				Name:   "3.10.2",
				Family: "alpine",
			},
			expectedRepo: ftypes.Repository{
				Family:  "alpine",
				Release: "3.10",
			},
			wantErr: false,
		},
		{
			name:      "sad path: tls verify",
			imageName: "ghcr.io/aquasecurity/trivy-test-images:alpine-310",
			imageFile: "../../../../integration/testdata/fixtures/images/alpine-310.tar.gz",
			option: ftypes.RemoteOptions{
				Credentials: []ftypes.Credential{
					{
						Username: registryUsername,
						Password: registryPassword,
					},
				},
			},
			wantErr: true,
		},
		{
			name:      "sad path: no credential",
			imageName: "ghcr.io/aquasecurity/trivy-test-images:alpine-310",
			imageFile: "../../../../integration/testdata/fixtures/images/alpine-310.tar.gz",
			option: ftypes.RemoteOptions{
				Insecure: true,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := testdocker.New()
			require.NoError(t, err)

			// 1. Load a test image from the tar file, tag it and push to the test registry.
			err = d.ReplicateImage(ctx, tc.imageName, tc.imageFile, config)
			require.NoError(t, err)

			if tc.login {
				err = d.Login(config)
				require.NoError(t, err)

				defer d.Logout(config)
			}

			// 2. Analyze it
			imageRef := fmt.Sprintf("%s/%s", registryURL.Host, tc.imageName)
			imageDetail, err := analyze(ctx, imageRef, tc.option)
			require.Equal(t, tc.wantErr, err != nil, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.expectedOS, imageDetail.OS)
			assert.Equal(t, &tc.expectedRepo, imageDetail.Repository)
		})
	}
}

func getRegistryURL(ctx context.Context, registryC testcontainers.Container, exposedPort nat.Port) (*url.URL, error) {
	ip, err := registryC.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := registryC.MappedPort(ctx, exposedPort)
	if err != nil {
		return nil, err
	}

	urlStr := fmt.Sprintf("https://%s:%s", ip, port.Port())
	return url.Parse(urlStr)
}

func analyze(ctx context.Context, imageRef string, opt ftypes.RemoteOptions) (*ftypes.ArtifactDetail, error) {
	d, err := ioutil.TempDir("", "TestRegistry-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(d)

	c, err := cache.NewFSCache(d)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	runtimes := image.WithRuntimes(types.AllRuntimes)
	img, cleanup, err := image.NewContainerImage(ctx, imageRef, opt, runtimes)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	ar, err := aimage.NewArtifact(img, c, artifact.Option{
		DisabledAnalyzers: []analyzer.Type{
			analyzer.TypeExecutable,
			analyzer.TypeLicenseFile,
		},
	})
	if err != nil {
		return nil, err
	}

	ap := applier.NewApplier(c)

	imageInfo, err := ar.Inspect(ctx)
	if err != nil {
		return nil, err
	}

	imageDetail, err := ap.ApplyLayers(imageInfo.ID, imageInfo.BlobIDs)
	if err != nil {
		return nil, err
	}
	return &imageDetail, nil
}
