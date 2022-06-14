//go:build integration
// +build integration

package integration

import (
	"github.com/aquasecurity/trivy/pkg/commands"
	"github.com/stretchr/testify/assert"
	"io"
	"path/filepath"
	"testing"
)

func TestCycloneDX(t *testing.T) {

	type args struct {
		input        string
		artifactType string
	}
	tests := []struct {
		name   string
		args   args
		golden string
	}{
		{
			name: "centos7-bom by trivy",
			args: args{
				input:        "testdata/fixtures/sbom/centos7-bom.json",
				artifactType: "cyclonedx",
			},
			golden: "testdata/centos7-bom.json.golden",
		},
		{
			name: "fluentd-multiple-lockfiles-bom by trivy",
			args: args{
				input:        "testdata/fixtures/sbom/fluentd-multiple-lockfiles-bom.json",
				artifactType: "cyclonedx",
			},
			golden: "testdata/fluentd-multiple-lockfiles-bom.json.golden",
		},
	}

	// Set up testing DB
	cacheDir := initDB(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osArgs := []string{
				"trivy", "--cache-dir", cacheDir, "sbom", "--skip-db-update",
				"--artifact-type", tt.args.artifactType,
			}

			// Setup the output file
			outputFile := filepath.Join(t.TempDir(), "output.json")
			if *update {
				outputFile = tt.golden
			}

			osArgs = append(osArgs, "--output", outputFile)
			osArgs = append(osArgs, tt.args.input)

			// Setup CLI App
			app := commands.NewApp("dev")
			app.Writer = io.Discard

			// Run "trivy fs"
			assert.Nil(t, app.Run(osArgs))

			// Compare want and got
			compareReports(t, tt.golden, outputFile)
		})
	}
}
