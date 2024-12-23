package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func main() {
	// Fetch GitLab token from the environment
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		log.Fatal("Environment variable GITLAB_TOKEN is not set")
	}

	// Configuration
	repoURL := "https://gitlab.cee.redhat.com/rkyatham/go-git.git" // HTTPS URL of the repository
	newBranch := "feature/test-branch"
	commitMessage := "Automated commit using GitLab API by Ram"

	// Inject token into the HTTPS URL
	authRepoURL := fmt.Sprintf("https://oauth2:%s@%s", token, repoURL[8:])

	// Determine the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a unique directory name for the repository clone
	timestamp := time.Now().Format("20060102_150405") // e.g., 20241222_173045
	repoPath := filepath.Join(currentDir, fmt.Sprintf("repo_clone_%s", timestamp))

	// Step 1: Clone the repository using the token
	cloneRepository(authRepoURL, repoPath)

	// Step 2: Switch to the cloned repository
	err = os.Chdir(repoPath)
	if err != nil {
		log.Fatalf("Failed to switch to repository: %v", err)
	}

	// Step 3: Create a new branch
	runGitCommand("git", "checkout", "-b", newBranch)

	// Step 4: Create a dummy file and commit changes
	dummyFile := filepath.Join(repoPath, "dummy.txt")
	err = os.WriteFile(dummyFile, []byte("This is a dummy file."), 0644)
	if err != nil {
		log.Fatalf("Failed to create dummy file: %v", err)
	}

	runGitCommand("git", "add", ".")
	runGitCommand("git", "commit", "-m", commitMessage)

	// Step 5: Push the branch to GitLab
	runGitCommand("git", "push", "-u", "origin", newBranch)

	fmt.Println("Branch pushed successfully using HTTPS and token authentication.")
}

// cloneRepository clones the Git repository to the specified path using HTTPS and token
func cloneRepository(repoURL, repoPath string) {
	fmt.Printf("Cloning repository %s into %s...\n", repoURL, repoPath)
	cmd := exec.Command("git", "clone", repoURL, repoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}
	fmt.Println("Repository cloned successfully.")
}

// runGitCommand runs a Git command
func runGitCommand(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Git command failed: %v", err)
	}
}
