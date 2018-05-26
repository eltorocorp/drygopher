package raw

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eltorocorp/drygopher/drygopher/coverage/hostiface"
	"github.com/eltorocorp/drygopher/drygopher/coverage/pckg"
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
	if !strings.Contains(rawPkgCoverageData[0], "mode:") {
		return nil, errors.New("coverage profile file malformed; missing 'mode:' in header")
	}
	rawPkgCoverageData = rawPkgCoverageData[1:]
	return rawPkgCoverageData, nil
}

// AggregateRawPackageAnalysisData reduces raw coveraeg data to a structured
// object that contains summarized information about the analysis.
func (a *API) AggregateRawPackageAnalysisData(pkg string, rawPkgCoverageData []string) (*pckg.Stats, error) {
	totalStatementCount := 0.0
	totalCoveredCount := 0.0
	for line, coverageReportItem := range rawPkgCoverageData {
		if line == 0 || coverageReportItem == "" {
			continue
		}

		statementCount, err := a.parseStatementCountFromRaw(coverageReportItem)
		if err != nil {
			return nil, err
		}

		totalStatementCount += statementCount
		callCount, err := a.parseCallCountFromRaw(coverageReportItem)
		if err != nil {
			return nil, err
		}

		// If a set of statements has at least one call, we set the number of covered
		// statements equal to the number of statements (rather than the number
		// of calls). This prevents coverage from going above 100% in cases
		// where a set of statements are called multiple times through the course
		// of running tests.
		covered := 0.0
		if callCount > 0 {
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
	}, nil
}

func (a *API) parseStatementCountFromRaw(rawDatum string) (float64, error) {
	return strconv.ParseFloat(strings.Split(rawDatum, " ")[1], 64)
}

func (a *API) parseCallCountFromRaw(rawDatum string) (float64, error) {
	return strconv.ParseFloat(strings.Split(rawDatum, " ")[2], 64)
}
