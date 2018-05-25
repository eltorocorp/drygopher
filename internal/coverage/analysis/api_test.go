package analysis_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"

	"github.com/eltorocorp/drygopher/internal/coverage/analysis"
	"github.com/eltorocorp/drygopher/internal/mocks"
)

func Test_GetRawCoverageAnalysisForPackage_Normally_ReturnsAnalysis(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)

	commandAPI.On("Output").Return([]byte("testresult"), nil)
	execAPI.On("Command", mock.Anything, mock.Anything, mock.Anything).Return(commandAPI)

	profileData := "headerthatshouldgetskipped\ndata"
	osioAPI.On("ReadFile", mock.Anything).Return([]byte(profileData), nil)

	analysisAPI := analysis.New(osioAPI, execAPI)
	result, err := analysisAPI.GetRawCoverageAnalysisForPackage("somepkg")
	assert.NoError(t, err)
	assert.Equal(t, []string{"data"}, result)
}

func Test_GetRawCoverageAnalysisForPackage_TestCmdError_ReturnsError(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)

	commandAPI.On("Output").Return(nil, errors.New("test error"))
	execAPI.On("Command", mock.Anything, mock.Anything, mock.Anything).Return(commandAPI)

	analysisAPI := analysis.New(osioAPI, execAPI)
	_, err := analysisAPI.GetRawCoverageAnalysisForPackage("somepkg")

	assert.EqualError(t, err, "test error")
}

func Test_GetRawCoverageAnalysisForPackage_NoTestsForPackage_ReturnsEmptySlice(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)

	// go test will output a line that starts with a question mark for packages
	// that have no test files.
	commandAPI.On("Output").Return([]byte("?"), nil)
	execAPI.On("Command", mock.Anything, mock.Anything, mock.Anything).Return(commandAPI)

	analysisAPI := analysis.New(osioAPI, execAPI)
	result, err := analysisAPI.GetRawCoverageAnalysisForPackage("somepkg")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, result)
}

func Test_GetRawCoverageAnalysisForPackage_ErrorReadingTmpFile_ReturnsError(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)

	commandAPI.On("Output").Return([]byte("testresult"), nil)
	execAPI.On("Command", mock.Anything, mock.Anything, mock.Anything).Return(commandAPI)

	osioAPI.On("ReadFile", mock.Anything).Return(nil, errors.New("test error"))

	analysisAPI := analysis.New(osioAPI, execAPI)
	_, err := analysisAPI.GetRawCoverageAnalysisForPackage("somepkg")

	assert.EqualError(t, err, "test error")
}
