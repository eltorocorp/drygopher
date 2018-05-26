package profile_test

import (
	"errors"
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
