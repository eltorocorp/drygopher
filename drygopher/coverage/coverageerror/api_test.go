package coverageerror_test

import (
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/coverageerror"
	"github.com/stretchr/testify/assert"
)

func Test_New_Normally_ReturnsNew(t *testing.T) {
	err := coverageerror.New(100, 50)
	assert.IsType(t, coverageerror.CoverageBelowStandard{}, err)
}

func Test_Error_Normally_ReturnsErrorMessage(t *testing.T) {
	err := coverageerror.New(98.3, 50.5)
	assert.EqualError(t, err, "coverage of 50.50% is below the standard of 98.30%")
}
