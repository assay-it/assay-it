//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/assay-it/assay/api"
	"github.com/fogfish/gurl"
	"github.com/spf13/cobra"
)

var (
	target string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&target, "url", "", "explicitly define target url to run quality check agains the deployment.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the quality assurance at assay.it",
	Long: `run the quality assurance at assay.it using test suites from
	supplied source code repository. The repository has to be linked with 
	the service beforehand.`,
	Example: `
assay run facebadge/sample.assay.it --key Z2l0aHV...bWhaQQ
assay run facebadge/sample.assay.it --url https://example.com --key Z2l0aHV...bWhaQQ
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	Args:         cobra.ExactArgs(1),
	RunE:         run,
}

func run(cmd *cobra.Command, args []string) error {
	sc := strings.Split(args[0], "/")
	if len(sc) < 2 {
		return fmt.Errorf("invalid source code identity: %s, :owner/:repo is required", args[0])
	}
	sourcecode := fmt.Sprintf("github:%s:%s", sc[0], sc[1])

	c := api.New(endpoint)
	return eval(
		gurl.Join(
			c.SignIn(digest),
			c.WebHookSourceCode(sourcecode, target),
		),
	)
}
