package cyclonedx

import (
	"strconv"
	"strings"

	"github.com/aquasecurity/trivy/pkg/sbom/cyclonedx/core"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/aquasecurity/trivy/pkg/digest"
	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/purl"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/samber/lo"
	"golang.org/x/xerrors"
)

const (
	PropertySchemaVersion = "SchemaVersion"
	PropertyType          = "Type"
	PropertyClass         = "Class"

	// Image properties
	PropertySize       = "Size"
	PropertyImageID    = "ImageID"
	PropertyRepoDigest = "RepoDigest"
	PropertyDiffID     = "DiffID"
	PropertyRepoTag    = "RepoTag"

	// Package properties
	PropertyPkgID           = "PkgID"
	PropertyPkgType         = "PkgType"
	PropertySrcName         = "SrcName"
	PropertySrcVersion      = "SrcVersion"
	PropertySrcRelease      = "SrcRelease"
	PropertySrcEpoch        = "SrcEpoch"
	PropertyModularitylabel = "Modularitylabel"
	PropertyFilePath        = "FilePath"
	PropertyLayerDigest     = "LayerDigest"
	PropertyLayerDiffID     = "LayerDiffID"

	// https://json-schema.org/understanding-json-schema/reference/string.html#dates-and-times
	timeLayout = "2006-01-02T15:04:05+00:00"
)

var (
	ErrInvalidBOMLink = xerrors.New("invalid bomLink format error")
)

type Marshaler struct {
	core *core.CycloneDX
}

func NewMarshaler(version string, opts ...core.Option) *Marshaler {
	return &Marshaler{
		core: core.NewCycloneDX(version, opts...),
	}
}

// Marshal converts the Trivy report to the CycloneDX format
func (e *Marshaler) Marshal(report types.Report) (*cdx.BOM, error) {
	// Convert
	root, err := e.marshalReport(report)
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal report: %w", err)
	}

	return e.core.Marshal(root), nil
}

// MarshalVulnerabilities converts the Trivy report to the CycloneDX format only with vulnerabilities.
// The output refers to another CycloneDX SBOM.
func (e *Marshaler) MarshalVulnerabilities(report types.Report) (*cdx.BOM, error) {
	return nil, nil
	//vulnMap := map[string]cdx.Vulnerability{}
	//for _, result := range report.Results {
	//	for _, vuln := range result.Vulnerabilities {
	//		ref, err := externalRef(report.CycloneDX.SerialNumber, vuln.PkgRef)
	//		if err != nil {
	//			return nil, err
	//		}
	//		if v, ok := vulnMap[vuln.VulnerabilityID]; ok {
	//			*v.Affects = append(*v.Affects, cdxAffects(ref, vuln.InstalledVersion))
	//		} else {
	//			vulnMap[vuln.VulnerabilityID] = toCdxVulnerability(ref, vuln)
	//		}
	//	}
	//}
	//vulns := maps.Values(vulnMap)
	//sort.Slice(vulns, func(i, j int) bool {
	//	return vulns[i].ID > vulns[j].ID
	//})
	//
	//bom := cdx.NewBOM()
	//bom.Metadata = e.cdxMetadata()
	//
	//// Fill the detected vulnerabilities
	//bom.Vulnerabilities = &vulns
	//
	//// Use the original component as is
	//bom.Metadata.Component = &cdx.Component{
	//	Name:    report.CycloneDX.Metadata.Component.Name,
	//	Version: report.CycloneDX.Metadata.Component.Version,
	//	Type:    cdx.ComponentType(report.CycloneDX.Metadata.Component.Type),
	//}
	//
	//// Overwrite the bom ref as it must be the BOM ref of the original CycloneDX.
	//// e.g.
	////  "metadata" : {
	////    "timestamp" : "2022-07-02T00:00:00Z",
	////    "component" : {
	////      "name" : "Acme Product",
	////      "version": "2.4.0",
	////      "type" : "application",
	////      "bom-ref" : "urn:cdx:f08a6ccd-4dce-4759-bd84-c626675d60a7/1"
	////    }
	////  },
	//if report.CycloneDX.SerialNumber != "" { // bomRef is optional field - https://cyclonedx.org/docs/1.4/json/#metadata_component_bom-ref
	//	bom.Metadata.Component.BOMRef = fmt.Sprintf("%s/%d", report.CycloneDX.SerialNumber, report.CycloneDX.Version)
	//}
	//return bom, nil
}

