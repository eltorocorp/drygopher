package coverage

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eltorocorp/drygopher/internal/coverage/interfaces"
	"github.com/eltorocorp/drygopher/internal/pckg"
)

// API contains methods for analyzing unit test coverage.
type API struct {
	tooling interfaces.CoverageToolingAPI
}

// New returns a reference to a new coverage API.
func New(tooling interfaces.CoverageToolingAPI) *API {
	return &API{
		tooling: tooling,
	}
}

// AnalyzeUnitTestCoverage analyzes unit test coverage across all packages.
func (a *API) AnalyzeUnitTestCoverage(exclusionPatterns []string, coverageStandard float64) (err error) {
	log.Println("Analyzing unit test coverage...")
	var (
		packages             []string
		testedPackageStats   pckg.Group
		untestedPackageStats pckg.Group
	)

	packages, err = a.tooling.GetPackages(exclusionPatterns)
	if err != nil {
		return
	}

	testedPackageStats, untestedPackageStats, err = a.getCoverageStatistics(packages)
	if err != nil {
		return
	}

	untestedPackageStats.SetEstimatedStmtCntFrom(testedPackageStats)
	allPackages := append(testedPackageStats, untestedPackageStats...)
	actualCoveragePercentage := allPackages.CoveragePercent()

	a.outputCoverageReport(allPackages, exclusionPatterns)

	if actualCoveragePercentage < coverageStandard {
		return fmt.Errorf("coverage of %v%% is below the standard of %v%%", actualCoveragePercentage, coverageStandard)
	}
	return nil
}

func (a *API) getCoverageStatistics(packages []string) (testedPackageStats, untestedPackageStats pckg.Group, err error) {
	log.Println("Aggregating packages stats...")

	for _, pkg := range packages {
		if pkg == "" {
			continue
		}

		var rawPkgCoverageData []string
		rawPkgCoverageData, err = a.tooling.GetRawCoverageAnalysisForPackage(pkg)
		if err != nil {
			return
		}

		if len(rawPkgCoverageData) == 0 {
			untestedPackageStats = append(untestedPackageStats, &pckg.Stats{
				Package:   pkg,
				Estimated: true,
			})
			continue
		}

		packageStats := a.aggregateRawPackageAnalysisData(pkg, rawPkgCoverageData)
		testedPackageStats = append(testedPackageStats, packageStats)

	}
	return
}

func (a *API) aggregateRawPackageAnalysisData(pkg string, rawPkgCoverageData []string) *pckg.Stats {
	totalStatementCount := 0.0
	totalCoveredCount := 0.0
	firstLine := true
	for _, coverageReportItem := range rawPkgCoverageData {
		if firstLine {
			firstLine = false
			continue
		}
		if coverageReportItem == "" {
			continue
		}
		statementCount := a.parseStatementCountFromRaw(coverageReportItem)
		totalStatementCount += statementCount
		covered := 0.0
		if a.parseCallCountFromRaw(coverageReportItem) > 0 {
			covered = statementCount
		}
		totalCoveredCount += covered
	}
	return &pckg.Stats{
		Covered:    totalCoveredCount,
		Estimated:  false,
		Package:    pkg,
		Statements: totalStatementCount,
		Uncovered:  totalStatementCount - totalCoveredCount,
	}
}

func (a *API) parseStatementCountFromRaw(rawDatum string) float64 {
	// raw: github.com/eltorocorp/cookiecrumbler-v2/src/shared/clock/api.go:6.34,6.72 1 1
	result, err := strconv.ParseFloat(strings.Split(rawDatum, " ")[1], 64)
	if err != nil {
		log.Println(err)
		panic("Error parsing raw data for statement count.")
	}
	return result
}

func (a *API) parseCallCountFromRaw(rawDatum string) float64 {
	// raw: github.com/eltorocorp/cookiecrumbler-v2/src/shared/clock/api.go:6.34,6.72 1 1
	result, err := strconv.ParseFloat(strings.Split(rawDatum, " ")[2], 64)
	if err != nil {
		log.Println(err)
		panic("Error parsing raw data for call count.")
	}
	return result
}
