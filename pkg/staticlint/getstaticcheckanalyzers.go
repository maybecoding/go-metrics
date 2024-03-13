package staticlint

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

func GetStaticCheckAnalyzers() []*analysis.Analyzer {
	as := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))

	for _, a := range staticcheck.Analyzers {
		as = append(as, a.Analyzer)
	}
	return as
}
