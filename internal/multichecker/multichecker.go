package multichecker

import (
	"fmt"
	"go/ast"

	"github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

const Doc = `check using os.Exit() in main package`

var Analyzer = &analysis.Analyzer{
	Name: "MulticheckerBeginner",
	Doc:  Doc,
	Run:  run,
}

func OSExit(f *ast.FuncDecl, pass *analysis.Pass) {
	ast.Inspect(f, func(n ast.Node) bool {
		funcCall, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if selExpr, ok := funcCall.Fun.(*ast.SelectorExpr); ok {
			if i, ok := selExpr.X.(*ast.Ident); ok {
				if i.Name == "os" && selExpr.Sel.Name == "Exit" {
					fmt.Print("test")
					pass.Reportf(funcCall.Pos(), "found os.Exit method in main function")
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
				OSExit(funcDecl, pass)
			}
			return true
		})
	}
	return nil, nil
}

type Maninmultichecker struct {
	analyzers []*analysis.Analyzer
}

func InitMultichecker() (m *Maninmultichecker) {

	m = &Maninmultichecker{}
	for _, a := range staticcheck.Analyzers {
		m.analyzers = append(m.analyzers, a.Analyzer)
	}
	m.analyzers = append(m.analyzers, printf.Analyzer)
	m.analyzers = append(m.analyzers, shadow.Analyzer)
	m.analyzers = append(m.analyzers, shift.Analyzer)
	m.analyzers = append(m.analyzers, structtag.Analyzer)
	m.analyzers = append(m.analyzers, analyzer.Analyzer)
	m.analyzers = append(m.analyzers, Analyzer)
	return m
}

func (m *Maninmultichecker) Start() {
	multichecker.Main(m.analyzers...)
}
