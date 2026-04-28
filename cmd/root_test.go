package cmd_test

import (
	"testing"

	"github.com/mlabate/fmatch/cmd"
)

// TestNewRootCmd_FlagDefaults verifies that all flags exist and have the correct default values.
func TestNewRootCmd_FlagDefaults(t *testing.T) {
	c := cmd.NewRootCmd()

	tests := []struct {
		name       string
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
