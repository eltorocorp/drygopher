package interfaces

import "github.com/eltorocorp/drygopher/internal/pckg"

// RawAPI is something that knows how to work with raw go coverage analysis data.
type RawAPI interface {
	GetRawCoverageAnalysisForPackage(pkg string) ([]string, error)
	AggregateRawPackageAnalysisData(pkg string, rawPkgCoverageData []string) (*pckg.Stats, error)
}
