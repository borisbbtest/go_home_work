package main

import (
	"go/ast"

	"github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

const Doc = `check using os.Exit() in main package`

var Analyzer = &analysis.Analyzer{
	Name: "osmainchecker",
	Doc:  Doc,
	Run:  run,
}

func hasOSExit(f *ast.FuncDecl, pass *analysis.Pass) {
	ast.Inspect(f, func(n ast.Node) bool {
		funcCall, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if selExpr, ok := funcCall.Fun.(*ast.SelectorExpr); ok {
			if i, ok := selExpr.X.(*ast.Ident); ok {
				if i.Name == "os" && selExpr.Sel.Name == "Exit" {
					pass.Reportf(funcCall.Pos(), "found os.Exit() call in main()")
				}
			}
		}
		return true
	})
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if funcDecl.Name.Name == "main" {
				hasOSExit(funcDecl, pass)
			}
			return true
		})
	}
	return nil, nil
}

func main() {
	var analyzers []*analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}
	analyzers = append(analyzers, printf.Analyzer)
	analyzers = append(analyzers, structtag.Analyzer)
	analyzers = append(analyzers, analyzer.Analyzer)
	analyzers = append(analyzers, analyzer.Analyzer)
	analyzers = append(analyzers, Analyzer)

	multichecker.Main(analyzers...)
}
