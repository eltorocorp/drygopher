package pckg

import (
	"sort"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// Group is a slice of package stats.
type Group []*Stats

// SetEstimatedStmtCountFrom sets the estimated statement count based on the
// average number of statements in other packages in the group.
func (g *Group) SetEstimatedStmtCountFrom(reference Group) {
	for _, p := range *g {
		p.Statements = reference.MedianStatementCount()
	}
}

// MedianStatementCount calculates the median number of statements in packages
// that currently have associated tests.
func (g *Group) MedianStatementCount() float64 {
	sc := g.StatementCounts()
	if floats.Sum(sc) == 0 {
		return 0
	}
	sort.Float64s(sc)
	return stat.Quantile(0.5, stat.Empirical, sc, nil)
}

// TotalStatementCount returns the sum of statements across all packaes that
// currently have associated tests.
func (g *Group) TotalStatementCount() float64 {
	return floats.Sum(g.StatementCounts())
}

// StatementCounts returns a list of statement counts from the group.
func (g *Group) StatementCounts() (statementCounts []float64) {
	for _, p := range *g {
		statementCounts = append(statementCounts, p.Statements)
	}
	return
}

// TotalCovered returns the total number of statements that have been covered
// by unit tests.
func (g *Group) TotalCovered() float64 {
	c := []float64{}
	for _, p := range *g {
		c = append(c, p.Covered)
	}
	return floats.Sum(c)
}

// TotalUncovered returns the total number of statements that have not been
// covered by unit tests.
func (g *Group) TotalUncovered() float64 {
	u := []float64{}
	for _, p := range *g {
		u = append(u, p.Uncovered)
	}
	return floats.Sum(u)
}

// CoveragePercent calculates the coverage percentage across all packages.
func (g *Group) CoveragePercent() (coveragePercent float64) {
	if g.TotalStatementCount() >= 1 {
		coveragePercent = g.TotalCovered() / g.TotalStatementCount()
	}
	return
}

// EstimateCount returns the number of packages who have had their coverage
// percentages estimated.
func (g *Group) EstimateCount() (estimatedCount float64) {
	for _, p := range *g {
		if p.Estimated {
			estimatedCount++
		}
	}
	return
}
