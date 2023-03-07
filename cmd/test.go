//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/assay-it/assay-it/internal/gocc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&testConfig, "config", "c", "", "path to assay-it config file (default .assay-it.json)")
	testCmd.Flags().StringVarP(&testBuildDir, "build-dir", "b", "", "build dir to cache packages and build artefact (default os temp)")
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "v", false, "enable verbose output of tests results")
	testCmd.Flags().StringVarP(&testSut, "sut", "u", "", "url to system under test")
}

var (
	testConfig   string
	testBuildDir string
	testVerbose  bool
	testSut      string
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "run quality job for either latests or specific commit of source code repository",
	Long: `
run the quality assurance at assay.it using test suites from
supplied source code repository. The repository has to be linked with 
the service beforehand.`,
	Example: `
assay run facebadge/sample.assay.it --key Z2l0aHV...bWhaQQ
assay run facebadge/sample.assay.it/master/8c7ec...dc59 --key Z2l0aHV...bWhaQQ
assay run facebadge/sample.assay.it --url https://example.com --key Z2l0aHV...bWhaQQ
	`,
	SilenceUsage: true,
	RunE:         test,
}

func test(cmd *cobra.Command, args []string) error {
	suites, err := testConfigFile(args)
	if err != nil {
		return fmt.Errorf("no suites are defined or config file is missing: %w", err)
	}

	//
	// Config Sandbox
	//

	if testBuildDir == "" {
		testBuildDir = filepath.Join(os.TempDir(), "assay-it")
	}

	sandbox, err := gocc.NewSandbox(os.Stderr, testBuildDir)
	if err != nil {
		return fmt.Errorf("unable to config build-dir %s: %w", testBuildDir, err)
	}
	stderr.Info("==> env %s\n", testBuildDir)

	//
	// Config Package
	//

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	pkg, err := gocc.NewPackage(sandbox.Path, filepath.Base(dir))
	if err != nil {
		return fmt.Errorf("unable to config package %s: %w", filepath.Base(dir), err)
	}

	//
	// Associate suites with package
	//
	units := []string{}

	for _, suite := range suites {
		stderr.Info("use: %s\n", suite)
		seq, err := pkg.CopyFrom(suite)
		if err != nil {
			return fmt.Errorf("unable to copy %s: %w", suite, err)
		}
		units = append(units, seq...)
	}

	//
	// Crate mandatory files (main.go & mod)
	//
	err = pkg.CreateRunner(testSut, units)
	if err != nil {
		return err
	}

	err = pkg.CreateMod()
	if err != nil {
		return err
	}

	//
	// Compile
	//

	stderr.Info("\n==> compile \n")
	err = sandbox.Compile(pkg)
	if err != nil {
		return err
	}

	//
	// Execute
	//

	stdout.Notice("\n==> testing \n")
	buf := bytes.Buffer{}
	err = pkg.Run(&buf)
	if err != nil {
		panic(err)
	}

	return testWriteResults(buf.Bytes())
}

func testConfigFile(args []string) ([]string, error) {
	type Config struct {
		Suites []string `json:"suites"`
	}

	if testConfig == "" && len(args) != 0 {
		return args, nil
	}

	if testConfig == "" {
		testConfig = ".assay-it.json"
	}

	bytes, err := os.ReadFile(testConfig)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return config.Suites, nil
}

func testWriteResults(data []byte) error {
	type Status struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Duration string `json:"duration"`
		Reason   string `json:"reason,omitempty"`
		Payload  string `json:"payload"`
	}

	var seq []Status
	if err := json.Unmarshal(data, &seq); err != nil {
		return err
	}

	pass := true
	for _, unit := range seq {
		if unit.Status == "success" {
			stdout.Success("PASS: %s\n", unit.ID)
		} else {
			pass = false
			stdout.Error("FAIL: %s\n", unit.ID)
			stdout.Warning("%s\n", unit.Reason)
		}
		if testVerbose {
			stdout.FormattedJSON(unit.Payload)
		}
	}

	if pass {
		stdout.Success("\nPASS\n")
	} else {
		stdout.Error("\nFAIL\n")
		os.Exit(1)
	}
	return nil
}
