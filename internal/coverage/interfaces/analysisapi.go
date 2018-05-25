package interfaces

import "github.com/eltorocorp/drygopher/internal/pckg"

// AnalysisAPI represents anything that knows how to deal with coverage statistics.
type AnalysisAPI interface {
	GetCoverageStatistics(packages []string) (testedPackageStats, untestedPackageStats pckg.Group, err error)
}
