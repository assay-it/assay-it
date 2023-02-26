//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package send

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/fogfish/gurl/v2/http"
)

// Creates ø.URI node
func URI(url string, opts ...any) ast.Expr {
	args := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(url),
		},
	}

	// TODO: fill args by opts
	// if len(opts) != 0 {
	// }
	// ast.Comment

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("ø"),
			Sel: ast.NewIdent("URI"),
		},
		Args: args,
	}
}

func Header[T http.ReadableHeaderValues](header string, value T) ast.Expr {
	args := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(header),
		},
	}

	switch v := any(value).(type) {
	case string:
		args = append(args,
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(v),
			},
		)
	case int:
		args = append(args,
			&ast.BasicLit{
				Kind:  token.INT,
				Value: strconv.Itoa(v),
			},
		)
	}

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("ø"),
			Sel: ast.NewIdent("Header"),
		},
		Args: args,
	}
}
