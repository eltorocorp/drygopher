// Package raw exposes lower level methods that interact with the 'go' CLI.
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

// GetRawCoverageAnalysisForPackage shells out a go test command and returns
// three values:
// - The raw coverage results for the package
// - Whether or not the package's unit tests failed.
// - An error, should one have occurred during processing.
func (a *API) GetRawCoverageAnalysisForPackage(pkg string) (rawResult []string, testFailed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an unhandled panic was encountered while executing a test against package '%v'", pkg)
			testFailed = true
			rawResult = []string{""}
		}
		a.osioAPI.MustRemove("tmp.out")
	}()

	covermode := "atomic"
	// the use of -count=1 is important, as it prevents go test from attempting
	// to use cached test results (which are not supported by drygopher)
	analyzeCmdText := "go test -count=1 -race -covermode=%v -coverprofile=tmp.out %v"
	analyzeCmdText = fmt.Sprintf(analyzeCmdText, covermode, pkg)
	analyzeCoverageCmd := a.execAPI.Command("sh", "-c", analyzeCmdText)

	var resultBytes []byte
	resultBytes, err = analyzeCoverageCmd.Output()
	result := string(resultBytes)

	// analyzeCoverageCmd.Output() will return an error when a test fails.
	// We check for this situation first, as, in this case, we will
	// disregard the error and instead process the test failure.
	if strings.Contains(result, "FAIL:") {
		log.Printf("---> Failed test: %v\n", result)
		rawResult, err = a.retrieveResultsFromTmpFile(result)
		testFailed = true
		return
	}

	// If we're here, and there is an error, we know the error didn't result from
	// a failed unit test. So we bubble the error up.
	if err != nil {
		return nil, false, err
	}

	log.Printf("---> Package result: %v", string(result))
	rawResult, err = a.retrieveResultsFromTmpFile(result)
	return rawResult, false, err
}

func (a *API) retrieveResultsFromTmpFile(result string) ([]string, error) {
	if result[0] == '?' {
		return []string{}, nil
	}
	var tmpOut []byte
	tmpOut, err := a.osioAPI.ReadFile("tmp.out")
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
	for _, coverageReportItem := range rawPkgCoverageData {
		if coverageReportItem == "" {
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
