package main

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/maybecoding/go-metrics.git/cmd/staticlint/mainosexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	var analyzers []*analysis.Analyzer
	// Add all analyzers from golang.org/x/tools/go/analysis/passes
	analyzers = append(analyzers, getPassesAnalyzers()...)
	// Add static check analyzers from class SA
	analyzers = append(analyzers, getStaticCheckAnalyzers()...)
	// Add check analyzers for code simplicity
	analyzers = append(analyzers, getSsimpleAnalyzers()...)

	//Add third party analyzers and identifying pkg:main dcl:main os.Exit
	analyzers = append(analyzers, errcheck.Analyzer, ineffassign.Analyzer, mainosexit.Analyzer)

	multichecker.Main(analyzers...)
}
