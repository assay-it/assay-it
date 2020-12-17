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
	head   string
	base   string
	title  string
	number string
)

func init() {
	//
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.PersistentFlags().StringVar(&title, "title", "", "human readable description of the change")
	webhookCmd.PersistentFlags().StringVar(&number, "number", "", "unique reference number of the change (e.g. pull request)")
	webhookCmd.PersistentFlags().StringVar(&target, "url", "", "explicitly define target url to run quality check agains the deployment.")

	webhookCmd.Flags().StringVar(&head, "head", "", "head reference of the change using the format :owner/:repo/:branch/:commit")
	webhookCmd.MarkFlagRequired("head")

	webhookCmd.Flags().StringVar(&base, "base", "", "base reference of the change using the format :owner/:repo/:branch/:commit")
	webhookCmd.MarkFlagRequired("base")

	//
	webhookCmd.AddCommand(webhookSourceCmd)

	//
	webhookCmd.AddCommand(webhookCommitCmd)

	//
	webhookCmd.AddCommand(webhookReleaseCmd)
}

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "run quality job via webhook at assay.it",
	Long: `the command triggers quality assurance at assay.it using webhook api,
the command faciliatets CI/CD integration use-cases.`,
	Example: `
assay webhook source facebadge/sample.assay.it
assay webhook commit facebadge/sample.assay.it/main/8c7ec...dc59
assay webhook branch --base facebadge/sample.assay.it/main/8c7ec...dc59 --head faceb.../feature/8c7ec...dc59
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	Args:         cobra.NoArgs,
	RunE:         webhook,
}

func webhook(cmd *cobra.Command, args []string) error {
	chead := strings.Split(head, "/")
	cbase := strings.Split(base, "/")

	if len(chead) != 4 {
		return fmt.Errorf("invalid head identity: %s", head)
	}

	if len(cbase) != 4 {
		return fmt.Errorf("invalid base identity: %s", base)
	}

	c := api.New(endpoint)
	return eval(
		assay.Join(
			c.SignIn(digest),
			c.WebHook(
				api.Hook{
					Base:        api.Commit{ID: fmt.Sprintf("[github:%s/%s/%s/%s]", cbase[0], cbase[1], cbase[2], cbase[3])},
					Head:        api.Commit{ID: fmt.Sprintf("[github:%s/%s/%s/%s]", chead[0], chead[1], chead[2], chead[3])},
					URL:         target,
					PullRequest: mkPullRequest(),
				},
			),
		),
	)

}

//
//
var webhookSourceCmd = &cobra.Command{
	Use:   "source",
	Short: "run quality job for latest commit of source code repository",
	Long: `
the command triggers quality assurance at assay.it using webhook api,
it schedule quality job for latest commit of source code repository.`,
	Example: `
assay webhook source facebadge/sample.assay.it
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	Args:         cobra.ExactArgs(1),
	RunE:         webhookSource,
}

func webhookSource(cmd *cobra.Command, args []string) error {
	sc := strings.Split(args[0], "/")
	if len(sc) != 2 {
		return fmt.Errorf("invalid source code identity: %s", args[0])
	}

	c := api.New(endpoint)
	return eval(
		assay.Join(
			c.SignIn(digest),
			c.WebHookSource(
				api.SourceCodeID{
					ID:          fmt.Sprintf("[github:%s/%s]", sc[0], sc[1]),
					URL:         target,
					PullRequest: mkPullRequest(),
				},
			),
		),
	)
}

//
//
var webhookCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "run quality job for specific commit of source code repository",
	Long: `
the command triggers quality assurance at assay.it using webhook api,
it schedule quality job for specific commit of source code repository.`,
	Example: `
assay webhook commit facebadge/sample.assay.it/main/8c7ec...dc59
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	Args:         cobra.ExactArgs(1),
	RunE:         webhookCommit,
}

func webhookCommit(cmd *cobra.Command, args []string) error {
	sc := strings.Split(args[0], "/")
	if len(sc) != 4 {
		return fmt.Errorf("invalid commit identity: %s", args[0])
	}

	c := api.New(endpoint)
	return eval(
		assay.Join(
			c.SignIn(digest),
			c.WebHookCommit(
				api.SourceCodeID{
					ID:          fmt.Sprintf("[github:%s/%s/%s/%s]", sc[0], sc[1], sc[2], sc[3]),
					URL:         target,
					PullRequest: mkPullRequest(),
				},
			),
		),
	)
}

var webhookReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "run quality job for specific release (tag) of source code repository",
	Long: `the command triggers quality assurance at assay.it using webhook api,
the command faciliatets CI/CD integration use-cases.`,
	Example: `
assay webhook release facebadge/sample.assay.it/main/v0
	`,
	SilenceUsage: true,
	PreRunE:      requiredFlagKey,
	Args:         cobra.ExactArgs(1),
	RunE:         webhookRelease,
}

func webhookRelease(cmd *cobra.Command, args []string) error {
	sc := strings.Split(args[0], "/")
	if len(sc) != 4 {
		return fmt.Errorf("invalid commit identity: %s", args[0])
	}

	c := api.New(endpoint)
	return eval(
		assay.Join(
			c.SignIn(digest),
			c.WebHookRelease(
				api.SourceCodeID{
					ID:          fmt.Sprintf("[github:%s/%s/%s/%s]", sc[0], sc[1], sc[2], sc[3]),
					URL:         target,
					PullRequest: mkPullRequest(),
				},
			),
		),
	)
}

//
//
func mkPullRequest() *api.PullRequest {
	if number != "" && title != "" {
		return &api.PullRequest{
			Number: number,
			Title:  title,
		}
	}
	return nil
}
