package interfaces

// HostAPI describes methods that need to interact directly with the host system.
// or are otherwise side-effecting.
type HostAPI interface {
	GetRawCoverageAnalysisForPackage(pkg string) ([]string, error)
	GetPackages(exclusionPatterns []string) ([]string, error)
	PrintExcludedPackages(exclusionPattern string)
	GetFileNamesForPackage(pkg string) ([]string, error)
	SaveCoverageProfile(fileName string, rawData []string) error
}
