package main

import (
	"fmt"
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
  bxcleaner -h  		# Show this help message
`

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 && (argsWithoutProg[0] == "--help" || argsWithoutProg[0] == "-h") {
		fmt.Print(usage)
		return
	}

	defaultBranch := "main"
	if len(argsWithoutProg) > 0 {
		defaultBranch = argsWithoutProg[0]
	}

	cmd := exec.Command("git", "branch", "--show-current")
	currentBranchBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error getting current branch:", err)
		os.Exit(1)
	}
	currentBranch := strings.TrimSpace(string(currentBranchBytes))

	if currentBranch != defaultBranch {
		fmt.Printf("you need to be in chosen branch, use git checkout %s\n", defaultBranch)
		os.Exit(1)
	}
}
