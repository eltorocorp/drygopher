// Package analysis contains methods for building coverage statistics.
package analysis

import (
	"log"
	"strings"

	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis/analysistypes"
	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis/interfaces"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
)

// API contains methods gathering coverage statistics.
type API struct {
	raw interfaces.RawAPI
}

// New returns a reference to an API
func New(rawAPI interfaces.RawAPI) *API {
	return &API{
		raw: rawAPI,
	}
}

// GetCoverageStatistics gathers and returns coverage statistics for the specified packages.
func (a *API) GetCoverageStatistics(packages []string) (result analysistypes.GetCoverageStatisticsOutput, err error) {
	log.Println("Aggregating packages stats...")

	var testedPackageStats pckg.Group
	var untestedPackageStats pckg.Group
	testFailuresEncountered := false
	for _, pkg := range packages {
		if len(strings.TrimSpace(pkg)) == 0 {
			continue
		}

		failedTest := false
		var rawPkgCoverageData []string
		rawPkgCoverageData, failedTest, err = a.raw.GetRawCoverageAnalysisForPackage(pkg)
		if err != nil {
			return
		}

		if failedTest == true {
			testFailuresEncountered = true
		}

		if len(rawPkgCoverageData) == 0 {
			untestedPackageStats = append(untestedPackageStats, &pckg.Stats{
				Package:   pkg,
				Estimated: true,
			})
			continue
		}

		var packageStats *pckg.Stats
		packageStats, err = a.raw.AggregateRawPackageAnalysisData(pkg, rawPkgCoverageData)
		if err != nil {
			return
		}
		testedPackageStats = append(testedPackageStats, packageStats)
	}

	result.TestedPackageStats = testedPackageStats
	result.UntestedPackageStats = untestedPackageStats
	result.TestFailuresEncountered = testFailuresEncountered
	return
}
