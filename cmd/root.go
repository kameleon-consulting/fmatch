package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/mlabate/fmatch/internal/comparator"
	"github.com/mlabate/fmatch/internal/output"
	"github.com/spf13/cobra"
)

// Version is injected at build time via ldflags:
//
//	-X github.com/mlabate/fmatch/cmd.Version=<version>
var Version = "dev"

// ExitError represents a deliberate exit with a specific code.
// Returned from RunE instead of calling os.Exit directly, to allow testing.
type ExitError struct {
	Code int    // exit code: 1 (different) or 2 (error)
	Msg  string // message to print to stderr; empty = no message
}

func (e *ExitError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return fmt.Sprintf("exit status %d", e.Code)
}

// NewRootCmd creates and configures the root cobra command.
// Exported to allow use in tests.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "fmatch <path_a> <path_b>",
		Version:       Version,
		Short:         "Compare two files or directories for exact equality",
		Long: `fmatch compares two files or directories and determines
whether they are exactly identical.

Exit codes:
  0 - files/directories are identical
  1 - differences found
  2 - error (file not found, permission denied, type mismatch)`,
		Args:          cobra.ExactArgs(2),
		SilenceErrors: true, // error printing handled by Execute()
		SilenceUsage:  true, // don't print usage on runtime errors
		RunE:          runE,
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

// runE is the core command logic. Returns *ExitError for controlled exits (code 1 or 2).
func runE(cmd *cobra.Command, args []string) error {
	pathA, pathB := args[0], args[1]

	// Read flags.
	quiet, _ := cmd.Flags().GetBool("quiet")
	verboseCount, _ := cmd.Flags().GetCount("verbose")
	noColor, _ := cmd.Flags().GetBool("no-color")

	verbosity := resolveVerbosity(quiet, verboseCount)

	// Stat both paths — exit 2 if either does not exist or is unreadable.
	infoA, err := os.Stat(pathA)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}
	infoB, err := os.Stat(pathB)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}

	// Type mismatch: cannot compare file with directory.
	if infoA.IsDir() != infoB.IsDir() {
		return &ExitError{Code: 2, Msg: "fmatch: cannot compare file with directory"}
	}

	// Directory comparison — not yet implemented (Step 7).
	if infoA.IsDir() {
		return &ExitError{Code: 2, Msg: "fmatch: directory comparison not yet implemented"}
	}

	// File comparison.
	result, err := comparator.CompareFiles(pathA, pathB)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}

	// Format and print output.
	opts := output.Options{
		Level:   verbosity,
		NoColor: noColor,
		PathA:   pathA,
		PathB:   pathB,
	}
	out, err := output.Format(result, opts)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}
	if out != "" {
		fmt.Fprintln(cmd.OutOrStdout(), out)
	}

	// Exit 1 if files differ.
	if !result.Identical {
		return &ExitError{Code: 1}
	}
	return nil
}

// resolveVerbosity maps CLI flags to an output.Verbosity level.
func resolveVerbosity(quiet bool, count int) output.Verbosity {
	switch {
	case quiet:
		return output.VerbosityQuiet
	case count >= 2:
		return output.VerbosityVV
	case count == 1:
		return output.VerbosityVerbose
	default:
		return output.VerbosityNormal
	}
}

// rootCmd is the package-level root command instance.
var rootCmd = NewRootCmd()

// Execute runs the root command. Called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			if exitErr.Msg != "" {
				fmt.Fprintln(os.Stderr, exitErr.Msg)
			}
			os.Exit(exitErr.Code)
		}
		// Unexpected error (e.g. flag parsing).
		fmt.Fprintf(os.Stderr, "fmatch: %v\n", err)
		os.Exit(2)
	}
}

