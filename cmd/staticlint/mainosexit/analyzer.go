// Package mainosexit analyzes if os.Exit called in main function of main package
package mainosexit

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "mainosexit",
	Doc:  "checks for usage os.Exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	findMainDcl := func(f *ast.File) *ast.FuncDecl {
		mainFound := false
		var dcl *ast.FuncDecl
		// using ast.Inspect we can find main func declaration
		ast.Inspect(f, func(node ast.Node) bool {
			// find
			if dlr, ok := node.(*ast.FuncDecl); ok {
				if dlr.Name.Name == "main" {
					mainFound = true
					dcl = dlr
				}
			}
			// if found stop walking
			return !mainFound
		})
		return dcl
	}

	// identify if call of os.Exit
	isOsExit := func(x *ast.CallExpr) (bool, token.Pos) {
		se, ok := x.Fun.(*ast.SelectorExpr)
		if !ok {
			return false, 0
		}
		if se.Sel.Name != "Exit" {
			return false, 0
		}
		pkg, ok := se.X.(*ast.Ident)
		if !ok || pkg.Name != "os" {
			return false, 0
		}
		return true, se.Sel.NamePos
	}

	// for every file
	for _, file := range pass.Files {
		if isGenerated(file) {
			continue
		}
		// with package main
		if file.Name.Name != "main" {
			continue
		}

		// find function main
		mainDcl := findMainDcl(file)
		if mainDcl == nil {
			continue
		}

		// find os exit
		ast.Inspect(mainDcl, func(node ast.Node) bool {
			if ce, ok := node.(*ast.CallExpr); ok {
				if osEx, pos := isOsExit(ce); osEx {
					pass.Reportf(pos, "call of os.Exit in func main of package main")
				}
			}
			return true
		})
	}
	return nil, nil
}

func isGenerated(file *ast.File) bool {
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// Code generated ") && strings.HasSuffix(c.Text, " DO NOT EDIT.") {
				return true
			}
		}
	}

	return false
}
