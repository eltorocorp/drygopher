package interfaces

// PackageAPI describes anything that knows how to interact with packages.
type PackageAPI interface {
	GetPackages(exclusionPatterns []string) ([]string, error)
	GetFileNamesForPackage(pkg string) ([]string, error)
}
