package customanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ErrUsingOsExit which definition error of calling os.Exit in main func.
var ErrUsingOsExit = &analysis.Analyzer{
	Name: "ErrUsingOsExit",
	Doc:  "check of error wich call os.Exit in main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {

		// Go tree.
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.File:
				// Check packege main.
				if x.Name.Name == "main" {
					return true
				}
				return false // end of Inspect.

			case *ast.FuncDecl:
				// Check func of main.
				if x.Name.Name == "main" {
					return true
				}
				return false // end of Inspect

			case *ast.CallExpr:
				// Проверяем, вызван ли Exit из пакета os
				if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok {
						if ident.Name == "os" && sel.Sel.Name == "Exit" {

							pass.Reportf(x.Pos(), "func main doesn't include os.Exit()")
							return false // end of Inspect
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
