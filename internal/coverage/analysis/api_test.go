package analysis_test

import (
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
