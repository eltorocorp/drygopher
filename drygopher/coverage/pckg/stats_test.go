package pckg_test

import (
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/stretchr/testify/assert"
)

func Test_StatsCoveragePercent_Normally_ReturnsCoverage(t *testing.T) {
	stats := pckg.Stats{
		Covered:    1.0,
		Statements: 3.0,
	}

	assert.InEpsilon(t, 0.3333, stats.CoveragePercent(), 0.001)
}
