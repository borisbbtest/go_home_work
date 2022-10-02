package main

import (
	mainmultichecker "github.com/borisbbtest/go_home_work/internal/multichecker"
	"github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var m []*analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		m = append(m, a.Analyzer)
	}
	m = append(m, printf.Analyzer)
	m = append(m, shadow.Analyzer)
	m = append(m, shift.Analyzer)
	m = append(m, structtag.Analyzer)
	m = append(m, analyzer.Analyzer)
	m = append(m, mainmultichecker.Analyzer)
	multichecker.Main(m...)
}
