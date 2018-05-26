package report_test

import (
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/eltorocorp/drygopher/drygopher/coverage/report"
	"github.com/eltorocorp/drygopher/drygopher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_BuildCoverageReport_Normally_BuildsAReport(t *testing.T) {
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)
	commandAPI.On("Run").Return(nil)
	execAPI.On("Command", mock.Anything, mock.Anything, mock.Anything).Return(commandAPI)

	reportAPI := report.New(execAPI)
	group := pckg.Group{
		&pckg.Stats{
			Package: "somepackage",
		},
	}
	exclusions := []string{
		"excludedpackage",
	}
	report, err := reportAPI.BuildCoverageReport(group, exclusions)

	expectedReport := `
Coverage Report
Packages Excluded From Coverage
===============================

excludedpackage
---------------

Analyzed Packages
-----------------
package    	stmts	cvrd	!cvrd	cvrg	est
somepackage	0	0	0	0.0%	no
           	0	0	0	0.0%	0
`
	assert.Equal(t, expectedReport, report)
	assert.NoError(t, err)
}
