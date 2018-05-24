package interfaces

// CoverageToolingAPI describes methods that know how to retrieve data from
// go's native tooling.
type CoverageToolingAPI interface {
	GetRawCoverageAnalysisForPackage(pkg string) ([]string, error)
	GetPackages(exclusionPatterns []string) ([]string, error)
	PrintExcludedPackages(exclusionPattern string)
}
