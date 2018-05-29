package coverage_test

import (
	"errors"
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage"
	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis/analysistypes"
	"github.com/eltorocorp/drygopher/drygopher/coverage/coverageerrors"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/eltorocorp/drygopher/drygopher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AnalyzeUnitTestCoverage_Normally_ReturnsWithoutError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(analysistypes.GetCoverageStatisticsOutput{}, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", nil)

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile", true)
	assert.NoError(t, err)
}

func Test_AnalyzeUnitTestCoverage_GetPackagesErrors_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return(nil, errors.New("test error"))

	analysisAPI := new(mocks.AnalysisAPI)
	profileAPI := new(mocks.ProfileAPI)
	reportAPI := new(mocks.ReportAPI)
	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile", true)
	assert.EqualError(t, err, "test error")
}

func Test_AnalyzeUnitTestCoverage_GetCoverageStatsReturnError_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(analysistypes.GetCoverageStatisticsOutput{}, errors.New("test error"))

	profileAPI := new(mocks.ProfileAPI)
	reportAPI := new(mocks.ReportAPI)
	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile", true)
	assert.EqualError(t, err, "test error")
}

func Test_AnalyzeTestCoverage_CoverageBelowStandard_ReturnsCoverageError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)

	getCoverageStatisticsResult := analysistypes.GetCoverageStatisticsOutput{
		TestedPackageStats: pckg.Group{
			&pckg.Stats{
				Covered:    0,
				Estimated:  false,
				Statements: 1,
				Uncovered:  1,
			},
		},
	}

	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(getCoverageStatisticsResult, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", nil)

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 100, false, "profile", true)
	assert.Error(t, err)
	assert.IsType(t, coverageerrors.CoverageBelowStandard{}, err)
}

func Test_AnalyzeTestCoverage_ErrorBuildingCoverageReport_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(analysistypes.GetCoverageStatisticsOutput{}, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", errors.New("test error"))

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 100, false, "profile", true)

	assert.EqualError(t, err, "test error")
}

func Test_AnalyzeTestCoverage_ErrorBuildingCoverageProfile_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(analysistypes.GetCoverageStatisticsOutput{}, nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(errors.New("test error"))

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 100, false, "profile", true)

	assert.EqualError(t, err, "test error")
}

func Test_AnalyzeTestCoverage_ErrorOutputingPercentageFile_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(analysistypes.GetCoverageStatisticsOutput{}, nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)
	profileAPI.On("OutputPercentageFile", mock.Anything).Return(errors.New("test error"))

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 100, false, "profile", false)

	assert.EqualError(t, err, "test error")
}

func Test_AnalyzeUnitTestCoverage_UnitTestError_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)
	getCoverageStatisticsResult := analysistypes.GetCoverageStatisticsOutput{
		TestFailuresEncountered: true,
	}
	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(getCoverageStatisticsResult, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", nil)

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 0, false, "profile", true)

	assert.EqualError(t, err, coverageerrors.NewUnitTestFailedError().Error())
}

func Test_AnalyzeUnitTestCoverage_UnitTestError_SupersedesCoverageError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	packageAPI.On("GetPackages", mock.Anything).Return([]string{}, nil)

	analysisAPI := new(mocks.AnalysisAPI)

	getCoverageStatisticsResult := analysistypes.GetCoverageStatisticsOutput{
		TestedPackageStats: pckg.Group{
			&pckg.Stats{
				Covered:    0,
				Estimated:  false,
				Statements: 1,
				Uncovered:  1,
			},
		},
		TestFailuresEncountered: true,
	}

	analysisAPI.On("GetCoverageStatistics", mock.Anything).Return(getCoverageStatisticsResult, nil)

	profileAPI := new(mocks.ProfileAPI)
	profileAPI.On("BuildAndSaveCoverageProfile", mock.Anything, mock.Anything).Return(nil)

	reportAPI := new(mocks.ReportAPI)
	reportAPI.On("BuildCoverageReport", mock.Anything, mock.Anything).Return("", nil)

	coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

	err := coverageAPI.AnalyzeUnitTestCoverage([]string{}, 100, false, "profile", true)

	assert.Error(t, err)
	assert.IsType(t, coverageerrors.UnitTestFailed{}, err)
}
