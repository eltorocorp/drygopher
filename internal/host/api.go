package host

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// API contains methods that interact directly with the host environment.
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
	rawPkgCoverageData = rawPkgCoverageData[1:]
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

// GetFileNamesForPackage returns a list package URIs with associated filenames.
func (a *API) GetFileNamesForPackage(pkg string) ([]string, error) {
	var gopath string
	var set bool
	if gopath, set = os.LookupEnv("GOPATH"); !set {
		return nil, errors.New("GOPATH not set")
	}
	packagePath := gopath + "/src/" + pkg
	files, err := ioutil.ReadDir(packagePath)
	if err != nil {
		return nil, err
	}
	fileNames := []string{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".go" {
			fullName := pkg + "/" + file.Name()
			fileNames = append(fileNames, fullName)
		}
	}
	return fileNames, nil
}

// SaveCoverageProfile saves the supplied raw data to the desired file.
func (a *API) SaveCoverageProfile(fileName string, rawData []string) error {
	for i, r := range rawData {
		if strings.TrimSpace(r) == "" {
			rawData = append(rawData[:i], rawData[i+1:]...)
		}
	}
	sort.StringSlice(rawData).Sort()
	profile := "mode: count\n" + strings.Join(rawData, "\n")
	return ioutil.WriteFile(fileName, []byte(profile), os.ModePerm)
}
