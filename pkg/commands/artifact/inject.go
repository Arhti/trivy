//go:build wireinject
// +build wireinject

package artifact

import (
	"context"

	"github.com/google/wire"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/aquasecurity/fanal/cache"
	"github.com/aquasecurity/fanal/types"
	"github.com/aquasecurity/trivy/pkg/result"
	"github.com/aquasecurity/trivy/pkg/rpc/client"
	"github.com/aquasecurity/trivy/pkg/scanner"
)

func initializeDockerScanner(ctx context.Context, imageName string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, dockerOpt types.DockerOption, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneDockerSet)
	return scanner.Scanner{}, nil, nil
}

func initializeArchiveScanner(ctx context.Context, filePath string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, error) {
	wire.Build(scanner.StandaloneArchiveSet)
	return scanner.Scanner{}, nil
}

func initializeFilesystemScanner(ctx context.Context, path string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneFilesystemSet)
	return scanner.Scanner{}, nil, nil
}

func initializeFilesystemRemoteScanner(ctx context.Context, path string, artifactCache cache.ArtifactCache,
	customHeaders client.CustomHeaders, url client.RemoteURL, insecure client.Insecure,
	artifactOption artifact.Option, configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneFilesystemRemoteSet)
	return scanner.Scanner{}, nil, nil
}
func initializeRepositoryScanner(ctx context.Context, url string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneRepositorySet)
	return scanner.Scanner{}, nil, nil
}

func initializeResultClient() result.Client {
	wire.Build(result.SuperSet)
	return result.Client{}
}

func initializeRemoteDockerScanner(ctx context.Context, imageName string, artifactCache cache.ArtifactCache, customHeaders client.CustomHeaders,
	url client.RemoteURL, insecure client.Insecure, dockerOpt types.DockerOption, artifactOption artifact.Option, configScannerOption config.ScannerOption) (
	scanner.Scanner, func(), error) {
	wire.Build(scanner.RemoteDockerSet)
	return scanner.Scanner{}, nil, nil
}

func initializeRemoteArchiveScanner(ctx context.Context, filePath string, artifactCache cache.ArtifactCache,
	customHeaders client.CustomHeaders, url client.RemoteURL, insecure client.Insecure, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, error) {
	wire.Build(scanner.RemoteArchiveSet)
	return scanner.Scanner{}, nil
}

func initializeRemoteResultClient() result.Client {
	wire.Build(result.SuperSet)
	return result.Client{}
}
