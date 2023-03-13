//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/assay-it/assay-it/internal/printer"
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
	Use:   "assay-it",
	Short: "Confirm Quality and Eliminate Risk by Testing Microservices in Production.",
	Long: `
Confirm Quality and Eliminate Risk by Testing Microservices in Production.
Runs testing of services across environment and deployments.
	`,
	Run:     root,
	Version: "v1.1.0",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// helper function outputs return of testing
// built over JSON protocol github.com/fogfish/gurl/v2/http.WriteOnce
func stdoutTestResults(silent, verbose bool, data []byte) error {
	type Status struct {
		ID       string        `json:"id"`
		Status   string        `json:"status"`
		Duration time.Duration `json:"duration"`
		Reason   string        `json:"reason,omitempty"`
		Payload  string        `json:"payload"`
	}

	var suites []Status
	if err := json.Unmarshal(data, &suites); err != nil {
		return err
	}

	suitesByPkg := map[string][]Status{}
	for _, suite := range suites {
		pkg := strings.TrimSuffix(suite.ID, filepath.Ext(suite.ID))
		if filepath.Ext(suite.ID) == "" {
			pkg = "main"
		}

		if _, has := suitesByPkg[pkg]; !has {
			suitesByPkg[pkg] = []Status{}
		}
		suitesByPkg[pkg] = append(suitesByPkg[pkg], suite)
	}

	hasFailed := false
	for pkg, seq := range suitesByPkg {
		durPackage := time.Duration(0)
		hasPackageFailed := false
		for _, unit := range seq {
			durPackage = durPackage + unit.Duration
			if unit.Status == "success" {
				if !silent {
					stdout.Success("==> PASS: %s (%s)\n", unit.ID, unit.Duration)
				}
			} else {
				hasPackageFailed = true
				stdout.Error("==> FAIL: %s (%s)\n", unit.ID, unit.Duration)
				stdout.Warning("%s\n\n", unit.Reason)
			}
			if verbose {
				stdout.FormattedJSON(unit.Payload)
			}
		}

		if hasPackageFailed {
			stdout.Error("FAIL\t%s (%s)\n", pkg, durPackage)
			hasFailed = true
		} else {
			stdout.Success("PASS\t%s (%s)\n", pkg, durPackage)
		}
	}

	if hasFailed {
		return errors.New("failed")
	}

	return nil
}
