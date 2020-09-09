//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"github.com/assay-it/assay/api"
	"github.com/spf13/cobra"
)

var (
	secretCheck bool
)

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.Flags().BoolVarP(&secretCheck, "check", "c", false, "validates secret key")
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "checks access credentials",
	Long:  `checks access credentials`,
	Example: `
assay secret --check --key Z2l0aHV...bWhaQQ
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	RunE:         secret,
}

func secret(cmd *cobra.Command, args []string) error {
	if secretCheck {
		return eval(api.New(endpoint).SignIn(digest))
	}

	return nil
}
