package packages_test

import (
	"errors"
	"os"
	"testing"

	"github.com/eltorocorp/drygopher/drygopher/coverage/packages"
	"github.com/eltorocorp/drygopher/drygopher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type FileInfo interface {
	os.FileInfo
}

func Test_GetFileNamesForPackage_Normally_ReturnsFileNamesWithoutError(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)

	osioAPI.On("LookupEnv", mock.Anything).Return("path", true)
	fileInfo := new(mocks.FileInfo)
	fileInfo.On("Name").Return("filename.go")
	files := []os.FileInfo{
		fileInfo,
	}

	osioAPI.On("ReadDir", mock.Anything).Return(files, nil)

	packageAPI := packages.New(execAPI, osioAPI)
	fileNames, err := packageAPI.GetFileNamesForPackage("packagename")

	assert.Equal(t, []string{"path/src/packagename/filename.go"}, fileNames)
	assert.NoError(t, err)
}

func Test_GetFilesNamesForPackages_NoGOPATH_ReturnsError(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)

	osioAPI.On("LookupEnv", mock.Anything).Return("", false)

	packageAPI := packages.New(execAPI, osioAPI)
	fileNames, err := packageAPI.GetFileNamesForPackage("packagename")

	assert.Nil(t, fileNames)
	assert.EqualError(t, err, "GOPATH not set")
}

func Test_GetFileNamesForPackages_ErrorReadingDirectory_ReturnsError(t *testing.T) {
	osioAPI := new(mocks.OSIOAPI)
	execAPI := new(mocks.ExecAPI)

	osioAPI.On("LookupEnv", mock.Anything).Return("path", true)
	osioAPI.On("ReadDir", mock.Anything).Return(nil, errors.New("test error"))

	packageAPI := packages.New(execAPI, osioAPI)
	fileNames, err := packageAPI.GetFileNamesForPackage("packagename")

	assert.Nil(t, fileNames)
	assert.EqualError(t, err, "test error")
}
