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
	host interfaces.HostAPI
}

// New returns a reference to a new coverage API.
func New(hostAPI interfaces.HostAPI) *API {
	return &API{
		host: hostAPI,
	}
}

// AnalyzeUnitTestCoverage analyzes unit test coverage across all packages.
func (a *API) AnalyzeUnitTestCoverage(exclusionPatterns []string, coverageStandard float64, suppressProfile bool, coverageProfileName string) (err error) {
	log.Println("Analyzing unit test coverage...")
	var (
		packages             []string
		testedPackageStats   pckg.Group
		untestedPackageStats pckg.Group
	)

	packages, err = a.host.GetPackages(exclusionPatterns)
	if err != nil {
		return
	}

	testedPackageStats, untestedPackageStats, err = a.getCoverageStatistics(packages)
	if err != nil {
		return
	}

	untestedPackageStats.SetEstimatedStmtCountFrom(testedPackageStats)
	allPackages := append(testedPackageStats, untestedPackageStats...)
	actualCoveragePercentage := allPackages.CoveragePercent()

	a.outputCoverageReport(allPackages, exclusionPatterns)

	if !suppressProfile {
		err = a.buildAndSaveCoverageProfile(allPackages, coverageProfileName)
	}

	if actualCoveragePercentage < coverageStandard {
		return fmt.Errorf("coverage of %v%% is below the standard of %v%%", actualCoveragePercentage, coverageStandard)
	}
	return
}

func (a *API) buildAndSaveCoverageProfile(allPackages pckg.Group, coverageProfileName string) error {
	profileData := []string{}
	for _, pkg := range allPackages {
		if pkg.Estimated {
			fileNames, err := a.host.GetFileNamesForPackage(pkg.Package)
			if err != nil {
				return err
			}
			for _, fileName := range fileNames {
				fileProfileInfo := fileName + ":0.0,0.0 0 0"
				pkg.RawCoverageData = append(pkg.RawCoverageData, fileProfileInfo)
			}
		}
		profileData = append(profileData, pkg.RawCoverageData...)
	}
	return a.host.SaveCoverageProfile(coverageProfileName, profileData)
}

func (a *API) getCoverageStatistics(packages []string) (testedPackageStats, untestedPackageStats pckg.Group, err error) {
	log.Println("Aggregating packages stats...")

	for _, pkg := range packages {
		if pkg == "" {
			continue
		}

		var rawPkgCoverageData []string
		rawPkgCoverageData, err = a.host.GetRawCoverageAnalysisForPackage(pkg)
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
		Covered:         totalCoveredCount,
		Estimated:       false,
		Package:         pkg,
		Statements:      totalStatementCount,
		Uncovered:       totalStatementCount - totalCoveredCount,
		RawCoverageData: rawPkgCoverageData,
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
