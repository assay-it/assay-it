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
	"github.com/assay-it/sdk-go/assay"
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
	PreRunE:      requiredFlagKey,
	Args:         cobra.ExactArgs(1),
	RunE:         run,
}

func run(cmd *cobra.Command, args []string) error {
	sc := strings.Split(args[0], "/")

	c := api.New(endpoint)

	if len(sc) == 2 {
		return eval(
			assay.Join(
				c.SignIn(digest),
				c.WebHookSource(
					api.SourceCodeID{
						ID:  fmt.Sprintf("[github:%s/%s]", sc[0], sc[1]),
						URL: target,
					},
				),
			),
		)
	}

	if len(sc) == 4 {
		return eval(
			assay.Join(
				c.SignIn(digest),
				c.WebHookSource(
					api.SourceCodeID{
						ID:  fmt.Sprintf("[github:%s/%s/%s/%s]", sc[0], sc[1], sc[2], sc[3]),
						URL: target,
					},
				),
			),
		)
	}

	return fmt.Errorf("invalid source code identity: %s", args[0])
}
