package swift

import (
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_swiftLockAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      *analyzer.AnalysisResult
	}{
		{
			name:      "happy path",
			inputFile: "testdata/happy/Package.resolved",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.Swift,
						FilePath: "testdata/happy/Package.resolved",
						Libraries: types.Packages{

							{
								Name:    "https://github.com/Quick/Nimble.git",
								Version: "9.2.1",
								Locations: []types.Location{
									{
										StartLine: 4,
										EndLine:   12,
									},
								},
							},
							{
								Name:    "https://github.com/Quick/Quick.git",
								Version: "7.0.0",
								Locations: []types.Location{
									{
										StartLine: 13,
										EndLine:   21,
									},
								},
							},
							{
								Name:    "https://github.com/ReactiveCocoa/ReactiveSwift",
								Version: "7.1.1",
								Locations: []types.Location{
									{
										StartLine: 22,
										EndLine:   30,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "empty file",
			inputFile: "testdata/empty/Package.resolved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer f.Close()

			a := swiftLockAnalyzer{}
			got, err := a.Analyze(nil, analyzer.AnalysisInput{
				FilePath: tt.inputFile,
				Content:  f,
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
