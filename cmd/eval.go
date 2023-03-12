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
	"path/filepath"

	"github.com/assay-it/assay-it/internal/tester"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(evalCmd)
	evalCmd.Flags().StringVarP(&evalBuildDir, "build-dir", "b", "", "build dir to cache packages and build artefact (default os temp)")
	evalCmd.Flags().BoolVarP(&evalVerbose, "verbose", "v", false, "enable verbose output of tests results")
	evalCmd.Flags().StringVarP(&testSut, "target", "t", "", "default url for the system under test")
}

var (
	evalBuildDir string
	evalVerbose  bool
	evalSilent   bool
	evalSut      string
)

var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "evaluate quality services and deployments",
	Long: `
'assay-it eval' command automates testing of services and deployments using
testing suites for protocol endpoints. The command prints a summary of the
quality assessment in the format:

	PASS: TestPetShopList (3.329073ms)
	PASS: TestPetShopListWithCursor (2.329073ms)
	FAIL: TestPetShopCreate (6.329073ms)

These scenarios check the correctness and makes the formal proof of
quality in loosely coupled topologies such as serverless applications,
microservices and other systems that rely on interface syntaxes and its
behaviors. Testing suites are type safe and pure functional test specification
of protocol endpoints exposed by software components. The command requires the
suite development using Golang syntax but limited functionality is supported
with Markdown documents.

The command 'assay-it eval' requires Golang development environment for
execution. It builds a "sandbox" workplace to compiles suites,links the
executable and run with the main test binary. 

By default, the command runs in the local directory mode, occurs when
'assay-it eval' is invoked with no testing suite file arguments (for example
'assay-it eval' or 'assay-it eval -v'). In this mode, the command searches for
the the configuration file '.assay-it.json' in current directory and then runs
the resulting test binary. After the test finishes, 'assay-it eval' prints a
summary line showing the test status ('PASS' or 'FAIL'), package name, and
elapsed time.

The configuration file (.assay-it.json) just list all testing suites to be
executed. The config file is plain JSON document, the key 'suites' enumerates
all files to be used by the utility:

	{
	  "suites": [
	    "suites/petshop.go"
	  ]
	}

The second mode occurs when when 'assay-it eval' is invoked with explicit
definition of test suite arguments (for example 'assay-it eval petshop.go',
'assay-it test petshop.md test/vetclinic.go'). In this mode, the command build
testing binary from these files. 


assumes listed files to be used for testing. 

each package contains own configuration file '.assay-it.json'. Each of
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
	assay-it eval
	assay-it eval -v -t https://assay.it
	assay-it eval petshop.go
	`,
	SilenceUsage: true,
	RunE:         eval,
}

func eval(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	tt, err := tester.NewEvaler(evalBuildDir, filepath.Base(dir), args)
	if err != nil {
		return err
	}

	stderr.Info("\n==> prepare \n")
	if err := tt.AnalyzeSourceCode(stderr); err != nil {
		return err
	}

	if err := tt.CreateRunner(); err != nil {
		return err
	}

	stderr.Info("\n==> compile \n")
	if err := tt.Compile(); err != nil {
		return err
	}

	stdout.Notice("\n==> testing \n")
	out, err := tt.Test(evalSut)
	if err != nil {
		stderr.Write(out)
		return err
	}

	return stdoutTestResults(evalSilent, evalVerbose, out)
}
