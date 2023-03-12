//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"os"

	"github.com/assay-it/assay-it/internal/tester"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
	// TODO: timeout
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "v", false, "enable verbose output of tests results")
	testCmd.Flags().StringVarP(&testSut, "target", "t", "", "default url for the system under test")
}

var (
	testVerbose bool
	testSilent  bool
	testSut     string
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "runs testing of services and deployments using Golang module.",
	Long: `
'assay-it test' command automates testing of services and deployments using the
Golang module. It prints a summary of the test results in the format: 

	PASS: TestPetShopList (3.329073ms)
	PASS: TestPetShopListWithCursor (2.329073ms)
	FAIL: TestPetShopCreate (6.329073ms)

The command uses the configuration file (.assay-it.json) from current directory
to identify all suites to be executed against deployments. The config file is
plain JSON document, the key 'suites' enumerates all files to be used by
the utility:

	{
	  "suites": [
	    "suites/petshop.go"
	  ]
	}

These files can contain test specifications, data structure and utility
functions written on Golang that scripts the quality assessment of deployments.

The utility generates runner.go source code to bootstrap testing, compiles
modules listed in the configuration file, links the executable and run with
the main test binary.

As part of building a test binary, 'assay-it test' runs standard go tools.
If go tools finds any problems, assay-it reports those and does not run the test
binary.

Only status and summary lines are printed to standard output. It also print the
root cause if test case is failed. In the verbose mode (-v flag) the command
prints the response payload returned by deployments. 

By default, the command runs in the local directory mode, occurs when
'assay-it test' is invoked with no the configuration file arguments (for example
'assay-it test' or 'assay-it test -v'). In this mode, the command searches for
the the configuration file '.assay-it.json' in current directory and then runs
the resulting test binary. After the test finishes, 'assay-it test' prints a
summary line showing the test status ('PASS' or 'FAIL'), package name, and
elapsed time.

The second mode occurs when when 'assay-it test' is invoked with explicit 
configuration file arguments (for example 'assay-it test petshop', 'assay-it
test petshop vetclinic', and even 'assay-it test ./...'). In this mode, the command
assumes each package contains own configuration file '.assay-it.json'. Each of
the packages listed on the command line is tested independently either
sequentially or parallel. If a package test passes, 'assay-it test' prints only
the final 'PASS' summary line. If a package test fails, 'assay-it test' prints
the full test output. If invoked with -v flag, the full output is printed for 
passing packages. On the exits, the command prints a summary line showing the
test status ('PASS' or 'FAIL') and elapsed time.

For more about testing suites, see 'assay-it help testspec'.

Visit  https://assay.it
`,
	Example: `
assay-it test
assay-it test -v -t https://assay.it
assay-it test ./...
	`,
	SilenceUsage: true,
	RunE:         test,
}

func test(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		testSilent = false
		if err := testModeLocal(cmd, ""); err != nil {
			stdout.Error("\nFAIL\n")
			os.Exit(1)
		}

		stdout.Success("\nPASS\n")
		return nil
	}

	passed := true
	for _, pkg := range args {
		testSilent = true
		if err := testModeLocal(cmd, pkg); err != nil {
			passed = false
		}
	}

	if !passed {
		stdout.Error("\nFAIL\n")
		os.Exit(1)
	}

	stdout.Success("\nPASS\n")
	return nil
}

func testModeLocal(cmd *cobra.Command, pkg string) error {
	tt, err := tester.NewTester(pkg)
	if err != nil {
		return err
	}

	if err := tt.AnalyzeSourceCode(nil); err != nil {
		return err
	}

	if err := tt.CreateRunner(); err != nil {
		return err
	}

	out, err := tt.Test(testSut)
	if err != nil {
		stderr.Write(out)
		return err
	}

	return stdoutTestResults(testSilent, testVerbose, out)
}
