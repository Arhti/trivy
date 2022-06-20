package cyclonedx

import (
	"encoding/json"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/aquasecurity/trivy/pkg/purl"
	"github.com/samber/lo"

	"golang.org/x/xerrors"

	ftypes "github.com/aquasecurity/fanal/types"
	"github.com/aquasecurity/trivy/pkg/sbom"
)

const (
	FormatCycloneDX sbom.SBOMFormat = "cyclonedx"
)

type CycloneDX struct {
	extension string
	cdx.BOM
}

func New(name string) sbom.Parser {
	c := cdx.NewBOM()

	return &CycloneDX{
		BOM:       *c,
		extension: filepath.Ext(name),
	}
}

func (c CycloneDX) Parse(r io.Reader) (string, *ftypes.OS, []ftypes.PackageInfo, []ftypes.Application, error) {
	switch filepath.Ext(c.extension) {
	case ".json":
		if err := json.NewDecoder(r).Decode(&c); err != nil {
			return "", nil, nil, nil, xerrors.Errorf("failed to json decode: %w", err)
		}
	case ".xml":
		// TODO: not supported yet
	}
	return c.parse()
}

func (c CycloneDX) Type() sbom.SBOMFormat {
	return FormatCycloneDX
}

func parseOSPkgs(component cdx.Component, componentMap map[string]cdx.Component,
	depMap map[string][]string, seen map[string]struct{}) (ftypes.PackageInfo, error) {
	components := walkDependencies(component.BOMRef, nil, componentMap, depMap)
	pkgs, err := parsePkgs(components, seen)
	if err != nil {
		return ftypes.PackageInfo{}, xerrors.Errorf("failed to parse os package: %w", err)
	}

	return ftypes.PackageInfo{
		Packages: pkgs,
	}, nil
}

func parseLangPkgs(component cdx.Component, componentMap map[string]cdx.Component, depMap map[string][]string, seen map[string]struct{}) (*ftypes.Application, error) {
	components := walkDependencies(component.BOMRef, nil, componentMap, depMap)
	components = lo.UniqBy(components, func(c cdx.Component) string {
		return c.BOMRef
	})

	app := toApplication(component)
	pkgs, err := parsePkgs(components, seen)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse language-specific packages: %w", err)
	}
	app.Libraries = pkgs

	return app, nil
}

func parsePkgs(components []cdx.Component, seen map[string]struct{}) ([]ftypes.Package, error) {
	var pkgs []ftypes.Package
	for _, com := range components {
		seen[com.BOMRef] = struct{}{}
		_, pkg, err := toPackage(com)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse language package: %w", err)
		}
		pkgs = append(pkgs, *pkg)
	}
	return pkgs, nil
}

func (c CycloneDX) parse() (string, *ftypes.OS, []ftypes.PackageInfo, []ftypes.Application, error) {
	depMap := c.dependenciesMap()
	componentMap := c.componentMap()

	var (
		OS       *ftypes.OS
		apps     []ftypes.Application
		pkgInfos []ftypes.PackageInfo
	)
	seen := make(map[string]struct{})
	for bomRef := range depMap {
		component := componentMap[bomRef]
		switch component.Type {
		case cdx.ComponentTypeOS: // OS info and OS packages
			OS = toOS(component)
			pkgInfo, err := parseOSPkgs(component, componentMap, depMap, seen)
			if err != nil {
				return "", nil, nil, nil, xerrors.Errorf("failed to parse os packages: %w", err)
			}
			pkgInfos = append(pkgInfos, pkgInfo)
		case cdx.ComponentTypeApplication: // It would be a lock file in a CycloneDX report generated by Trivy
			if getProperty(component.Properties, PropertyType) == "" {
				continue
			}
			app, err := parseLangPkgs(component, componentMap, depMap, seen)
			if err != nil {
				return "", nil, nil, nil, xerrors.Errorf("failed to parse language packages: %w", err)
			}
			apps = append(apps, *app)
		case cdx.ComponentTypeLibrary:
			// It is an individual package not associated with any lock files and should be processed later.
			// e.g. .gemspec, .egg and .wheel
			continue
		}
	}

	for bomRef := range seen {
		delete(componentMap, bomRef)
	}

	libComponents := lo.Filter(lo.Values(componentMap), func(c cdx.Component, _ int) bool {
		return c.Type == cdx.ComponentTypeLibrary
	})

	aggregatedApps, err := c.Aggregate(libComponents)
	if err != nil {
		return "", nil, nil, nil, xerrors.Errorf("failed to aggregate packages: %w", err)
	}
	apps = append(apps, aggregatedApps...)

	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Type != apps[j].Type {
			return apps[i].Type < apps[j].Type
		}
		return apps[i].FilePath < apps[j].FilePath
	})

	return c.SerialNumber, OS, pkgInfos, apps, nil
}

