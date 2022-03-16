package artifact

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/aquasecurity/fanal/cache"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/rpc/client"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
)

func initializeClientDockerScanner(ctx context.Context, target string, ac cache.ArtifactCache, lac cache.LocalArtifactCache,
	remoteAddr string, customHeaders http.Header,
	insecure bool, artifactOpt artifact.Option, configScannerOptions config.ScannerOption) (scanner.Scanner, func(), error) {
	// Scan an image in Docker Engine or Docker Registry
	dockerOpt, err := types.GetDockerOption(insecure)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}

	s, cleanup, err := initializeRemoteDockerScanner(ctx, target, ac, client.CustomHeaders(customHeaders),
		client.RemoteURL(remoteAddr), client.Insecure(insecure), dockerOpt, artifactOpt, configScannerOptions)
	if err != nil {
		return scanner.Scanner{}, nil, xerrors.Errorf("unable to initialize the docker scanner: %w", err)
	}
	return s, cleanup, nil
}

func initializeClientTarScanner(ctx context.Context, target string, ac cache.ArtifactCache, lac cache.LocalArtifactCache,
	remoteAddr string, customHeaders http.Header,
	insecure bool, artifactOpt artifact.Option, configScannerOptions config.ScannerOption) (scanner.Scanner, func(), error) {
	// Scan tar file
	s, err := initializeRemoteArchiveScanner(ctx, target, ac, client.CustomHeaders(customHeaders),
		client.RemoteURL(remoteAddr), client.Insecure(insecure), artifactOpt, configScannerOptions)
	if err != nil {
		return scanner.Scanner{}, nil, xerrors.Errorf("unable to initialize the archive scanner: %w", err)
	}
	return s, func() {}, nil
}

func ClientRun(cliCtx *cli.Context) error {
	opt, err := initOption(cliCtx)
	if err != nil {
		return xerrors.Errorf("option error: %w", err)
	}
	ctx, cancel := context.WithTimeout(cliCtx.Context, opt.Timeout)
	defer cancel()

	// Disable the lock file scanning
	opt.DisabledAnalyzers = analyzer.TypeLockfiles

	if opt.Input != "" {
		err = runWithTimeout(ctx, opt, initializeClientTarScanner, initCache)
	} else {
		err = runWithTimeout(ctx, opt, initializeClientDockerScanner, initCache)
	}
	if xerrors.Is(err, context.DeadlineExceeded) {
		log.Logger.Warn("Increase --timeout value")
	}
	return err
}
