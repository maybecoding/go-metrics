package staticlint

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/simple"
)

func GetSimpleAnalyzers() []*analysis.Analyzer {
	as := make([]*analysis.Analyzer, 0, len(simple.Analyzers))

	for _, a := range simple.Analyzers {
		as = append(as, a.Analyzer)
	}
	return as
}
