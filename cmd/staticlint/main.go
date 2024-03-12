package main

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/maybecoding/go-metrics.git/pkg/staticlint"
	"github.com/maybecoding/go-metrics.git/pkg/staticlint/mainosexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	var analyzers []*analysis.Analyzer
	// Add all analyzers from golang.org/x/tools/go/analysis/passes
	analyzers = append(analyzers, staticlint.GetPassesAnalyzers()...)
	// Add static check analyzers from class SA
	analyzers = append(analyzers, staticlint.GetStaticCheckAnalyzers()...)
	// Add check analyzers for code simplicity
	analyzers = append(analyzers, staticlint.GetSimpleAnalyzers()...)

	// Add third party analyzers and identifying pkg:main dcl:main os.Exit
	analyzers = append(analyzers, errcheck.Analyzer, ineffassign.Analyzer, mainosexit.Analyzer)

	multichecker.Main(analyzers...)
}
