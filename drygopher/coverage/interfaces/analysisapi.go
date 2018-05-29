package interfaces

import (
	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis/analysistypes"
)

// AnalysisAPI represents anything that knows how to deal with coverage statistics.
type AnalysisAPI interface {
	GetCoverageStatistics(packages []string) (results analysistypes.GetCoverageStatisticsOutput, err error)
}
