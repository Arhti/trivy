package gradle

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/go-dep-parser/pkg/gradle/lockfile"
	godeptypes "github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer/language"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/utils/fsutils"
)

func init() {
	analyzer.RegisterPostAnalyzer(analyzer.TypeGradleLock, newGradleLockAnalyzer)
}

const (
	version        = 2
	fileNameSuffix = "gradle.lockfile"
)

// gradleLockAnalyzer analyzes '*gradle.lockfile'
type gradleLockAnalyzer struct {
	parser godeptypes.Parser
}

func newGradleLockAnalyzer(_ analyzer.AnalyzerOptions) (analyzer.PostAnalyzer, error) {
	return &gradleLockAnalyzer{
		parser: lockfile.NewParser(),
	}, nil
}

func (a gradleLockAnalyzer) PostAnalyze(_ context.Context, input analyzer.PostAnalysisInput) (*analyzer.AnalysisResult, error) {
	poms, err := parsePoms()
	if err != nil {
		log.Logger.Warnf("Unable to get licenses: %s", err)
	}

	required := func(path string, d fs.DirEntry) bool {
		return a.Required(path, nil)
	}

	var apps []types.Application
	err = fsutils.WalkDir(input.FS, ".", required, func(path string, _ fs.DirEntry, r io.Reader) error {
		var app *types.Application
		app, err = language.Parse(types.Gradle, path, r, a.parser)
		if err != nil {
			return xerrors.Errorf("unable to parse %q: %w", path, err)
		}

		if app == nil {
			return nil
		}

		libs := lo.SliceToMap(app.Libraries, func(lib types.Package) (string, types.Package) {
			return lib.ID, lib
		})

		for i, lib := range app.Libraries {
			pom := poms[lib.ID]

			if len(pom.Licenses.License) > 0 {
				app.Libraries[i].Licenses = lo.Map(pom.Licenses.License, func(license License, _ int) string {
					return license.Name
				})
			}

			var deps []string
			for _, dep := range pom.Dependencies.Dependency {
				id := packageID(dep.GroupID, dep.ArtifactID, dep.Version)
				if _, ok := libs[id]; ok {
					deps = append(deps, id)
				}
			}
			app.Libraries[i].DependsOn = deps
		}

		sort.Sort(app.Libraries)
		apps = append(apps, *app)
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("walk error: %w", err)
	}

	return &analyzer.AnalysisResult{
		Applications: apps,
	}, nil
}

func (a gradleLockAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return strings.HasSuffix(filePath, fileNameSuffix)
}

func (a gradleLockAnalyzer) Type() analyzer.Type {
	return analyzer.TypeGradleLock
}

func (a gradleLockAnalyzer) Version() int {
	return version
}

func packageID(groupId, artifactId, ver string) string {
	return fmt.Sprintf("%s:%s:%s", groupId, artifactId, ver)
}
