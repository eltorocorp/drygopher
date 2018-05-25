package interfaces

import "github.com/eltorocorp/drygopher/internal/pckg"

// ProfileAPI represents anything that knows how to work with coverage profiles.
type ProfileAPI interface {
	SaveCoverageProfile(fileName string, rawData []string) error
	BuildAndSaveCoverageProfile(allPackages pckg.Group, coverageProfileName string) error
}
