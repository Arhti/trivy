package node

import (
	"golang.org/x/xerrors"

	npm "github.com/aquasecurity/go-npm-version/pkg"
	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/detector/library/comparer"
)

// NpmComparer represents a comparer for npm
type NpmComparer struct{}

// IsVulnerable checks if the package version is vulnerable to the advisory.
func (n NpmComparer) IsVulnerable(ver string, advisory dbTypes.Advisory) bool {
	return comparer.IsVulnerable(ver, advisory, n.MatchVersion)
}

// MatchVersion checks if the package version satisfies the given constraint.
func (n NpmComparer) MatchVersion(currentVersion, constraint string) (bool, error) {
	v, err := npm.NewVersion(currentVersion)
	if err != nil {
		return false, xerrors.Errorf(" npm version error (%s): %s", currentVersion, err)
	}

	c, err := npm.NewConstraints(constraint)
	if err != nil {
		return false, xerrors.Errorf("npm constraint error (%s): %s", constraint, err)
	}

	return c.Check(v), nil
}
