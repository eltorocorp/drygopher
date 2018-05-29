// Package profile exposes methods for working with coverage profile files and other profile related files.
package profile

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/eltorocorp/drygopher/drygopher/coverage/hostiface"
	"github.com/eltorocorp/drygopher/drygopher/coverage/interfaces"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
)

// API contains methods for creating a coverage profile.
type API struct {
	packages interfaces.PackageAPI
	osioAPI  hostiface.OSIOAPI
}

// New returns a reference to a profile api.
func New(packageAPI interfaces.PackageAPI, osioAPI hostiface.OSIOAPI) *API {
	return &API{
		packages: packageAPI,
		osioAPI:  osioAPI,
	}
}

// BuildAndSaveCoverageProfile builds and saves a coverage profile that can be consumed by go tool cover -html.
func (a *API) BuildAndSaveCoverageProfile(allPackages pckg.Group, coverageProfileName string) error {
	profileData := []string{}
	for _, pkg := range allPackages {
		if pkg.Estimated {
			fileNames, err := a.packages.GetFileNamesForPackage(pkg.Package)
			if err != nil {
				return err
			}
			for _, fileName := range fileNames {
				fileProfileInfo := fileName + ":0.0,0.0 0 0"
				pkg.RawCoverageData = append(pkg.RawCoverageData, fileProfileInfo)
			}
		}
		profileData = append(profileData, pkg.RawCoverageData...)
	}
	return a.saveCoverageProfile(coverageProfileName, profileData)
}

// OutputPercentageFile saves a file containing the supplied coverage percentage.
func (a *API) OutputPercentageFile(coveragePercentage float64) error {
	pctString := fmt.Sprintf("%.2f", coveragePercentage)
	return a.osioAPI.WriteFile("coveragepct", []byte(pctString), os.ModePerm)
}

// SaveCoverageProfile saves the supplied raw data to the desired file.
func (a *API) saveCoverageProfile(fileName string, rawData []string) error {
	sort.StringSlice(rawData).Sort()
	profile := "mode: count\n" + strings.TrimSpace(strings.Join(rawData, "\n"))
	return a.osioAPI.WriteFile(fileName, []byte(profile), os.ModePerm)
}
