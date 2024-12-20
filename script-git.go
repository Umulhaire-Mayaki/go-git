package main

import (
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	gitlab "https://gitlab.cee.redhat.com/umayaki/clusterimagesets"
)

func main() {
	// Step 1: Clone the repository or open an existing one
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Failed to open repo: %v", err)
	}

	// Step 2: Stage all changes
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to get worktree: %v", err)
	}

	err = worktree.AddGlob(".")
	if err != nil {
		log.Fatalf("Failed to add changes: %v", err)
	}

	// Step 3: Commit the changes
	commit, err := worktree.Commit("Automated commit message", &git.CommitOptions{})
	if err != nil {
		log.Fatalf("Failed to commit changes: %v", err)
	}

	// Step 4: Push the changes to remote
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		log.Fatalf("Failed to push changes: %v", err)
	}
	fmt.Println("Changes pushed successfully.")

	// Step 5: Create a Merge Request using GitLab API
	token := os.Getenv("GITLAB_TOKEN")
	gitlabClient, err := gitlab.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	mergeRequest, _, err := gitlabClient.MergeRequests.CreateMergeRequest("your-namespace/your-repo", &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String("Automated MR"),
		SourceBranch: gitlab.String("your-branch"),
		TargetBranch: gitlab.String("main"),
		Description:  gitlab.String("This MR was created automatically."),
	})
	if err != nil {
		log.Fatalf("Failed to create MR: %v", err)
	}

	fmt.Printf("Merge Request created: %s\n", mergeRequest.WebURL)
}

