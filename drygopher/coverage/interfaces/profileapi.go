package interfaces

import "github.com/eltorocorp/drygopher/drygopher/coverage/pckg"

// ProfileAPI represents anything that knows how to work with coverage profiles.
type ProfileAPI interface {
	BuildAndSaveCoverageProfile(allPackages pckg.Group, coverageProfileName string) error
}
