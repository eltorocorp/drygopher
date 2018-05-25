package analysis

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eltorocorp/drygopher/internal/hostiface"
	"github.com/eltorocorp/drygopher/internal/pckg"
)

// API contains methods gathering coverage statistics.
type API struct {
	osioAPI hostiface.OSIOAPI
	execAPI hostiface.ExecAPI
}

// New returns a reference to an API
func New(osioAPI hostiface.OSIOAPI, execAPI hostiface.ExecAPI) *API {
	return &API{
		osioAPI: osioAPI,
		execAPI: execAPI,
	}
}

// GetRawCoverageAnalysisForPackage shells out a go test command and returns the
// raw data from the resulting coverage profile.
func (a *API) GetRawCoverageAnalysisForPackage(pkg string) ([]string, error) {
	covermode := "count"
	analyzeCmdText := "go test -covermode=%v -coverprofile=tmp.out %v"
	analyzeCmdText = fmt.Sprintf(analyzeCmdText, covermode, pkg)
	analyzeCoverageCmd := a.execAPI.Command("sh", "-c", analyzeCmdText)
	var result []byte
	result, err := analyzeCoverageCmd.Output()
	if err != nil {
		log.Println("Error issuing command to analyze package.")
		log.Println(analyzeCmdText)
		log.Println(err)
		return nil, err
	}
	log.Printf("---> Package result: %v", string(result))
	if result[0] == '?' {
		return []string{}, nil
	}
	var tmpOut []byte
	tmpOut, err = a.osioAPI.ReadFile("tmp.out")
	if err != nil {
		return nil, err
	}
	rawPkgCoverageData := strings.Split(string(tmpOut), "\n")
	rawPkgCoverageData = rawPkgCoverageData[1:]
	return rawPkgCoverageData, nil
}

// GetCoverageStatistics gathers and returns coverage statistics for the specified packages.
func (a *API) GetCoverageStatistics(packages []string) (testedPackageStats, untestedPackageStats pckg.Group, err error) {
	log.Println("Aggregating packages stats...")

	for _, pkg := range packages {
		if pkg == "" {
			continue
		}

		var rawPkgCoverageData []string
		rawPkgCoverageData, err = a.GetRawCoverageAnalysisForPackage(pkg)
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

		packageStats := a.aggregateRawPackageAnalysisData(pkg, rawPkgCoverageData)
		testedPackageStats = append(testedPackageStats, packageStats)

	}
	return
}

func (a *API) aggregateRawPackageAnalysisData(pkg string, rawPkgCoverageData []string) *pckg.Stats {
	totalStatementCount := 0.0
	totalCoveredCount := 0.0
	firstLine := true
	for _, coverageReportItem := range rawPkgCoverageData {
		if firstLine {
			firstLine = false
			continue
		}
		if coverageReportItem == "" {
			continue
		}
		statementCount := a.parseStatementCountFromRaw(coverageReportItem)
		totalStatementCount += statementCount
		covered := 0.0
		if a.parseCallCountFromRaw(coverageReportItem) > 0 {
			covered = statementCount
		}
		totalCoveredCount += covered
	}
	return &pckg.Stats{
		Covered:         totalCoveredCount,
		Estimated:       false,
		Package:         pkg,
		Statements:      totalStatementCount,
		Uncovered:       totalStatementCount - totalCoveredCount,
		RawCoverageData: rawPkgCoverageData,
	}
}

func (a *API) parseStatementCountFromRaw(rawDatum string) float64 {
	result, err := strconv.ParseFloat(strings.Split(rawDatum, " ")[1], 64)
	if err != nil {
		log.Println(err)
		panic("Error parsing raw data for statement count.")
	}
	return result
}

func (a *API) parseCallCountFromRaw(rawDatum string) float64 {
	result, err := strconv.ParseFloat(strings.Split(rawDatum, " ")[2], 64)
	if err != nil {
		log.Println(err)
		panic("Error parsing raw data for call count.")
	}
	return result
}
