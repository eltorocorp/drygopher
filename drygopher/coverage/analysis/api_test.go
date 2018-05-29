package analysis_test

import (
	"errors"
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/eltorocorp/drygopher/drygopher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Ensuring that stats are properly returned for a package that has tests
// and a package that does not have tests.
func Test_GetCoverageStatistics_Normally_ReturnsStats(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	rawPackageData := []string{"testedpkg:0.0,0.0 0 0"}
	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return(rawPackageData, false, nil).
		Once()

	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return([]string{}, false, nil).
		Once()

	rawAPI.On("AggregateRawPackageAnalysisData", mock.Anything, mock.Anything).
		Once().
		Return(&pckg.Stats{}, nil)

	analysisAPI := analysis.New(rawAPI)
	result, err := analysisAPI.GetCoverageStatistics([]string{"testedpkg", "untestedpkg"})
	assert.Len(t, result.TestedPackageStats, 1)
	assert.Len(t, result.UntestedPackageStats, 1)
	assert.NoError(t, err)
	rawAPI.AssertExpectations(t)
}

func Test_GetCoverageStatistics_EmptyOrWhiteSpacePackageName_SkipsPackage(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	analysisAPI := analysis.New(rawAPI)
	packageListWithWhiteSpacePackageNames := []string{
		"",
		"   ",
		"\t",
	}
	result, err := analysisAPI.GetCoverageStatistics(packageListWithWhiteSpacePackageNames)

	assert.Len(t, result.TestedPackageStats, 0)
	assert.Len(t, result.UntestedPackageStats, 0)
	assert.NoError(t, err)
	rawAPI.AssertNotCalled(t, "GetRawCoverageAnalysisForPackage")
}

func Test_GetCoverageStatistics_ErrorGettingRawAnalysis_ReturnsError(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return(nil, false, errors.New("test error"))

	analysisAPI := analysis.New(rawAPI)
	result, err := analysisAPI.GetCoverageStatistics([]string{"somepackage"})
	assert.Nil(t, result.TestedPackageStats)
	assert.Nil(t, result.UntestedPackageStats)
	assert.EqualError(t, err, "test error")
}

func Test_GetCoverageStatistics_ErrorAggregatingData_ReturnsError(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	rawPackageData := []string{"testedpkg:0.0,0.0 0 0"}
	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return(rawPackageData, false, nil)

	rawAPI.On("AggregateRawPackageAnalysisData", mock.Anything, mock.Anything).
		Return(nil, errors.New("test error"))

	analysisAPI := analysis.New(rawAPI)
	result, err := analysisAPI.GetCoverageStatistics([]string{"somepacakge"})
	assert.Nil(t, result.TestedPackageStats)
	assert.Nil(t, result.UntestedPackageStats)
	assert.EqualError(t, err, "test error")
}
