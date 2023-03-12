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

	"github.com/assay-it/assay-it/internal/ast/http"
	"github.com/assay-it/assay-it/internal/katt"
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

func (pkg *Package) CopyFrom(file string) (*Suite, error) {
	switch {
	case filepath.Ext(file) == ".go":
		return pkg.copyFromGo(file)
	case filepath.Ext(file) == ".md":
		return pkg.copyFromMd(file)
	default:
		return nil, fmt.Errorf("not supported %s", file)
	}
}

func (pkg *Package) copyFromGo(file string) (*Suite, error) {
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

	return &Suite{
		File:  file,
		Pkg:   "main",
		Units: astLookupUnitTest(code),
	}, nil
}

func (pkg *Package) copyFromMd(file string) (*Suite, error) {
	gofile := filepath.Join(pkg.SourceCode, strings.ReplaceAll(filepath.Base(file), ".md", ".go"))
	cc := http.New("main", gofile)
	if err := katt.Decode(file, cc); err != nil {
		return nil, err
	}

	if err := cc.Write(); err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &Suite{
		File:  file,
		Pkg:   "main",
		Units: astLookupUnitTest(code),
	}, nil
}

func (pkg *Package) CreateRunner(tests []*Suite) error {
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
	if err != nil {
		return err
	}

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

func (pkg *Package) CreateMod() error {
	run := filepath.Join(pkg.SourceCode, "go.mod")

	err := os.WriteFile(run, []byte("module "+pkg.Name), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (pkg *Package) Run(stdout io.Writer, defaultHost string) error {
	run := exec.Command("./"+pkg.Name, defaultHost)
	run.Dir = pkg.SourceCode
	run.Stderr, run.Stdout = stdout, stdout

	err := run.Run()
	if err != nil {
		return err
	}

	return nil
}
