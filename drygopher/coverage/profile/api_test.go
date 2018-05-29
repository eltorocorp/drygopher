package profile_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
	"github.com/eltorocorp/drygopher/drygopher/coverage/profile"
	"github.com/eltorocorp/drygopher/drygopher/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_BuildAndSaveCoverageProfile_Normally_NoError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	osioAPI := new(mocks.OSIOAPI)

	packageAPI.On("GetFileNamesForPackage", mock.Anything).Return([]string{"somepackage.go"}, nil)
	osioAPI.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	profileAPI := profile.New(packageAPI, osioAPI)
	group := pckg.Group{
		&pckg.Stats{
			Estimated: true,
		},
	}
	err := profileAPI.BuildAndSaveCoverageProfile(group, "coverage.out")

	assert.NoError(t, err)
}

func Test_BuildAndSaveCoverageProfile_ErrorGettingFileNames_ReturnsError(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	osioAPI := new(mocks.OSIOAPI)

	packageAPI.On("GetFileNamesForPackage", mock.Anything).Return(nil, errors.New("test error"))

	profileAPI := profile.New(packageAPI, osioAPI)
	group := pckg.Group{
		&pckg.Stats{
			Estimated: true,
		},
	}
	err := profileAPI.BuildAndSaveCoverageProfile(group, "coverage.out")

	assert.EqualError(t, err, "test error")
}

func Test_OutputPercentageFile_Normally_WritesCorrectlyToFile(t *testing.T) {
	packageAPI := new(mocks.PackageAPI)
	osioAPI := new(mocks.OSIOAPI)
	osioAPI.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	profileAPI := profile.New(packageAPI, osioAPI)
	suppliedPercentage := 12.3

	err := profileAPI.OutputPercentageFile(suppliedPercentage)

	expectedFileName := "coveragepct"
	expectedFilePermissions := os.ModePerm
	expectedPercentage := []byte(fmt.Sprintf("%.2f", suppliedPercentage))
	assert.NoError(t, err)
	osioAPI.AssertCalled(t, "WriteFile", expectedFileName, expectedPercentage, expectedFilePermissions)
}
