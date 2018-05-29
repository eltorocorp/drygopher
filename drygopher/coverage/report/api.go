// Package report exposes methods for generating a coverage report.
package report

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eltorocorp/drygopher/drygopher/coverage/hostiface"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/gonum/floats"
	"github.com/willf/pad"
)

// API contains methods generating and printing a coverage report.
type API struct {
	execAPI hostiface.ExecAPI
}

// New returns a reference to a profile api.
func New(execAPI hostiface.ExecAPI) *API {
	return &API{
		execAPI: execAPI,
	}
}

var sb = new(strings.Builder)
var pl = func(a ...interface{}) { sb.WriteString(fmt.Sprintln(a...)) }
var pf = func(format string, v ...interface{}) { sb.WriteString(fmt.Sprintf(format, v...)) }

// BuildCoverageReport generates coverage report.
func (a *API) BuildCoverageReport(allPackages pckg.Group, exclusionPatterns []string) (string, error) {

	pl()
	pl("Coverage Report")
	pl("Packages Excluded From Coverage")
	pl("===============================")

	for _, exclusionPattern := range exclusionPatterns {
		pl() // space
		pl(exclusionPattern)
		pl(pad.Right("", len(exclusionPattern), "-"))
		output, err := a.getExcludedPackages(exclusionPattern)
		if err != nil {
			return "", err
		}
		pl(output)
	}

	pl()
	pl("Analyzed Packages")
	pl("-----------------")
	longestName := 0
	for _, p := range allPackages {
		if len(p.Package) > longestName {
			longestName = len(p.Package)
		}
	}
	format := "%v\t%v\t%v\t%v\t%v\t%v\n"
	pf(format, pad.Right("package", longestName, " "), "stmts", "cvrd", "!cvrd", "cvrg", "est")
	for _, p := range allPackages {
		packageName := pad.Right(p.Package, longestName, " ")
		pct := ftoa(p.CoveragePercent()*100) + "%"
		est := "no"
		if p.Estimated {
			est = "yes"
		}
		pf(format, packageName, p.Statements, p.Covered, p.Uncovered, pct, est)
	}

	pf(format,
		pad.Right("", longestName, " "),
		allPackages.TotalStatementCount(),
		allPackages.TotalCovered(),
		allPackages.TotalUncovered(),
		ftoa(allPackages.CoveragePercent()*100)+"%",
		allPackages.EstimateCount(),
	)
	return sb.String(), nil
}

func (a *API) getExcludedPackages(exclusionPattern string) (string, error) {
	// Appending the the '|| echo none' is important as it both squeches a
	// non-zero exit code should grep not return any results, and improves the
	// UX by including a value for exclusions that are noops.
	cmdtxt := fmt.Sprintf("go list ./... | grep %v || echo none", exclusionPattern)
	cmd := a.execAPI.Command("sh", "-c", cmdtxt)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(cmdtxt)
		return "", err
	}
	return string(output), err
}

func ftoa(f float64) string {
	return strconv.FormatFloat(floats.Round(f, 1), 'f', 1, 64)
}
