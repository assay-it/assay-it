package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initGoCmd)
}

var initGoCmd = &cobra.Command{
	Use:   "init",
	Short: "creates an example suites into current directory",
	Long:  `creates an example suites into current directory`,
	Example: `
	assay-it init
	`,
	SilenceUsage: true,
	RunE:         initGo,
}

func initGo(cmd *cobra.Command, args []string) error {
	err := os.WriteFile(".assay-it.json", []byte(testSpecConfig()), 0644)
	if err != nil {
		return err
	}

	err = os.MkdirAll("suites", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile("suites/suite.go", []byte(testSpecGo()), 0644)
	if err != nil {
		return err
	}

	stdout.Notice("\n\nrun `go mod tidy` finalize config.\n")

	return nil
}
