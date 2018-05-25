package report

import (
	"fmt"
	"log"
	"strconv"

	"github.com/eltorocorp/drygopher/internal/hostiface"
	"github.com/eltorocorp/drygopher/internal/pckg"
	"github.com/willf/pad"
	"gonum.org/v1/gonum/floats"
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

// OutputCoverageReport generates and outputs a coverage report.
func (a *API) OutputCoverageReport(allPackages pckg.Group, exclusionPatterns []string) {
	fmt.Println() // space
	log.Println("Coverage Report")
	log.Println("Packages Excluded From Coverage")
	log.Println("===============================")

	for _, exclusionPattern := range exclusionPatterns {
		fmt.Println() // space
		log.Println(exclusionPattern)
		log.Println(pad.Right("", len(exclusionPattern), "-"))
		a.PrintExcludedPackages(exclusionPattern)
	}

	fmt.Println() // space
	log.Println("Analyzed Packages")
	log.Println("-----------------")
	longestName := 0
	for _, p := range allPackages {
		if len(p.Package) > longestName {
			longestName = len(p.Package)
		}
	}
	format := "\t%v\t%v\t%v\t%v\t%v\t%v\n"
	log.Printf(format, pad.Right("package", longestName, " "), "stmts", "cvrd", "!cvrd", "cvrg", "est")
	for _, p := range allPackages {
		packageName := pad.Right(p.Package, longestName, " ")
		pct := ftoa(p.CoveragePercent()*100) + "%"
		est := "no"
		if p.Estimated {
			est = "yes"
		}
		log.Printf(format, packageName, p.Statements, p.Covered, p.Uncovered, pct, est)
	}

	log.Printf(format,
		pad.Right("", longestName, " "),
		allPackages.TotalStatementCount(),
		allPackages.TotalCovered(),
		allPackages.TotalUncovered(),
		ftoa(allPackages.CoveragePercent()*100)+"%",
		allPackages.EstimateCount(),
	)
}

// PrintExcludedPackages shells out a go list command and sends the results
// of the command directly to stdout.
func (a *API) PrintExcludedPackages(exclusionPattern string) {
	cmd := a.execAPI.Command("sh", "-c", fmt.Sprintf("go list ./... | grep -v /vendor/ | grep %v", exclusionPattern))
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func ftoa(f float64) string {
	return strconv.FormatFloat(floats.Round(f, 1), 'f', 1, 64)
}