func externalRef(bomLink string, bomRef string) (string, error) {
	// bomLink is optional field: https://cyclonedx.org/docs/1.4/json/#vulnerabilities_items_bom-ref
	if bomLink == "" {
		return bomRef, nil
	}
	if !strings.HasPrefix(bomLink, "urn:uuid:") {
		return "", xerrors.Errorf("%q: %w", bomLink, ErrInvalidBOMLink)
	}
	return fmt.Sprintf("%s/%d#%s", strings.Replace(bomLink, "uuid", "cdx", 1), cdx.BOMFileFormatJSON, bomRef), nil
}

func (e *Marshaler) marshalResult(metadata types.Metadata, result types.Result) ([]*core.Component, error) {
	if result.Type == ftypes.NodePkg || result.Type == ftypes.PythonPkg ||
		result.Type == ftypes.GemSpec || result.Type == ftypes.Jar || result.Type == ftypes.CondaPkg {
		// If a package is language-specific package that isn't associated with a lock file,
		// it will be a dependency of a component under "metadata".
		// e.g.
		//   Container component (alpine:3.15) ----------------------- #1
		//     -> Library component (npm package, express-4.17.3) ---- #2
		//     -> Library component (python package, django-4.0.2) --- #2
		//     -> etc.
		// ref. https://cyclonedx.org/use-cases/#inventory

		// Dependency graph from #1 to #2
		components, err := e.marshalPackages(metadata, result)
		if err != nil {
			return nil, err
		}
		return components, nil
	} else if result.Class == types.ClassOSPkg || result.Class == types.ClassLangPkg {
		// If a package is OS package, it will be a dependency of "Operating System" component.
		// e.g.
		//   Container component (alpine:3.15) --------------------- #1
		//     -> Operating System Component (Alpine Linux 3.15) --- #2
		//       -> Library component (bash-4.12) ------------------ #3
		//       -> Library component (vim-8.2)   ------------------ #3
		//       -> etc.
		//
		// Else if a package is language-specific package associated with a lock file,
		// it will be a dependency of "Application" component.
		// e.g.
		//   Container component (alpine:3.15) ------------------------ #1
		//     -> Application component (/app/package-lock.json) ------ #2
		//       -> Library component (npm package, express-4.17.3) --- #3
		//       -> Library component (npm package, lodash-4.17.21) --- #3
		//       -> etc.

		// #2
		appComponent := e.resultComponent(result, metadata.OS)

		// #3
		components, err := e.marshalPackages(metadata, result)
		if err != nil {
			return nil, err
		}

		// Dependency graph from #2 to #3
		appComponent.Components = components

		// Dependency graph from #1 to #2
		return []*core.Component{appComponent}, nil
	}
	return nil, nil
}

func (e *Marshaler) marshalPackages(metadata types.Metadata, result types.Result) ([]*core.Component, error) {
	// Get dependency parents first
	parents := ftypes.Packages(result.Packages).ParentDeps()

	// Group vulnerabilities by package ID
	vulns := lo.GroupBy(result.Vulnerabilities, func(v types.DetectedVulnerability) string {
		return v.PkgID
	})

	// Create package map
	pkgs := lo.SliceToMap(result.Packages, func(pkg ftypes.Package) (string, Package) {
		return pkg.ID, Package{
			Type:            result.Type,
			Metadata:        metadata,
			Package:         pkg,
			Vulnerabilities: vulns[pkg.ID],
		}
	})

	var directComponents []*core.Component
	for _, pkg := range pkgs {
		// Skip indirect dependencies
		if pkg.Indirect && len(parents[pkg.ID]) != 0 {
			continue
		}

		// Recursive packages from direct dependencies
		component, err := e.marshalPackage(pkg, pkgs, map[string]*core.Component{})
		if err != nil {
			return nil, nil
		}
		directComponents = append(directComponents, component)
	}

	return directComponents, nil
}

type Package struct {
	ftypes.Package
	Type            string
	Metadata        types.Metadata
	Vulnerabilities []types.DetectedVulnerability
}

func (e *Marshaler) marshalPackage(pkg Package, pkgs map[string]Package, components map[string]*core.Component,
) (*core.Component, error) {
	if _, ok := components[pkg.ID]; ok {
		return nil, nil
	}

	component, err := pkgComponent(pkg)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse pkg: %w", err)
	}
	components[pkg.ID] = component

	// Iterate dependencies
	for _, dep := range pkg.DependsOn {
		childPkg, ok := pkgs[dep]
		if !ok {
			continue
		}

		child, err := e.marshalPackage(childPkg, pkgs, components)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse pkg: %w", err)
		}
		component.Components = append(component.Components, child)
	}
	return component, nil
}

