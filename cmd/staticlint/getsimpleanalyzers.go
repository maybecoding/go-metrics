package main

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/simple"
)

func getSsimpleAnalyzers() []*analysis.Analyzer {
	as := make([]*analysis.Analyzer, 0, len(simple.Analyzers))

	for _, a := range simple.Analyzers {
		as = append(as, a.Analyzer)
		//fmt.Println(a.Analyzer.Name, strings.Split(a.Analyzer.Doc, "\n")[0])
	}
	return as
}
