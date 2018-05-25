package coverageerror

import "fmt"

// CoverageBelowStandard is an error returned when unit test coverage is below
// a set standard.
type CoverageBelowStandard struct {
	actualPercentage   float64
	standardPercentage float64
}

// New returns a reference to an CoverageBelowStandard error.
func New(standardPercentage, actualPercentage float64) CoverageBelowStandard {
	return CoverageBelowStandard{
		actualPercentage:   actualPercentage,
		standardPercentage: standardPercentage,
	}
}

func (e CoverageBelowStandard) Error() string {
	return fmt.Sprintf("coverage of %v%% is below the standard of %v%%", e.actualPercentage, e.standardPercentage)
}
