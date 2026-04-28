package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Version is injected at build time via ldflags:
//
//	-X github.com/mlabate/fmatch/cmd.Version=<version>
var Version = "dev"

// NewRootCmd creates and configures the root cobra command.
// Exported to allow use in tests.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fmatch <path_a> <path_b>",
		Version: Version,
		Short:   "Compare two files or directories for exact equality",
		Long: `fmatch compares two files or directories and determines
whether they are exactly identical.

Exit codes:
  0 - files/directories are identical
  1 - differences found
  2 - error (file not found, permission denied, type mismatch)`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implemented in Step 5 (integration)
			return nil
		},
	}

	// Verbosity
	cmd.Flags().BoolP("quiet", "q", false, "Quiet mode: exit code only, no output")
	cmd.Flags().CountP("verbose", "v", "Verbose output (repeatable: -vv for extra detail)")

	// Directory options
	cmd.Flags().IntP("depth", "d", -1, "Maximum depth for directory traversal (-1 = unlimited)")

	// Ignore patterns
	cmd.Flags().StringArrayP("ignore", "i", []string{}, "Additional pattern to ignore (repeatable)")
	cmd.Flags().String("ignore-file", ".fmatchignore", "Path to pattern ignore file")
	cmd.Flags().Bool("no-ignore", false, "Disable .fmatchignore file")

	// Symlinks
	cmd.Flags().Bool("no-follow-symlinks", false, "Do not follow symlinks (default: follow)")

	// Output
	cmd.Flags().Bool("no-color", false, "Disable colored output")

	return cmd
}

// rootCmd is the package-level root command instance.
var rootCmd = NewRootCmd()

// Execute runs the root command. Called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