func (e *Marshaler) marshalReport(r types.Report) (*core.Component, error) {
	// Metadata component
	root, err := e.rootComponent(r)
	if err != nil {
		return nil, err
	}

	for _, result := range r.Results {
		components, err := e.marshalResult(r.Metadata, result)
		if err != nil {
			return nil, err
		}
		root.Components = append(root.Components, components...)
	}
	return root, nil
}

func (e *Marshaler) rootComponent(r types.Report) (*core.Component, error) {
	root := &core.Component{
		Name: r.ArtifactName,
	}

	props := map[string]string{
		PropertySchemaVersion: strconv.Itoa(r.SchemaVersion),
	}

	switch r.ArtifactType {
	case ftypes.ArtifactContainerImage:
		root.Type = cdx.ComponentTypeContainer
		props[PropertyImageID] = r.Metadata.ImageID

		p, err := purl.NewPackageURL(purl.TypeOCI, r.Metadata, ftypes.Package{})
		if err != nil {
			return nil, xerrors.Errorf("failed to new package url for oci: %w", err)
		} else if p.Type != "" {
			root.PackageURL = &p
		}

	case ftypes.ArtifactVM:
		root.Type = cdx.ComponentTypeContainer
	case ftypes.ArtifactFilesystem, ftypes.ArtifactRemoteRepository:
		root.Type = cdx.ComponentTypeApplication
	}

	if r.Metadata.Size != 0 {
		props[PropertySize] = strconv.FormatInt(r.Metadata.Size, 10)
	}

	if len(r.Metadata.RepoDigests) > 0 {
		props[PropertyRepoDigest] = strings.Join(r.Metadata.RepoDigests, ",")
	}
	if len(r.Metadata.DiffIDs) > 0 {
		props[PropertyDiffID] = strings.Join(r.Metadata.DiffIDs, ",")
	}
	if len(r.Metadata.RepoTags) > 0 {
		props[PropertyRepoTag] = strings.Join(r.Metadata.RepoTags, ",")
	}

	root.Properties = filterProperties(props)

	return root, nil
}

func (e *Marshaler) resultComponent(r types.Result, osFound *ftypes.OS) *core.Component {
	component := &core.Component{
		Name: r.Target,
		Properties: map[string]string{
			PropertyType:  r.Type,
			PropertyClass: string(r.Class),
		},
	}

	switch r.Class {
	case types.ClassOSPkg:
		// UUID needs to be generated since Operating System Component cannot generate PURL.
		// https://cyclonedx.org/use-cases/#known-vulnerabilities
		if osFound != nil {
			component.Name = osFound.Family
			component.Version = osFound.Name
		}
		component.Type = cdx.ComponentTypeOS
	case types.ClassLangPkg:
		// UUID needs to be generated since Application Component cannot generate PURL.
		// https://cyclonedx.org/use-cases/#known-vulnerabilities
		component.Type = cdx.ComponentTypeApplication
	}

	return component
}

func pkgComponent(pkg Package) (*core.Component, error) {
	pu, err := purl.NewPackageURL(pkg.Type, pkg.Metadata, pkg.Package)
	if err != nil {
		return nil, xerrors.Errorf("failed to new package purl: %w", err)
	}

	properties := map[string]string{
		PropertyPkgID:           pkg.ID,
		PropertyPkgType:         pkg.Type,
		PropertyFilePath:        pkg.FilePath,
		PropertySrcName:         pkg.SrcName,
		PropertySrcVersion:      pkg.SrcVersion,
		PropertySrcRelease:      pkg.SrcRelease,
		PropertySrcEpoch:        strconv.Itoa(pkg.SrcEpoch),
		PropertyModularitylabel: pkg.Modularitylabel,
		PropertyLayerDigest:     pkg.Layer.Digest,
		PropertyLayerDiffID:     pkg.Layer.DiffID,
	}

	return &core.Component{
		Type:            cdx.ComponentTypeLibrary,
		Name:            pkg.Name,
		Version:         pu.Version,
		PackageURL:      &pu,
		Supplier:        pkg.Maintainer,
		Licenses:        pkg.Licenses,
		Hashes:          lo.Ternary(pkg.Digest == "", nil, []digest.Digest{pkg.Digest}),
		Properties:      filterProperties(properties),
		Vulnerabilities: pkg.Vulnerabilities,
	}, nil
}

func filterProperties(props map[string]string) map[string]string {
	return lo.OmitBy(props, func(key string, value string) bool {
		return value == "" || (key == PropertySrcEpoch && value == "0")
	})
}
