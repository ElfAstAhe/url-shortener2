package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer - iteration 20 task
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osExitChecker",
	Doc:  "check for os.Exit usage in package main and main function.",
	Run:  runOsExitUsage,
}

func runOsExitUsage(pass *analysis.Pass) (interface{}, error) {
	// We only care about the "main" package.
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// This function finds os.Exit calls only within the scope of the main function.
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != "main" {
				// If it's not the 'main' function declaration,
				// we don't need to descend into its children to check for os.Exit calls *in main*.
				return true // Continue inspecting siblings/other top-level nodes
			}

			// We found the 'main' function. Now inspect *only* its body for the call.
			ast.Inspect(funcDecl.Body, func(bodyN ast.Node) bool {
				if callExpr, isCall := bodyN.(*ast.CallExpr); isCall {
					checkOsExit(callExpr, pass)
				}
				return true // Keep looking inside the main function body
			})

			return false // Stop descending into the original main function node (already handled body)
		})
	}

	return nil, nil
}

// checkOsExit checks CallExpr for a call to os.Exit()
func checkOsExit(expr *ast.CallExpr, pass *analysis.Pass) {
	fun, ok := expr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	x, ok := fun.X.(*ast.Ident)
	if !ok || x.Name != "os" || fun.Sel.Name != "Exit" {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:     expr.Pos(),
		Message: "Using call os.Exit() into main function.",
	})
}
