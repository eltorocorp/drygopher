//Package analysistypes contains types used by the analysis package.
package analysistypes

import "github.com/eltorocorp/drygopher/drygopher/coverage/pckg"

// GetCoverageStatisticsOutput contains results from the GetCoverageStatistics method.
type GetCoverageStatisticsOutput struct {
	TestedPackageStats      pckg.Group
	UntestedPackageStats    pckg.Group
	TestFailuresEncountered bool
}
