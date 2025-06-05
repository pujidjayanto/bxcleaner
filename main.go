package main

import (
	"bufio"
	"bytes"
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

	cmd = exec.Command("git", "branch")
	listBranchBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error getting branches:", err)
		os.Exit(1)
	}

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
		fmt.Println("no branches to delete. everything done, bye!")
		return
	}

	for _, branch := range branchesToDelete {
		fmt.Printf("deleting branch: %s\n", branch)
		cmdDel := exec.Command("git", "branch", "-D", branch)
		cmdDel.Stdout = os.Stdout
		cmdDel.Stderr = os.Stderr
		if err := cmdDel.Run(); err != nil {
			fmt.Printf("failed to delete branch %s: %v\n", branch, err)
		}
	}
}
