package coverage

import (
	"log"

	"github.com/eltorocorp/drygopher/coverageerror"
	"github.com/eltorocorp/drygopher/internal/coverage/interfaces"
	"github.com/eltorocorp/drygopher/internal/pckg"
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
func (a *API) AnalyzeUnitTestCoverage(exclusionPatterns []string, coverageStandard float64, suppressProfile bool, coverageProfileName string) (err error) {
	log.Println("Analyzing unit test coverage...")
	var (
		packages             []string
		testedPackageStats   pckg.Group
		untestedPackageStats pckg.Group
	)

	packages, err = a.packages.GetPackages(exclusionPatterns)
	if err != nil {
		return
	}

	testedPackageStats, untestedPackageStats, err = a.analysis.GetCoverageStatistics(packages)
	if err != nil {
		return
	}

	untestedPackageStats.SetEstimatedStmtCountFrom(testedPackageStats)
	allPackages := append(testedPackageStats, untestedPackageStats...)
	actualCoveragePercentage := allPackages.CoveragePercent()

	a.report.OutputCoverageReport(allPackages, exclusionPatterns)

	if !suppressProfile {
		err = a.profile.BuildAndSaveCoverageProfile(allPackages, coverageProfileName)
	}

	if actualCoveragePercentage < coverageStandard {
		err = coverageerror.New(coverageStandard, actualCoveragePercentage)
	}
	return
}
