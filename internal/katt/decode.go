//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package katt

import (
	"bufio"
	"go/ast"
	"log"
	"os"
	"strings"

	"github.com/assay-it/assay/internal/ast/http"
)

type decoder interface {
	parse(*katt, string) error
	flush(*katt, *http.Mod) error
}

type katt struct {
	state decoder

	unittest string
	method   func(args []ast.Expr) ast.Expr
	expr     []ast.Expr
}

func newKatt() *katt {
	return &katt{
		state: decoderIdle{},
		expr:  []ast.Expr{},
	}
}

func (spec *katt) parse(line string) error {
	return spec.state.parse(spec, line)
}

func (spec *katt) flush(mod *http.Mod) error {
	return spec.state.flush(spec, mod)
}

type decoderIdle struct{}

func (decoderIdle) parse(spec *katt, line string) error {
	if !parsableUnitTestID(line) {
		return nil
	}

	id, err := parseUnitTestID(line)
	if err != nil {
		return err
	}

	*spec = katt{
		state:    decoderMethod{},
		unittest: *id,
		expr:     []ast.Expr{},
	}

	return nil
}

func (decoderIdle) flush(spec *katt, mod *http.Mod) error {
	return nil
}

type decoderMethod struct{}

func (decoderMethod) parse(spec *katt, line string) error {
	if !parsableRequestMethod(line) {
		return nil
	}

	method, expr, err := parseRequestMethod(line)
	if err != nil {
		return err
	}

	spec.method = method
	spec.expr = append(spec.expr, expr)
	spec.state = decoderRequestHeader{}

	return nil
}

func (decoderMethod) flush(spec *katt, mod *http.Mod) error {
	return nil
}

type decoderRequestHeader struct{}

func (decoderRequestHeader) parse(spec *katt, line string) error {
	if parsableStatusCode(line) {
		spec.state = decoderStatusCode{}
		return spec.parse(line)
	}

	if !parsableRequestHeader(line) {
		return nil
	}

	expr, err := parseRequestHeader(line)
	if err != nil {
		return err
	}

	spec.expr = append(spec.expr, expr)
	return nil
}

func (decoderRequestHeader) flush(spec *katt, mod *http.Mod) error {
	return nil
}

type decoderStatusCode struct{}

func (decoderStatusCode) parse(spec *katt, line string) error {
	if !parsableStatusCode(line) {
		return nil
	}

	expr, err := parseStatusCode(line)
	if err != nil {
		return err
	}

	spec.expr = append(spec.expr, expr)
	spec.state = decoderResponseHeader{}
	return nil
}

func (decoderStatusCode) flush(spec *katt, mod *http.Mod) error {
	mod.Add(spec.unittest, spec.method(spec.expr))
	return nil
}

type decoderResponseHeader struct{}

func (decoderResponseHeader) parse(spec *katt, line string) error {
	if parsablePayload(line) {
		spec.state = &decoderResponsePayload{}
		return spec.parse(line)
	}

	if !parsableResponseHeader(line) {
		return nil
	}

	expr, err := parseResponseHeader(line)
	if err != nil {
		return err
	}

	spec.expr = append(spec.expr, expr)
	return nil
}

func (decoderResponseHeader) flush(spec *katt, mod *http.Mod) error {
	mod.Add(spec.unittest, spec.method(spec.expr))
	return nil
}

type decoderResponsePayload struct {
	buf strings.Builder
}

func (c *decoderResponsePayload) parse(spec *katt, line string) error {
	c.buf.WriteString(line)

	if parsablePayloadEOF(line) {
		expr, err := parsePayload(c.buf.String())
		if err != nil {
			return err
		}
		spec.expr = append(spec.expr, expr)
		spec.state = decoderResponseHeader{}
	}

	return nil
}

func (*decoderResponsePayload) flush(spec *katt, mod *http.Mod) error {
	return nil
}

// Decodes KATT file into Golang module
func Decode(filename string, mod *http.Mod) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := newKatt()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if parsableUnitTestID(line) {
			if err := decoder.flush(mod); err != nil {
				return err
			}
		}

		err := decoder.parse(line)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := decoder.flush(mod); err != nil {
		return err
	}

	return nil
}
