package installed

import (
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_composerInstalledAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      *analyzer.AnalysisResult
		wantErr   string
	}{
		{
			name:      "happy path",
			inputFile: "testdata/happy/installed.json",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.ComposerInstalled,
						FilePath: "testdata/happy/installed.json",
						Libraries: []types.Package{
							{
								ID:       "pear/log@1.13.3",
								Name:     "pear/log",
								Version:  "1.13.3",
								Indirect: false,
								Licenses: []string{"MIT"},
								Locations: []types.Location{
									{
										StartLine: 3,
										EndLine:   65,
									},
								},
								DependsOn: []string{"pear/pear_exception@v1.0.2"},
							},
							{
								ID:       "pear/pear_exception@v1.0.2",
								Name:     "pear/pear_exception",
								Version:  "v1.0.2",
								Indirect: false,
								Licenses: []string{"BSD-2-Clause"},
								Locations: []types.Location{
									{
										StartLine: 66,
										EndLine:   127,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "sad path",
			inputFile: "testdata/sad/installed.json",
			wantErr:   "decode error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer func() {
				err = f.Close()
				assert.NoError(t, err)
			}()

			a := composerInstalledAnalyzer{}
			got, err := a.Analyze(nil, analyzer.AnalysisInput{
				FilePath: tt.inputFile,
				Content:  f,
			})

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_composerInstalledAnalyzer_Required(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{
			name:     "happy path",
			filePath: "app/vendor/composer/installed.json",
			want:     true,
		},
		{
			name:     "sad path",
			filePath: "composer.json",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := composerInstalledAnalyzer{}
			got := a.Required(tt.filePath, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
