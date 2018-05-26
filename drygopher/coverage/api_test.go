package coverage_test

import (
	"errors"
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage"
	"github.com/eltorocorp/drygopher/drygopher/coverage/coverageerror"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/eltorocorp/drygopher/drygopher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AnalyzeUnitTestCoverage_Normally_ReturnsWithoutError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(pckg.Group{}, pckg.Group{}, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("OutputCoverageReport", mock.Anything, mock.Anything).Return(nil)

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile")
	assert.NoError(t, err)
}

func Test_AnalyzeUnitTestCoverage_GetPackagesErrors_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return(nil, errors.New("test error"))

	analysisAPI := new(mocks.AnalysisAPI)
	profileAPI := new(mocks.ProfileAPI)
	reportAPI := new(mocks.ReportAPI)
	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile")
	assert.EqualError(t, err, "test error")
}
func Test_AnalyzeUnitTestCoverage_GetCoverageStatsReturnError_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(nil, nil, errors.New("test error"))

	profileAPI := new(mocks.ProfileAPI)
	reportAPI := new(mocks.ReportAPI)
	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile")
	assert.EqualError(t, err, "test error")
}

func Test_AnalyzeTestCoverage_CoverageBelowStandard_ReturnsCoverageError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	testedPackageStats := pckg.Group{
		&pckg.Stats{
			Covered:    0,
			Estimated:  false,
			Statements: 1,
			Uncovered:  1,
		},
	}
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(testedPackageStats, pckg.Group{}, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("OutputCoverageReport", mock.Anything, mock.Anything).Return(nil)

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 100, false, "profile")
	assert.Error(t, err)
	assert.IsType(t, coverageerror.CoverageBelowStandard{}, err)

}
