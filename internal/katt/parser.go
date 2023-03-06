//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package katt

import (
	"fmt"
	"go/ast"
	"net/http"
	"regexp"
	"strings"

	ghttp "github.com/assay-it/assay-it/internal/ast/http"
	"github.com/assay-it/assay-it/internal/ast/http/recv"
	"github.com/assay-it/assay-it/internal/ast/http/send"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func parsableUnitTestID(spec string) bool {
	return strings.HasPrefix(spec, "##")
}

func parseUnitTestID(spec string) (*string, error) {
	if !parsableUnitTestID(spec) {
		return nil, fmt.Errorf("invalid unit test id: %s", spec)
	}

	nowhitespace := strings.Trim(spec[2:], " ")
	uppercase := cases.Title(language.AmericanEnglish).String(nowhitespace)
	unittestid := strings.ReplaceAll(uppercase, " ", "")

	return &unittestid, nil
}

func parsableRequestMethod(spec string) bool {
	return strings.HasPrefix(spec, http.MethodGet)
}

func parseRequestMethod(spec string) (func(args []ast.Expr) ast.Expr, ast.Expr, error) {
	if !parsableRequestMethod(spec) {
		return nil, nil, fmt.Errorf("invalid request: %s", spec)
	}

	nowhitespace := strings.Trim(spec, " ")
	seq := strings.Split(nowhitespace, " ")
	if len(seq) != 2 {
		return nil, nil, fmt.Errorf("invalid request: %s", spec)
	}
	mthd := ghttp.GET
	expr := send.URI(seq[1])

	return mthd, expr, nil
}

func parsableRequestHeader(spec string) bool {
	return strings.HasPrefix(spec, "> ")
}

func parseRequestHeader(spec string) (ast.Expr, error) {
	if !parsableRequestHeader(spec) {
		return nil, fmt.Errorf("invalid header spec: %s", spec)
	}

	nowhitespace := strings.Trim(spec[2:], " ")
	seq := strings.SplitN(nowhitespace, ": ", 2)
	if len(seq) != 2 {
		return nil, fmt.Errorf("invalid header spec: %s", spec)
	}
	expr := send.Header(seq[0], seq[1])

	return expr, nil
}

func parsableStatusCode(spec string) bool {
	is, err := regexp.MatchString("< [0-9]{3}.*", spec)
	return err == nil && is
}

func parseStatusCode(spec string) (ast.Expr, error) {
	if !parsableStatusCode(spec) {
		return nil, fmt.Errorf("invalid status code: %s", spec)
	}

	nowhitespace := strings.Trim(spec[2:], " ")
	seq := strings.SplitN(nowhitespace, " ", 2)
	if len(seq) != 2 {
		return nil, fmt.Errorf("invalid status code: %s", spec)
	}
	expr := recv.StatusCode(seq[1])

	return expr, nil
}

func parsableResponseHeader(spec string) bool {
	return strings.HasPrefix(spec, "< ")
}

func parseResponseHeader(spec string) (ast.Expr, error) {
	if !parsableResponseHeader(spec) {
		return nil, fmt.Errorf("invalid header spec: %s", spec)
	}

	nowhitespace := strings.Trim(spec[2:], " ")
	seq := strings.SplitN(nowhitespace, ": ", 2)
	if len(seq) != 2 {
		return nil, fmt.Errorf("invalid header spec: %s", spec)
	}
	expr := recv.Header(seq[0], seq[1])

	return expr, nil
}

func parsablePayload(spec string) bool {
	return strings.HasPrefix(spec, "{")
}

func parsablePayloadEOF(spec string) bool {
	return strings.HasPrefix(spec, "}")
}

func parsePayload(spec string) (ast.Expr, error) {
	expr := recv.Match(spec)
	return expr, nil
}
