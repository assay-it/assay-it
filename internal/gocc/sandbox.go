//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package gocc

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Sandbox struct {
	Path  string
	Cache string
}

func NewSandbox(root string) (*Sandbox, error) {
	if root == "" {
		tmp, err := os.MkdirTemp(os.TempDir(), "assay")
		if err != nil {
			return nil, err
		}
		root = tmp
	}

	err := os.MkdirAll(root, 0755)
	if err != nil {
		return nil, err
	}

	return &Sandbox{
		Path:  root,
		Cache: filepath.Join(root, "cache"),
	}, nil
}

func (box *Sandbox) Compile(stdout io.Writer, pkg *Package) error {
	gcc := exec.Command("go", "build", "-mod=mod")
	gcc.Dir = pkg.SourceCode
	gcc.Env = []string{
		"GO111MODULE=on",
		"GOPATH=" + box.Path,
		"GOCACHE=" + box.Cache,
	}
	gcc.Stderr, gcc.Stdout = stdout, stdout

	err := gcc.Run()
	if err != nil {
		return err
	}

	return nil
}
