package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eltorocorp/drygopher/internal/pckg"
)

// AnalyzeUnitTestCoverage analyzes unit test coverage across all packages.
func AnalyzeUnitTestCoverage(exclusionPatterns []string, coverageStandard float64) (err error) {
	log.Println("Analyzing unit test coverage...")
	var (
		packages         []string
		testedPackages   pckg.Group
		untestedPackages pckg.Group
	)

	packages, err = getPackages(exclusionPatterns)
	if err != nil {
		return
	}

	testedPackages, untestedPackages, err = aggregatePackageStats(packages)
	if err != nil {
		return
	}

	untestedPackages.SetEstimatedStmtCntFrom(testedPackages)
	allPackages := append(testedPackages, untestedPackages...)
	actualCoveragePercentage := allPackages.CoveragePercent()
	outputCoverageReport(allPackages, exclusionPatterns)
	if actualCoveragePercentage < coverageStandard {
		return fmt.Errorf("coverage of %v%% is below the standard of %v%%", actualCoveragePercentage, coverageStandard)
	}
	return nil
}

func aggregatePackageStats(packages []string) (testedPackages, untestedPackages pckg.Group, err error) {
	log.Println("Aggregating packages stats...")
	covermode := "count"
	initCmdText := fmt.Sprintf("echo 'mode: %v' > coverage-all.out", covermode)
	initializeCoverageOutCmd := exec.Command("sh", "-c", initCmdText)

	err = initializeCoverageOutCmd.Run()
	if err != nil {
		return
	}
	for _, pkg := range packages {
		if pkg == "" {
			continue
		}

		analyzeCmdText := "go test -covermode=%v -coverprofile=tmp.out %v"
		analyzeCmdText = fmt.Sprintf(analyzeCmdText, covermode, pkg)
		analyzeCoverageCmd := exec.Command("sh", "-c", analyzeCmdText)
		var result []byte
		result, err = analyzeCoverageCmd.Output()
		if err != nil {
			log.Println("Error issueing command to analyze package.")
			log.Println(analyzeCmdText)
			return
		}
		log.Printf("---> Package result: %v", string(result))
		if result[0] == '?' {
			untestedPackages = append(untestedPackages, &pckg.Stats{
				Package:   pkg,
				Estimated: true,
			})
			continue
		}
		var tmpOut []byte
		tmpOut, err = ioutil.ReadFile("tmp.out")
		if err != nil {
			return
		}
		rawPkgCoverageData := strings.Split(string(tmpOut), "\n")
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
			statementCount := parseStatementCountFromRaw(coverageReportItem)
			totalStatementCount += statementCount
			covered := 0.0
			if parseCallCountFromRaw(coverageReportItem) > 0 {
				covered = statementCount
			}
			totalCoveredCount += covered
		}
		testedPackages = append(testedPackages, &pckg.Stats{
			Covered:    totalCoveredCount,
			Estimated:  false,
			Package:    pkg,
			Statements: totalStatementCount,
			Uncovered:  totalStatementCount - totalCoveredCount,
		})
	}
	return
}

func getPackages(exclusionPatterns []string) (packages []string, err error) {
	var (
		out []byte
	)
	grepString := "go list ./... | grep -v /vendor/ "
	for _, exclusionPattern := range exclusionPatterns {
		grepString += fmt.Sprintf(" | grep -v %v", exclusionPattern)
	}
	cmd := exec.Command("sh", "-c", grepString)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		return
	}
	packages = strings.Split(string(out), "\n")
	return
}

func parseStatementCountFromRaw(rawDatum string) float64 {
	// raw: github.com/eltorocorp/cookiecrumbler-v2/src/shared/clock/api.go:6.34,6.72 1 1
	result, err := strconv.ParseFloat(strings.Split(rawDatum, " ")[1], 64)
	if err != nil {
		log.Println(err)
		panic("Error parsing raw data for statement count.")
	}
	return result
}

func parseCallCountFromRaw(rawDatum string) float64 {
	// raw: github.com/eltorocorp/cookiecrumbler-v2/src/shared/clock/api.go:6.34,6.72 1 1
	result, err := strconv.ParseFloat(strings.Split(rawDatum, " ")[2], 64)
	if err != nil {
		log.Println(err)
		panic("Error parsing raw data for call count.")
	}
	return result
}
