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

// Sandbox environment for Golang
type Sandbox struct {
	stdout io.Writer
	Path   string
	Cache  string
}

// Configure new sandbox environment
func NewSandbox(stdout io.Writer, root string) (*Sandbox, error) {
	if root == "" {
		tmp, err := os.MkdirTemp(os.TempDir(), org)
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
		stdout: stdout,
		Path:   root,
		Cache:  filepath.Join(root, "cache"),
	}, nil
}

// Compile package, producing binary code
func (box *Sandbox) Compile(pkg *Package) error {
	gcc := exec.Command("go", "build", "-mod=mod")
	gcc.Dir = pkg.SourceCode
	gcc.Env = []string{
		"GO111MODULE=on",
		"GOPATH=" + box.Path,
		"GOCACHE=" + box.Cache,
	}
	gcc.Stderr, gcc.Stdout = box.stdout, box.stdout

	err := gcc.Run()
	if err != nil {
		return err
	}

	return nil
}
