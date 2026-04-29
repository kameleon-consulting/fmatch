package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/kameleon-consulting/fmatch/internal/comparator"
	"github.com/kameleon-consulting/fmatch/internal/ignore"
	"github.com/kameleon-consulting/fmatch/internal/output"
	"github.com/spf13/cobra"
)

// Version is injected at build time via ldflags:
//
//	-X github.com/kameleon-consulting/fmatch/cmd.Version=<version>
var Version = "dev"

// ExitError represents a deliberate exit with a specific code.
// Returned from RunE instead of calling os.Exit directly, to allow testing.
type ExitError struct {
	Code int    // exit code: 1 (different/duplicates) or 2 (error)
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
		Use:     "fmatch [flags] <path_a> [path_b]",
		Version: Version,
		Short:   "Compare files or directories for exact equality, or find duplicates",
		Long: `fmatch compares two files or directories and determines
whether they are exactly identical.

With a single directory argument, fmatch finds duplicate files
within that directory, grouped by content (SHA-256 hash).

Exit codes:
  0 - identical files/directories, or no duplicates found
  1 - differences found, or duplicates found
  2 - error (file not found, permission denied, type mismatch)`,
		Args:          cobra.RangeArgs(1, 2),
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

// runE is the core command logic. Dispatches to runSingleArg or runTwoArgs
// based on the number of positional arguments.
func runE(cmd *cobra.Command, args []string) error {
	quiet, _ := cmd.Flags().GetBool("quiet")
	verboseCount, _ := cmd.Flags().GetCount("verbose")
	noColor, _ := cmd.Flags().GetBool("no-color")
	verbosity := resolveVerbosity(quiet, verboseCount)

	if len(args) == 1 {
		return runSingleArg(cmd, args[0], verbosity, noColor)
	}
	return runTwoArgs(cmd, args[0], args[1], verbosity, noColor)
}

// runSingleArg handles the single-argument invocation: fmatch <dir>.
// Only directories are valid; a single file argument is an error.
// Finds duplicate files within the directory grouped by content hash.
func runSingleArg(cmd *cobra.Command, path string, verbosity output.Verbosity, noColor bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}
	if !info.IsDir() {
		return &ExitError{Code: 2, Msg: "fmatch: single file argument requires a second path to compare against"}
	}

	matcher, err := loadMatcher(cmd)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}
	depth, _ := cmd.Flags().GetInt("depth")

	dupResult, err := comparator.FindDuplicates(path, comparator.DirOptions{
		Matcher: matcher,
		Depth:   depth,
	})
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}

	opts := output.Options{
		Level:   verbosity,
		NoColor: noColor,
		PathA:   path,
	}
	out, err := output.FormatDuplicates(dupResult, opts)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}
	if out != "" {
		fmt.Fprintln(cmd.OutOrStdout(), out)
	}

	// Exit 1 if duplicates were found.
	if dupResult.HasDuplicates {
		return &ExitError{Code: 1}
	}
	return nil
}

// runTwoArgs handles the two-argument invocation: fmatch <path_a> <path_b>.
// Compares two files (byte-by-byte) or two directories (hash-based).
// Mixing a file with a directory is an error (exit 2).
func runTwoArgs(cmd *cobra.Command, pathA, pathB string, verbosity output.Verbosity, noColor bool) error {
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

	// Directory comparison (hash-based).
	if infoA.IsDir() {
		matcher, err := loadMatcher(cmd)
		if err != nil {
			return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
		}
		depth, _ := cmd.Flags().GetInt("depth")

		dirResult, err := comparator.CompareDir(pathA, pathB, comparator.DirOptions{
			Matcher: matcher,
			Depth:   depth,
		})
		if err != nil {
			return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
		}

		opts := output.Options{
			Level:   verbosity,
			NoColor: noColor,
			PathA:   pathA,
			PathB:   pathB,
		}
		out, err := output.FormatDirCompare(dirResult, opts)
		if err != nil {
			return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
		}
		if out != "" {
			fmt.Fprintln(cmd.OutOrStdout(), out)
		}
		if !dirResult.Identical {
			return &ExitError{Code: 1}
		}
		return nil
	}

	// File comparison (byte-by-byte).
	result, err := comparator.CompareFiles(pathA, pathB)
	if err != nil {
		return &ExitError{Code: 2, Msg: fmt.Sprintf("fmatch: %v", err)}
	}

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

// loadMatcher builds an ignore.Matcher from the --ignore-file, --no-ignore and -i flags.
// If --no-ignore is set, returns an empty Matcher (nothing ignored).
// Otherwise combines the ignore file patterns with any -i inline patterns.
func loadMatcher(cmd *cobra.Command) (*ignore.Matcher, error) {
	noIgnore, _ := cmd.Flags().GetBool("no-ignore")
	if noIgnore {
		return ignore.LoadPatterns([]string{}), nil
	}
	ignoreFile, _ := cmd.Flags().GetString("ignore-file")
	extra, _ := cmd.Flags().GetStringArray("ignore")
	return ignore.LoadFileAndPatterns(ignoreFile, extra)
}
