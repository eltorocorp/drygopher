package analysis

import (
	"log"

	"github.com/eltorocorp/drygopher/internal/coverage/analysis/interfaces"
	"github.com/eltorocorp/drygopher/internal/pckg"
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
func (a *API) GetCoverageStatistics(packages []string) (testedPackageStats, untestedPackageStats pckg.Group, err error) {
	log.Println("Aggregating packages stats...")

	for _, pkg := range packages {
		if pkg == "" {
			continue
		}

		var rawPkgCoverageData []string
		rawPkgCoverageData, err = a.raw.GetRawCoverageAnalysisForPackage(pkg)
		if err != nil {
			return
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
	return
}
