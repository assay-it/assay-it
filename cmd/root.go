//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/assay-it/assay/internal/printer"
	"github.com/spf13/cobra"
)

// Execute is entry point for cobra cli application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var (
	stdout = printer.New(os.Stdout)
	stderr = printer.New(os.Stderr)
)

var rootCmd = &cobra.Command{
	Use:     "assay",
	Short:   "command line interface to https://assay.it",
	Long:    `command line interface to https://assay.it`,
	Run:     root,
	Version: "v0",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}
