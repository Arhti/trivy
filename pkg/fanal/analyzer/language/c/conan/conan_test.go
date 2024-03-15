package conan

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
)

func Test_conanLockAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		cacheDir string
		want     *analyzer.AnalysisResult
	}{
		{
			name: "happy path",
			dir:  "testdata/happy",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Conan,
						FilePath: "conan.lock",
						Libraries: types.Packages{
							{
								ID:      "openssl/3.0.5",
								Name:    "openssl",
								Version: "3.0.5",
								DependsOn: []string{
									"zlib/1.2.12",
								},
								Locations: []types.Location{
									{
										StartLine: 12,
										EndLine:   21,
									},
								},
							},
							{
								ID:       "zlib/1.2.12",
								Name:     "zlib",
								Version:  "1.2.12",
								Indirect: true,
								Locations: []types.Location{
									{
										StartLine: 22,
										EndLine:   28,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "happy path with cache dir",
			dir:      "testdata/happy",
			cacheDir: "testdata/cacheDir",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Conan,
						FilePath: "conan.lock",
						Libraries: types.Packages{
							{
								ID:      "openssl/3.0.5",
								Name:    "openssl",
								Version: "3.0.5",
								Licenses: []string{
									"Apache-2.0",
								},
								DependsOn: []string{
									"zlib/1.2.12",
								},
								Locations: []types.Location{
									{
										StartLine: 12,
										EndLine:   21,
									},
								},
							},
							{
								ID:      "zlib/1.2.12",
								Name:    "zlib",
								Version: "1.2.12",
								Licenses: []string{
									"Zlib",
								},
								Indirect: true,
								Locations: []types.Location{
									{
										StartLine: 22,
										EndLine:   28,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "empty file",
			dir:  "testdata/empty",
			want: &analyzer.AnalysisResult{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cacheDir != "" {
				t.Setenv("CONAN_USER_HOME", tt.cacheDir)
			}
			a, err := newConanLockAnalyzer(analyzer.AnalyzerOptions{})
			require.NoError(t, err)

			got, err := a.PostAnalyze(context.Background(), analyzer.PostAnalysisInput{
				FS: os.DirFS(tt.dir),
			})

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_conanLockAnalyzer_Required(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{
			name:     "default name",
			filePath: "conan.lock",
			want:     true,
		},
		{
			name:     "name with prefix",
			filePath: "pkga_deps.lock",
			want:     false,
		},
		{
			name:     "txt",
			filePath: "test.txt",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			f, err := os.Create(filepath.Join(dir, tt.filePath))
			require.NoError(t, err)
			defer f.Close()

			fi, err := f.Stat()
			require.NoError(t, err)

			a := conanLockAnalyzer{}
			got := a.Required("", fi)
			assert.Equal(t, tt.want, got)
		})
	}
}
