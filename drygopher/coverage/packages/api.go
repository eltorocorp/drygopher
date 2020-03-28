// Package packages exposes methods for listing packages URIs, filenames, etc...
package packages

import (
	"fmt"
	"go/build"
	"log"
	"path/filepath"
	"strings"

	"github.com/eltorocorp/drygopher/drygopher/coverage/hostiface"
)

// API contains methods that interact directly with the host environment.
type API struct {
	osioAPI hostiface.OSIOAPI
	execAPI hostiface.ExecAPI
}

// New returns a reference to an API
func New(execAPI hostiface.ExecAPI, osioAPI hostiface.OSIOAPI) *API {
	return &API{
		osioAPI: osioAPI,
		execAPI: execAPI,
	}
}

// GetPackages shells out a go list command that retusn a list of all packages
// below the current directory.
func (a *API) GetPackages(exclusionPatterns []string) (packages []string, err error) {
	var (
		out []byte
	)
	grepString := "go list ./..."
	for _, exclusionPattern := range exclusionPatterns {
		grepString += fmt.Sprintf(" | grep -v %v", exclusionPattern)
	}
	cmd := a.execAPI.Command("sh", "-c", grepString)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		return
	}
	packages = strings.Split(string(out), "\n")
	return
}

// GetFileNamesForPackage returns a list package URIs with associated filenames.
func (a *API) GetFileNamesForPackage(pkg string) ([]string, error) {
	fmt.Println("!!!", build.Default.GOPATH)
	gopath := a.osioAPI.GetGoPath()
	packagePath := gopath + "/src/" + pkg
	files, err := a.osioAPI.ReadDir(packagePath)
	if err != nil {
		return nil, err
	}
	fileNames := []string{}
	for _, file := range files {
		fileName := file.Name()
		if filepath.Ext(fileName) == ".go" {
			fullName := packagePath + "/" + fileName
			fileNames = append(fileNames, fullName)
		}
	}
	return fileNames, nil
}