// walkeDependencies takes all nested dependencies of the root component.
//
// e.g. Library A, B, C, D and E will be returned as dependencies of Application 1.
// type: Application 1
//   - type: Library A
//     - type: Library B
//   - type: Application 2
//     - type: Library C
//     - type: Application 3
//       - type: Library D
//       - type: Library E
func walkDependencies(rootRef string, components []cdx.Component, componentMap map[string]cdx.Component, dependenciesMap map[string][]string) []cdx.Component {
	deps, ok := dependenciesMap[rootRef]
	if !ok {
		return nil
	}
	for _, dep := range deps {
		component, ok := componentMap[dep]
		if !ok {
			continue
		}

		// Take only 'Libraries'
		if component.Type == cdx.ComponentTypeLibrary {
			components = append(components, component)
		}

		components = append(components, walkDependencies(dep, components, componentMap, dependenciesMap)...)
	}
	return components
}

func (c CycloneDX) componentMap() map[string]cdx.Component {
	componentMap := make(map[string]cdx.Component)

	if c.Components == nil {
		return componentMap
	}

	for _, component := range *c.Components {
		componentMap[component.BOMRef] = component
	}
	if c.Metadata != nil {
		componentMap[c.Metadata.Component.BOMRef] = *c.Metadata.Component
	}
	return componentMap
}

func (c CycloneDX) dependenciesMap() map[string][]string {
	dependencyMap := make(map[string][]string)

	if c.Dependencies == nil {
		return dependencyMap
	}
	for _, dep := range *c.Dependencies {
		if _, ok := dependencyMap[dep.Ref]; ok {
			continue
		}

		var refs []string
		if dep.Dependencies != nil {
			for _, d := range *dep.Dependencies {
				refs = append(refs, d.Ref)
			}
		}

		dependencyMap[dep.Ref] = refs
	}
	return dependencyMap
}

func (c CycloneDX) Aggregate(libs []cdx.Component) ([]ftypes.Application, error) {
	appMap := map[string]*ftypes.Application{}
	for _, lib := range libs {
		appType, pkg, err := toPackage(lib)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse purl to package: %w", err)
		}

		app, ok := appMap[appType]
		if !ok {
			// Pseudo application like Node.js and Ruby, not lock files
			app = &ftypes.Application{
				Type: appType,
			}
			appMap[appType] = app
		}

		app.Libraries = append(app.Libraries, *pkg)
	}

	var apps []ftypes.Application
	for _, app := range appMap {
		apps = append(apps, *app)
	}
	return apps, nil
}

func toOS(component cdx.Component) *ftypes.OS {
	return &ftypes.OS{
		Family: component.Name,
		Name:   component.Version,
	}
}

func toApplication(component cdx.Component) *ftypes.Application {
	return &ftypes.Application{
		Type:     getProperty(component.Properties, PropertyType),
		FilePath: component.Name,
	}
}

func toPackage(component cdx.Component) (string, *ftypes.Package, error) {
	p, err := purl.FromString(component.PackageURL)
	if err != nil {
		return "", nil, xerrors.Errorf("failed to parse purl: %w", err)
	}

	pkg := p.Package()
	pkg.Ref = component.BOMRef

	var licenses []string
	if component.Licenses != nil {
		for _, license := range *component.Licenses {
			licenses = append(licenses, license.Expression)
		}
	}
	// TODO: In Trivy's SBOM, Expression is singular
	pkg.License = strings.Join(licenses, ", ")

	if component.Properties == nil {
		return p.AppType(), pkg, nil
	}

	for _, p := range *component.Properties {
		if strings.HasPrefix(p.Name, Namespace) {
			switch strings.TrimPrefix(p.Name, Namespace) {
			case PropertySrcName:
				pkg.SrcName = p.Value
			case PropertySrcVersion:
				pkg.SrcVersion = p.Value
			case PropertySrcRelease:
				pkg.SrcRelease = p.Value
			case PropertySrcEpoch:
				pkg.SrcEpoch, err = strconv.Atoi(p.Value)
				if err != nil {
					return "", nil, xerrors.Errorf("failed to parse source epoch: %w", err)
				}
			case PropertyModularitylabel:
				pkg.Modularitylabel = p.Value
			case PropertyLayerDiffID:
				pkg.Layer.DiffID = p.Value
			}
		}
	}

	return p.AppType(), pkg, nil
}

func getProperty(properties *[]cdx.Property, key string) string {
	if properties == nil {
		return ""
	}

	for _, p := range *properties {
		if p.Name == Namespace+key {
			return p.Value
		}
	}
	return ""
}
