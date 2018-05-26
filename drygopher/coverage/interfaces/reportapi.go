package interfaces

import "github.com/eltorocorp/drygopher/drygopher/coverage/pckg"

// ReportAPI represents anything that knows how to work with coverage reports.
type ReportAPI interface {
	PrintExcludedPackages(exclusionPattern string)
	OutputCoverageReport(allPackages pckg.Group, exclusionPatterns []string)
}
