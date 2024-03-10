package main

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

func getStaticCheckAnalyzers() []*analysis.Analyzer {
	as := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))

	for _, a := range staticcheck.Analyzers {
		as = append(as, a.Analyzer)
		//fmt.Println(a.Analyzer.Name, strings.Split(a.Analyzer.Doc, "\n")[0])
	}
	return as
}
