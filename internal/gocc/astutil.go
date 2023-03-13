package gocc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
)

type Suite struct {
	File    string
	Pkg     string
	PkgPath string
	Units   []string
}

func NewSuite(module, file string) (*Suite, error) {
	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &Suite{
		File:    file,
		Pkg:     code.Name.Name,
		PkgPath: filepath.Dir(filepath.Join(module, file)),
		Units:   astLookupUnitTest(code),
	}, nil
}

//
// Parser Golang source
//

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

//
// Modify golang source
//

func astReWrite(node *ast.File, tests []*Suite) {
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

func astMain(tests []*Suite) *ast.BlockStmt {
	args := []ast.Expr{
		&ast.SelectorExpr{
			X:   ast.NewIdent("os"),
			Sel: ast.NewIdent("Stdout"),
		},
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{X: ast.NewIdent("http"), Sel: ast.NewIdent("New")},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun:  &ast.SelectorExpr{X: ast.NewIdent("http"), Sel: ast.NewIdent("WithMemento")},
					Args: []ast.Expr{},
				},
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{X: ast.NewIdent("http"), Sel: ast.NewIdent("WithDefaultHost")},
					Args: []ast.Expr{
						&ast.IndexExpr{
							X: &ast.SelectorExpr{X: ast.NewIdent("os"), Sel: ast.NewIdent("Args")},
							Index: &ast.BasicLit{
								Kind:  token.INT,
								Value: "1",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		for _, unit := range test.Units {
			if test.Pkg == "main" {
				args = append(args, ast.NewIdent(unit))
			} else {
				args = append(args, &ast.SelectorExpr{X: ast.NewIdent(test.Pkg), Sel: ast.NewIdent(unit)})
			}
		}
	}

	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  &ast.SelectorExpr{X: ast.NewIdent("http"), Sel: ast.NewIdent("WriteOnce")},
					Args: args,
				},
			},
		},
	}
}
