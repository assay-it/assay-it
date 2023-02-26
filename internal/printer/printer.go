//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

// Printer to console utility
type Printer struct {
	stdout io.Writer
}

func New(stdout io.Writer) *Printer {
	return &Printer{
		stdout: stdout,
	}
}

func (printer *Printer) Success(s string, args ...any) error {
	printer.stdout.Write([]byte(color.GreenString(s, args...)))
	return nil
}

func (printer *Printer) Error(s string, args ...any) error {
	printer.stdout.Write([]byte(color.RedString(s, args...)))
	return nil
}

func (printer *Printer) Warning(s string, args ...any) error {
	printer.stdout.Write([]byte(color.YellowString(s, args...)))
	return nil
}

func (printer *Printer) Notice(s string, args ...any) error {
	printer.stdout.Write([]byte(color.WhiteString(s, args...)))
	return nil
}

func (printer *Printer) Info(s string, args ...any) error {
	printer.stdout.Write([]byte(fmt.Sprintf(s, args...)))
	return nil
}

func (printer *Printer) FormattedJSON(data string) error {
	var obj map[string]any
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return err
	}

	f := colorjson.NewFormatter()
	f.Indent = 2
	f.KeyColor = color.New(color.FgBlue)

	encoded, err := f.Marshal(obj)
	if err != nil {
		return err
	}

	escaped := "| " + strings.ReplaceAll(string(encoded), "\n", "\n| ")
	_, err = printer.stdout.Write([]byte(escaped))
	if err != nil {
		return err
	}

	_, err = printer.stdout.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}
