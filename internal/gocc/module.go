package gocc

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"
)

type Module struct {
	Name       string
	SourceCode string
}

func NewModule(root string, pkg string) (*Module, error) {
	return &Module{
		Name:       pkg,
		SourceCode: root,
	}, nil
}

func (mod *Module) CopyFrom(file string) (*Suite, error) {
	switch {
	case filepath.Ext(file) == ".go":
		return mod.copyFromGo(file)
	default:
		return nil, fmt.Errorf("not supported %s", file)
	}
}

func (mod *Module) copyFromGo(file string) (*Suite, error) {
	return NewSuite(file)
}

func (mod *Module) CreateRunner(tests []*Suite) error {
	run := filepath.Join(mod.SourceCode, "runner.go")

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

	for _, test := range tests {
		if test.Pkg != "main" {
			astutil.AddNamedImport(fset, code, test.Pkg, test.PkgPath)
		}
	}

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

func (mod *Module) Run(stdout io.Writer, defaultHost string) error {
	run := exec.Command("go", "run", "runner.go", defaultHost)
	run.Dir = mod.SourceCode
	run.Stderr, run.Stdout = stdout, stdout

	err := run.Run()
	if err != nil {
		return err
	}

	return nil
}
