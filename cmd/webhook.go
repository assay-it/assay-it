//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package cmd

import (
	"path/filepath"
	"strings"

	"github.com/assay-it/assay/api"
	"github.com/assay-it/sdk-go/assay"
	"github.com/spf13/cobra"
)

var (
	head   string
	base   string
	title  string
	number string
)

func init() {
	rootCmd.AddCommand(webhookCmd)

	webhookCmd.Flags().StringVar(&head, "head", "", "head reference of the change using the format :owner/:repo/:branch/:commit")
	webhookCmd.MarkFlagRequired("head")

	webhookCmd.Flags().StringVar(&base, "base", "", "base reference of the change using the format :owner/:repo/:branch/:commit")
	webhookCmd.MarkFlagRequired("base")

	webhookCmd.Flags().StringVar(&title, "title", "", "human readable description of the change")
	webhookCmd.Flags().StringVar(&number, "number", "", "unique reference number of the change (e.g. pull request)")

	webhookCmd.Flags().StringVar(&target, "url", "", "explicitly define target url to run quality check agains the deployment.")
}

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "run quality job via webhook at assay.it",
	Long: `the command triggers quality assurance at assay.it using webhook api,
the command faciliatets CI/CD integration use-cases.`,
	Example: `
assay webhook --base facebadge/sample.assay.it/master/8c7ec...dc59 --head faceb.../master/8c7ec...dc59
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	Args:         cobra.NoArgs,
	RunE:         webhook,
}

func webhook(cmd *cobra.Command, args []string) error {
	base = strings.Join(strings.Split(filepath.Join("github", base), "/"), ":")
	head = strings.Join(strings.Split(filepath.Join("github", head), "/"), ":")

	hook := api.Hook{
		Base: api.Commit{ID: base},
		Head: api.Commit{ID: head},
		URL:  target,
	}
	if number != "" && title != "" {
		hook.PullRequest = &api.PullRequest{
			Number: number,
			Title:  title,
		}
	}

	c := api.New(endpoint)
	return eval(
		assay.Join(
			c.SignIn(digest),
			c.WebHook(hook),
		),
	)
}
