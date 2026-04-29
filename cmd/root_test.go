package cmd_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/kameleon-consulting/fmatch/cmd"
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

// TestNewRootCmd_ArgCount verifies that the command accepts 1 or 2 arguments (RangeArgs).
func TestNewRootCmd_ArgCount(t *testing.T) {
	c := cmd.NewRootCmd()

	cases := []struct {
		args    []string
		wantErr bool
	}{
		{[]string{}, true},              // 0 args: rejected
		{[]string{"a"}, false},          // 1 arg: valid (single-dir mode)
		{[]string{"a", "b"}, false},     // 2 args: valid
		{[]string{"a", "b", "c"}, true}, // 3 args: rejected
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
// ── Single-arg mode ──────────────────────────────────────────────────────────

// TestRunE_SingleFile_ExitCode2: a single file argument is an error (exit 2).
func TestRunE_SingleFile_ExitCode2(t *testing.T) {
	file := writeTemp(t, []byte("hello"))

	c := cmd.NewRootCmd()
	c.SetArgs([]string{file})

	err := c.Execute()
	var exitErr *cmd.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *cmd.ExitError, got %T: %v", err, err)
	}
	if exitErr.Code != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.Code)
	}
	if !strings.Contains(exitErr.Msg, "second path") {
		t.Errorf("expected message mentioning 'second path', got %q", exitErr.Msg)
	}
}

// TestRunE_SingleDir_NoDuplicates: single dir with no duplicates → exit 0.
func TestRunE_SingleDir_NoDuplicates(t *testing.T) {
	dir := t.TempDir()
	// Write two files with distinct content.
	if err := os.WriteFile(dir+"/a.txt", []byte("content a"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dir+"/b.txt", []byte("content b"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := cmd.NewRootCmd()
	var buf bytes.Buffer
	c.SetOut(&buf)
	c.SetArgs([]string{"--no-color", "--no-ignore", dir})

	err := c.Execute()
	if err != nil {
		t.Fatalf("expected exit 0 (no duplicates), got: %v", err)
	}
	if !strings.Contains(buf.String(), "No duplicates found") {
		t.Errorf("expected 'No duplicates found' in output, got %q", buf.String())
	}
}

// TestRunE_SingleDir_WithDuplicates: single dir with duplicates → exit 1.
func TestRunE_SingleDir_WithDuplicates(t *testing.T) {
	dir := t.TempDir()
	// Write two files with identical content.
	if err := os.WriteFile(dir+"/dup1.txt", []byte("same content"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dir+"/dup2.txt", []byte("same content"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := cmd.NewRootCmd()
	var buf bytes.Buffer
	c.SetOut(&buf)
	c.SetArgs([]string{"--no-color", "--no-ignore", dir})

	err := c.Execute()
	var exitErr *cmd.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *cmd.ExitError (exit 1), got %T: %v", err, err)
	}
	if exitErr.Code != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.Code)
	}
	if !strings.Contains(buf.String(), "duplicate group") {
		t.Errorf("expected 'duplicate group' in output, got %q", buf.String())
	}
}
