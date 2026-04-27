package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the main entry point for the fmatch CLI.
var rootCmd = &cobra.Command{
	Use:   "fmatch <path_a> <path_b>",
	Short: "Compare two files or directories for exact equality",
	Long: `fmatch compares two files or directories and determines
whether they are exactly identical.

Exit codes:
  0 - files/directories are identical
  1 - differences found
  2 - error (file not found, permission denied, type mismatch)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implemented in Step 5 (integration)
		return nil
	},
}

// Execute runs the root command. Called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
