package shelledcmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

// API contains methods that send commands to the shell.
type API struct{}

// GetRawCoverageAnalysisForPackage shells out a go test command and returns the
// raw data from the resulting coverage profile.
func (a *API) GetRawCoverageAnalysisForPackage(pkg string) ([]string, error) {
	covermode := "count"
	analyzeCmdText := "go test -covermode=%v -coverprofile=tmp.out %v"
	analyzeCmdText = fmt.Sprintf(analyzeCmdText, covermode, pkg)
	analyzeCoverageCmd := exec.Command("sh", "-c", analyzeCmdText)
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
	tmpOut, err = ioutil.ReadFile("tmp.out")
	if err != nil {
		return nil, err
	}
	rawPkgCoverageData := strings.Split(string(tmpOut), "\n")
	return rawPkgCoverageData, nil
}

// GetPackages shells out a go list command that retusn a list of all packages
// below the current directory.
func (a *API) GetPackages(exclusionPatterns []string) (packages []string, err error) {
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

// PrintExcludedPackages shells out a go list command and sends the results
// of the command directly to stdout.
func (a *API) PrintExcludedPackages(exclusionPattern string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("go list ./... | grep -v /vendor/ | grep %v", exclusionPattern))
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
