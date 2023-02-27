//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package recv

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/fogfish/gurl/v2/http"
)

func StatusCode(code string) ast.Expr {
	return &ast.SelectorExpr{
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("ƒ"),
			Sel: ast.NewIdent("Status"),
		},
		Sel: ast.NewIdent(code),
	}
}

func Header[T http.MatchableHeaderValues](header string, value T) ast.Expr {
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
			X:   ast.NewIdent("ƒ"),
			Sel: ast.NewIdent("Header"),
		},
		Args: args,
	}
}

func Match(value string) ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("ƒ"),
			Sel: ast.NewIdent("Match"),
		},
		Args: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(value),
			},
		},
	}
}
