package analysis_test

import (
	"errors"
	"testing"

	"github.com/eltorocorp/drygopher/internal/coverage/analysis"
	"github.com/eltorocorp/drygopher/internal/mocks"
	"github.com/eltorocorp/drygopher/internal/pckg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Ensuring that stats are properly returned for a package that has tests
// and a package that does not have tests.
func Test_GetCoverageStatistics_Normally_ReturnsStats(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	rawPackageData := []string{"testedpkg:0.0,0.0 0 0"}
	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return(rawPackageData, nil).
		Once()

	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).Return([]string{}, nil).
		Once()

	rawAPI.On("AggregateRawPackageAnalysisData", mock.Anything, mock.Anything).
		Once().
		Return(&pckg.Stats{}, nil)

	analysisAPI := analysis.New(rawAPI)
	tested, untested, err := analysisAPI.GetCoverageStatistics([]string{"testedpkg", "untestedpkg"})
	assert.Len(t, tested, 1)
	assert.Len(t, untested, 1)
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
	tested, untested, err := analysisAPI.GetCoverageStatistics(packageListWithWhiteSpacePackageNames)

	assert.Len(t, tested, 0)
	assert.Len(t, untested, 0)
	assert.NoError(t, err)
	rawAPI.AssertNotCalled(t, "GetRawCoverageAnalysisForPackage")
}

func Test_GetCoverageStatistics_ErrorGettingRawAnalysis_ReturnsError(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return(nil, errors.New("test error"))

	analysisAPI := analysis.New(rawAPI)
	tested, untested, err := analysisAPI.GetCoverageStatistics([]string{"somepackage"})
	assert.Nil(t, tested)
	assert.Nil(t, untested)
	assert.EqualError(t, err, "test error")
}

func Test_GetCoverageStatistics_ErrorAggregatingData_ReturnsError(t *testing.T) {
	rawAPI := new(mocks.RawAPI)

	rawPackageData := []string{"testedpkg:0.0,0.0 0 0"}
	rawAPI.On("GetRawCoverageAnalysisForPackage", mock.Anything).
		Return(rawPackageData, nil)

	rawAPI.On("AggregateRawPackageAnalysisData", mock.Anything, mock.Anything).
		Return(nil, errors.New("test error"))

	analysisAPI := analysis.New(rawAPI)
	tested, untested, err := analysisAPI.GetCoverageStatistics([]string{"somepacakge"})
	assert.Nil(t, tested)
	assert.Nil(t, untested)
	assert.EqualError(t, err, "test error")
}
