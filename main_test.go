package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCleanBranches(t *testing.T) {
	// Fake git commands
	mockGitCmd := func(args ...string) ([]byte, error) {
		// Simulate "git branch --show-current"
		if len(args) == 2 && args[0] == "branch" && args[1] == "--show-current" {
			return []byte("main\n"), nil
		}

		// Simulate "git branch" listing
		if len(args) == 1 && args[0] == "branch" {
			return []byte("* main\n  feature1\n  feature2\n"), nil
		}

		// Simulate deleting branch
		if len(args) == 2 && args[0] == "branch" && args[1] == "-D" {
			// Return no error
			return []byte{}, nil
		}

		return nil, nil
	}

	// Capture output
	var stdout, stderr bytes.Buffer

	// Run
	err := cleanBranches("main", mockGitCmd, &stdout, &stderr)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check output
	outStr := stdout.String()
	if !strings.Contains(outStr, "deleting branch: feature1") ||
		!strings.Contains(outStr, "deleting branch: feature2") ||
		!strings.Contains(outStr, "clean complete") {
		t.Errorf("unexpected stdout:\n%s", outStr)
	}

	if stderr.Len() != 0 {
		t.Errorf("unexpected stderr:\n%s", stderr.String())
	}
}
