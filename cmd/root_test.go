package cmd_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/mlabate/fmatch/cmd"
)

// writeTemp creates a temporary file with the given content.
func writeTemp(t *testing.T, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "fmatch-cmd-*")
	if err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		t.Fatalf("writeTemp write: %v", err)
	}
	return f.Name()
}

// ── Flag defaults ────────────────────────────────────────────────────────────

// TestNewRootCmd_FlagDefaults verifies that all flags exist and have the correct default values.
func TestNewRootCmd_FlagDefaults(t *testing.T) {
	c := cmd.NewRootCmd()

	tests := []struct {
		name        string
		wantDefault string
	}{
		{"quiet", "false"},
		{"verbose", "0"},
		{"depth", "-1"},
		{"ignore", "[]"},
		{"ignore-file", ".fmatchignore"},
		{"no-ignore", "false"},
		{"no-follow-symlinks", "false"},
		{"no-color", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := c.Flags().Lookup(tt.name)
			if f == nil {
				t.Fatalf("flag --%s not registered", tt.name)
			}
			if f.DefValue != tt.wantDefault {
				t.Errorf("flag --%s: want default %q, got %q", tt.name, tt.wantDefault, f.DefValue)
			}
		})
	}
}

// TestNewRootCmd_ExactArgs verifies that the command requires exactly 2 arguments.
func TestNewRootCmd_ExactArgs(t *testing.T) {
	c := cmd.NewRootCmd()

	cases := []struct {
		args    []string
		wantErr bool
	}{
		{[]string{}, true},
		{[]string{"a"}, true},
		{[]string{"a", "b"}, false},
		{[]string{"a", "b", "c"}, true},
	}

	for _, tc := range cases {
		err := c.Args(c, tc.args)
		if tc.wantErr && err == nil {
			t.Errorf("args %v: expected error, got nil", tc.args)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("args %v: expected no error, got %v", tc.args, err)
		}
	}
}

// TestVersion_NotEmpty verifies that the Version variable is set.
func TestVersion_NotEmpty(t *testing.T) {
	if cmd.Version == "" {
		t.Error("expected Version to be non-empty, got empty string")
	}
}

// ── Integration (RunE) ───────────────────────────────────────────────────────

func TestRunE_Identical(t *testing.T) {
	content := []byte("hello fmatch!")
	a, b := writeTemp(t, content), writeTemp(t, content)

	c := cmd.NewRootCmd()
	var buf bytes.Buffer
	c.SetOut(&buf)
	c.SetArgs([]string{"--no-color", a, b})

	err := c.Execute()
	if err != nil {
		t.Fatalf("expected no error for identical files, got: %v", err)
	}
	if !strings.Contains(buf.String(), "IDENTICAL") {
		t.Errorf("expected IDENTICAL in output, got %q", buf.String())
	}
}

func TestRunE_Different_ExitCode1(t *testing.T) {
	a := writeTemp(t, []byte("aaa"))
	b := writeTemp(t, []byte("bbb"))

	c := cmd.NewRootCmd()
	c.SetArgs([]string{"--no-color", a, b})

	err := c.Execute()
	if err == nil {
		t.Fatal("expected ExitError for different files, got nil")
	}
	var exitErr *cmd.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *cmd.ExitError, got %T: %v", err, err)
	}
	if exitErr.Code != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.Code)
	}
}

func TestRunE_FileNotFound_ExitCode2(t *testing.T) {
	c := cmd.NewRootCmd()
	c.SetArgs([]string{"/nonexistent/path/a", "/nonexistent/path/b"})

	err := c.Execute()
	var exitErr *cmd.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *cmd.ExitError, got %T: %v", err, err)
	}
	if exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
}

func TestRunE_TypeMismatch_ExitCode2(t *testing.T) {
	dir := t.TempDir()
	file := writeTemp(t, []byte("hello"))

	c := cmd.NewRootCmd()
	c.SetArgs([]string{dir, file})

	err := c.Execute()
	var exitErr *cmd.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *cmd.ExitError, got %T: %v", err, err)
	}
	if exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
}

func TestRunE_Quiet_Identical(t *testing.T) {
	content := []byte("hello quiet")
	a, b := writeTemp(t, content), writeTemp(t, content)

	c := cmd.NewRootCmd()
	var buf bytes.Buffer
	c.SetOut(&buf)
	c.SetArgs([]string{"-q", a, b})

	err := c.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "" {
		t.Errorf("quiet mode: expected empty output, got %q", buf.String())
	}
}

