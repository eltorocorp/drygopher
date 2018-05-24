package pckg

// Stats represent statistics for a package.
type Stats struct {
	Package    string
	Statements float64
	Covered    float64
	Uncovered  float64
	Estimated  bool
}

// CoveragePercent calculates the coverage percentage for the package stats.
func (ps *Stats) CoveragePercent() (coveragePercent float64) {
	if ps.Statements >= 1 {
		coveragePercent = ps.Covered / ps.Statements
	}
	return
}
