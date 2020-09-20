//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/assay-it/sdk-go/assay"
	"github.com/assay-it/sdk-go/http"
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
	endpoint string
	digest   string
	debug    bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&endpoint, "api", "api.assay.it", "assay's rest api endpoint")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output of http traffic with assay.it endpoint")
	rootCmd.PersistentFlags().StringVarP(&digest, "key", "t", "", "personal access key to authenticate with assay.it")
}

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

//
// requiredFlagKey check presence of secret key, use it with PreRunE
func requiredFlagKey(cmd *cobra.Command, args []string) error {
	if digest == "" {
		return errors.New("undefined secret key, obtain a new personal access key from assay.it and use it with --key flag")
	}
	return nil
}

//
// eval executed side-effect on http computation
func eval(f assay.Arrow) error {
	deb := assay.LogLevelNone
	if debug {
		deb = assay.LogLevelDebug
	}

	io := http.DefaultIO(assay.Logging(deb))
	return f(io).Fail
}
