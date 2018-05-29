package coverageerrors_test

import (
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/coverageerrors"
	"github.com/stretchr/testify/assert"
)

func Test_CoverageBelowStandard_Error_Normally_ReturnsErrorMessage(t *testing.T) {
	err := coverageerrors.NewCoverageBelowStandardError(98.3, 50.5)
	assert.EqualError(t, err, "coverage of 50.50% is below the standard of 98.30%")
}
