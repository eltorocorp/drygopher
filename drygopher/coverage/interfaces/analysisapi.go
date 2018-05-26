package interfaces

import "github.com/eltorocorp/drygopher/drygopher/coverage/pckg"

// AnalysisAPI represents anything that knows how to deal with coverage statistics.
type AnalysisAPI interface {
	GetCoverageStatistics(packages []string) (testedPackageStats, untestedPackageStats pckg.Group, err error)
}
