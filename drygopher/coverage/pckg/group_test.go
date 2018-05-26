package pckg_test

import (
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/stretchr/testify/assert"
)

func Test_MedianStatementCount_Normally_ReturnsMedian(t *testing.T) {
	group := pckg.Group{
		&pckg.Stats{
			Statements: 1.0,
		},
		&pckg.Stats{
			Statements: 2.0,
		},
		&pckg.Stats{
			Statements: 3.0,
		},
	}
	assert.Equal(t, 2.0, group.MedianStatementCount())
}

func Test_MedianStatementCount_StatementCountIsZero_ReturnsZero(t *testing.T) {
	group := pckg.Group{
		new(pckg.Stats),
	}
	assert.Equal(t, 0.0, group.MedianStatementCount())

}

func Test_SetEstimatedStmtCountFrom_Normally_SetsEstimate(t *testing.T) {
	groupOne := pckg.Group{
		&pckg.Stats{
			Statements: 1.0,
		},
	}
	groupTwo := pckg.Group{
		&pckg.Stats{
			Estimated: true,
		},
	}

	groupTwo.SetEstimatedStmtCountFrom(groupOne)

	assert.Equal(t, 1.0, groupTwo.EstimateCount())
}

func Test_TotalStatementCount_Normally_ReturnsTotal(t *testing.T) {
	group := pckg.Group{
		&pckg.Stats{
			Statements: 1.0,
		},
		&pckg.Stats{
			Statements: 2.0,
		},
		&pckg.Stats{
			Statements: 3.0,
		},
	}
	assert.Equal(t, 6.0, group.TotalStatementCount())

}

func Test_TotalCovered_Normally_ReturnsTotal(t *testing.T) {
	group := pckg.Group{
		&pckg.Stats{
			Covered: 1.0,
		},
		&pckg.Stats{
			Covered: 2.0,
		},
		&pckg.Stats{
			Covered: 3.0,
		},
	}
	assert.Equal(t, 6.0, group.TotalCovered())

}

func Test_TotalUncovered_Normally_ReturnsTotal(t *testing.T) {
	group := pckg.Group{
		&pckg.Stats{
			Uncovered: 1.0,
		},
		&pckg.Stats{
			Uncovered: 2.0,
		},
		&pckg.Stats{
			Uncovered: 3.0,
		},
	}
	assert.Equal(t, 6.0, group.TotalUncovered())

}

func Test_GroupCoveragePercent_Normally_ReturnsTotal(t *testing.T) {
	group := pckg.Group{
		&pckg.Stats{
			Statements: 1.0,
			Covered:    1.0,
		},
		&pckg.Stats{
			Statements: 2.0,
			Covered:    1.0,
		},
		&pckg.Stats{
			Statements: 3.0,
			Covered:    1.0,
		},
	}
	assert.Equal(t, 0.5, group.CoveragePercent())

}
