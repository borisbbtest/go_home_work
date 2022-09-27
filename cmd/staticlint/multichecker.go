package main

import (
	mmc "github.com/borisbbtest/go_home_work/interna/multichecker"
	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}
	analyzers = append(analyzers, printf.Analyzer)
	analyzers = append(analyzers, structtag.Analyzer)
	analyzers = append(analyzers, analyzer.Analyzer)
	analyzers = append(analyzers, analyzer.Analyzer)
	analyzers = append(analyzers, mmc.Analyzer)

	multichecker.Main(analyzers...)
}
