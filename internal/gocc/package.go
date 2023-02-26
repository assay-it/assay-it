//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package gocc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/assay-it/assay/internal/ast/http"
	"github.com/assay-it/assay/internal/katt"
)

const (
	hub = "github.com"
	org = "assay"
)

type Package struct {
	Name       string
	SourceCode string
}

func NewPackage(root string, pkg string) (*Package, error) {
	src := filepath.Join(root, "src", hub, org, pkg)
	err := os.MkdirAll(src, 0755)
	if err != nil {
		return nil, err
	}

	return &Package{
		Name:       pkg,
		SourceCode: src,
	}, nil
}

func (pkg *Package) CopyFrom(file string) ([]string, error) {
	switch {
	case filepath.Ext(file) == ".go":
		return pkg.copyFromGo(file)
	case filepath.Ext(file) == ".md":
		return pkg.copyFromMd(file)
	default:
		return nil, fmt.Errorf("Not supported %s", file)
	}
}

func (pkg *Package) copyFromGo(file string) ([]string, error) {
	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	code.Name = ast.NewIdent("main")

	f, err := os.Create(filepath.Join(pkg.SourceCode, filepath.Base(file)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = printer.Fprint(f, fset, code)
	if err != nil {
		return nil, err
	}

	return astLookupUnitTest(code), nil
}

func (pkg *Package) copyFromMd(file string) ([]string, error) {
	gf := filepath.Join(pkg.SourceCode, strings.ReplaceAll(filepath.Base(file), ".md", ".go"))
	cc := http.New("main", gf)
	if err := katt.Decode(file, cc); err != nil {
		return nil, err
	}

	if err := cc.Write(); err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, gf, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return astLookupUnitTest(code), nil
}

func astLookupUnitTest(node *ast.File) (units []string) {
	ast.Inspect(node, func(n ast.Node) bool {
		if f := astUnitTestFunc(n); f != nil {
			units = append(units, f.Name.Name)
		}
		return true
	})
	return
}

func astUnitTestFunc(node ast.Node) *ast.FuncDecl {
	switch f := node.(type) {
	case *ast.FuncDecl:
		if !unicode.IsUpper([]rune(f.Name.Name)[0]) || f.Type.Results == nil || len(f.Type.Params.List) != 0 || len(f.Type.Results.List) != 1 {
			return nil
		}
		if !strings.HasPrefix(f.Name.Name, "Test") {
			return nil
		}

		switch t := f.Type.Results.List[0].Type.(type) {
		case *ast.SelectorExpr:
			switch l := t.X.(type) {
			case *ast.Ident:
				if l.Name == "http" && t.Sel.Name == "Arrow" {
					return f
				}
			}
		}
	}
	return nil
}

func (pkg *Package) CreateRunner(tests []string) error {
	run := filepath.Join(pkg.SourceCode, "runner.go")

	err := os.WriteFile(run, []byte(`
	package main

	import (
		"os"
		"github.com/fogfish/gurl/v2/http"
	)
	
	func main(){}
	`), 0644)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, run, nil, parser.ParseComments)
	astReWrite(code, tests)

	f, err := os.Create(run)
	if err != nil {
		return err
	}
	defer f.Close()

	err = printer.Fprint(f, fset, code)
	if err != nil {
		return err
	}

	return nil
}

func astReWrite(node *ast.File, tests []string) {
	ast.Inspect(node, func(node ast.Node) bool {
		switch f := node.(type) {
		case *ast.FuncDecl:
			if f.Name.Name == "main" {
				f.Body = astMain(tests)
			}
		}
		return true
	})
}

func astMain(tests []string) *ast.BlockStmt {
	args := []ast.Expr{
		&ast.SelectorExpr{
			X:   &ast.Ident{Name: "os"},
			Sel: &ast.Ident{Name: "Stdout"},
		},
	}

	for _, test := range tests {
		args = append(args, ast.NewIdent(test))
	}

	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  &ast.SelectorExpr{X: &ast.Ident{Name: "http"}, Sel: &ast.Ident{Name: "WriteOnce"}},
					Args: args,
				},
			},
		},
	}
}

func (pkg *Package) CreateMod() error {
	run := filepath.Join(pkg.SourceCode, "go.mod")

	err := os.WriteFile(run, []byte("module "+pkg.Name), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (pkg *Package) Run(stdout io.Writer) error {
	run := exec.Command("./" + pkg.Name)
	run.Dir = pkg.SourceCode
	run.Stderr, run.Stdout = stdout, stdout

	err := run.Run()
	if err != nil {
		return err
	}

	return nil
}
