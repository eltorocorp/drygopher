package pckg

// Stats represent statistics for one or more packages.
type Stats struct {
	Package         string
	Statements      float64
	Covered         float64
	Uncovered       float64
	Estimated       bool
	RawCoverageData []string
}

// CoveragePercent calculates the coverage percentage for the package stats.
func (s *Stats) CoveragePercent() (coveragePercent float64) {
	if s.Statements >= 1 {
		coveragePercent = s.Covered / s.Statements
	}
	return
}
