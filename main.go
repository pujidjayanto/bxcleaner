package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const usage = `bxcleaner - A simple tool to delete local Git branches except the chosen branch

Usage:
  bxcleaner [default-branch]

Examples:
  bxcleaner         # Keeps 'main' by default, deletes other branches
  bxcleaner develop # Keeps 'develop' branch, deletes others
  bxcleaner --help  # Show this help message
  bxcleaner -h      # Show this help message
`

func main() {
	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Print(usage)
		return
	}

	defaultBranch := "main"
	if len(args) > 0 {
		defaultBranch = args[0]
	}

	// Use real git commands
	gitCmd := func(args ...string) ([]byte, error) {
		cmd := exec.Command("git", args...)
		return cmd.Output()
	}

	err := cleanBranches(defaultBranch, gitCmd, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func cleanBranches(defaultBranch string, gitCmd func(args ...string) ([]byte, error), stdout, stderr io.Writer) error {
	// Check if we're in a git repository first
	_, err := gitCmd("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	/* Get current branch
	1. git branch --show-current (since git v 2.22)
	2. git rev-parse --abbrev-ref HEAD
	3. git branch | sed -n '/\* /s///p' (i think this one can done programmatically)
	https://stackoverflow.com/questions/6245570/how-do-i-get-the-current-branch-name-in-git
	*/
	currentBranchBytes, err := gitCmd("branch", "--show-current")
	if err != nil {
		return fmt.Errorf("error getting current branch: %w", err)
	}
	currentBranch := strings.TrimSpace(string(currentBranchBytes))

	if currentBranch != defaultBranch {
		fmt.Fprintf(stderr, "you need to be in chosen branch, use git checkout %s\n", defaultBranch)
		return fmt.Errorf("not on chosen branch")
	}

	// List all branches
	listBranchBytes, err := gitCmd("branch")
	if err != nil {
		return fmt.Errorf("error getting branches: %w", err)
	}

	// Determine which to delete
	var branchesToDelete []string
	scanner := bufio.NewScanner(bytes.NewReader(listBranchBytes))
	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())
		branch = strings.TrimPrefix(branch, "* ")
		if branch != defaultBranch && branch != currentBranch {
			branchesToDelete = append(branchesToDelete, branch)
		}
	}

	if len(branchesToDelete) == 0 {
		fmt.Fprintln(stdout, "no branches to delete. clean complete, bye!")
		return nil
	}

	for _, branch := range branchesToDelete {
		fmt.Fprintf(stdout, "deleting branch: %s\n", branch)
		_, err := gitCmd("branch", "-D", branch)
		if err != nil {
			fmt.Fprintf(stderr, "failed to delete branch %s: %v\n", branch, err)
		}
	}

	fmt.Fprintln(stdout, "clean complete, bye!")
	return nil
}
