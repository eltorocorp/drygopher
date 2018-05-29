// Package coverageerrors contains error types specific to coverage analysis.
package coverageerrors

import "fmt"

// CoverageBelowStandard is an error returned when unit test coverage is below
// a set standard.
type CoverageBelowStandard struct {
	actualPercentage   float64
	standardPercentage float64
}

// NewCoverageBelowStandardError returns a CoverageBelowStandard error.
func NewCoverageBelowStandardError(standardPercentage, actualPercentage float64) CoverageBelowStandard {
	return CoverageBelowStandard{
		actualPercentage:   actualPercentage,
		standardPercentage: standardPercentage,
	}
}

func (e CoverageBelowStandard) Error() string {
	return fmt.Sprintf("coverage of %.2f%% is below the standard of %.2f%%", e.actualPercentage, e.standardPercentage)
}

// UnitTestFailed represents an error that results from one or more unit tests
// having failed during coverage analysis.
type UnitTestFailed struct{}

func (e UnitTestFailed) Error() string {
	return "One or more unit tests failed while conducting coverage analysis."
}

// NewUnitTestFailedError return a UnitTestFailed error.
func NewUnitTestFailedError() UnitTestFailed {
	return UnitTestFailed{}
}
