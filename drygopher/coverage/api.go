// Package coverage contains an API for working with unit test coverage analysis.
package coverage

import (
	"fmt"
	"log"

	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis/analysistypes"

	"github.com/eltorocorp/drygopher/drygopher/coverage/coverageerrors"
	"github.com/eltorocorp/drygopher/drygopher/coverage/interfaces"
)

// API contains methods for analyzing unit test coverage.
type API struct {
	profile  interfaces.ProfileAPI
	packages interfaces.PackageAPI
	report   interfaces.ReportAPI
	analysis interfaces.AnalysisAPI
}

// New returns a reference to a new coverage API.
func New(packageAPI interfaces.PackageAPI, analysisAPI interfaces.AnalysisAPI, profileAPI interfaces.ProfileAPI, reportAPI interfaces.ReportAPI) *API {
	return &API{
		packages: packageAPI,
		profile:  profileAPI,
		report:   reportAPI,
		analysis: analysisAPI,
	}
}

// AnalyzeUnitTestCoverage analyzes unit test coverage across all packages.
func (a *API) AnalyzeUnitTestCoverage(exclusionPatterns []string, coverageStandard float64, suppressProfile bool, coverageProfileName string, suppressPercentageFile bool) (err error) {
	log.Println("Analyzing unit test coverage...")
	var (
		packages      []string
		coverageStats analysistypes.GetCoverageStatisticsOutput
	)

	packages, err = a.packages.GetPackages(exclusionPatterns)
	if err != nil {
		return
	}

	coverageStats, err = a.analysis.GetCoverageStatistics(packages)
	if err != nil {
		return
	}

	coverageStats.UntestedPackageStats.SetEstimatedStmtCountFrom(coverageStats.TestedPackageStats)
	allPackages := append(coverageStats.TestedPackageStats, coverageStats.UntestedPackageStats...)
	actualCoveragePercentage := allPackages.CoveragePercent()

	var report string
	report, err = a.report.BuildCoverageReport(allPackages, exclusionPatterns)
	if err != nil {
		return
	}
	fmt.Println(report)

	if !suppressProfile {
		err = a.profile.BuildAndSaveCoverageProfile(allPackages, coverageProfileName)
	}
	if err != nil {
		return
	}

	if !suppressPercentageFile {
		err = a.profile.OutputPercentageFile(100.0 * actualCoveragePercentage)
	}
	if err != nil {
		return
	}

	if actualCoveragePercentage*100.0 < coverageStandard {
		err = coverageerrors.NewCoverageBelowStandardError(coverageStandard, 100.0*actualCoveragePercentage)
	}

	// A TestFailuresEncountered error is allowed to overwrite (supersede) a
	// CoverageBelowStandardError
	if coverageStats.TestFailuresEncountered {
		err = coverageerrors.NewUnitTestFailedError()
	}

	return
}
