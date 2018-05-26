package packages_test

import (
	"errors"
	"testing"

	"github.com/eltorocorp/drygopher/internal/coverage/packages"
	"github.com/eltorocorp/drygopher/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_GetPackages_Normally_ShellsCorrectListCommand(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)

	commandAPI.On("CombinedOutput").Return([]byte{}, nil)
	expectedCommandArgs := []interface{}{
		"sh",
		"-c",
		"go list ./... | grep -v /interfaces/ | grep -v /host/",
	}
	execAPI.On("Command", expectedCommandArgs...).Return(commandAPI)

	packageAPI := packages.New(execAPI, osioAPI)
	exclusionPatterns := []string{
		"/interfaces/",
		"/host/",
	}
	packageAPI.GetPackages(exclusionPatterns)

	execAPI.AssertExpectations(t)
}

func Test_GetPackages_ErrorFromShellOutput_ReturnsError(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)
	commandAPI := new(mocks.CommandAPI)

	commandAPI.On("CombinedOutput").Return(nil, errors.New("test error"))
	expectedCommandArgs := []interface{}{
		"sh",
		"-c",
		"go list ./...",
	}
	execAPI.On("Command", expectedCommandArgs...).Return(commandAPI)

	packageAPI := packages.New(execAPI, osioAPI)
	pckgs, err := packageAPI.GetPackages([]string{})

	assert.Nil(t, pckgs)
	assert.EqualError(t, err, "test error")
}
