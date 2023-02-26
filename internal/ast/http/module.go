//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package http

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"strconv"
)

type Mod struct {
	path string
	file *ast.File
}

func New(mod string, path string) *Mod {
	return &Mod{
		path: path,
		file: &ast.File{
			Name:  ast.NewIdent(mod),
			Decls: []ast.Decl{imports()},
		},
	}
}

func imports() ast.Decl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote("github.com/fogfish/gurl/v2/http"),
				},
			},
			&ast.ImportSpec{
				Name: ast.NewIdent("ø"),
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote("github.com/fogfish/gurl/v2/http/send"),
				},
			},
			&ast.ImportSpec{
				Name: ast.NewIdent("ƒ"),
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote("github.com/fogfish/gurl/v2/http/recv"),
				},
			},
		},
	}
}

func (mod *Mod) Add(name string, expr ast.Expr) {
	fun := &ast.FuncDecl{
		Name: ast.NewIdent(name),
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: []*ast.Field{}},
			Results: &ast.FieldList{List: []*ast.Field{
				{
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("http"),
						Sel: ast.NewIdent("Arrow"),
					},
				},
			}},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{expr},
				},
			},
		},
	}
	mod.file.Decls = append(mod.file.Decls, fun)
}

func (mod *Mod) Write() error {
	fset := token.NewFileSet()

	var buf bytes.Buffer
	err := format.Node(&buf, fset, mod.file)
	if err != nil {
		log.Fatal(err)
	}

	fd, _ := os.Create(mod.path)
	defer fd.Close()
	_, err = fd.Write(buf.Bytes())

	return err
}
