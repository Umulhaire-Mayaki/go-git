package main

import (
	"fmt"
	"log"
	"os"


	//gitlab "github.com/xanzy/go-gitlab"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func main() {
	// GitLab Project Details
	projectNamespace := "umayaki" 
	projectName := "clusterimagesets"          
	sourceBranch := "feature/mr-automation-changes"   
	targetBranch := "master"             
	commitMessage := "Automated commit message"
	mergeRequestTitle := "Automated Merge Request"
	mergeRequestDescription := "This MR was created automatically."

	// Step 1: Clone or open the repository
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
	commit, err := worktree.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		log.Fatalf("Failed to commit changes: %v", err)
	}
	fmt.Printf("Commit created: %s\n", commit.String())

	// Step 4: Push the changes to remote
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		log.Fatalf("Failed to push changes: %v", err)
	}
	fmt.Println("Changes pushed successfully.")

	// Step 5: Create a Merge Request using GitLab API
	token := os.Getenv("6-9yBqwxdN75XfBvconw") // Ensure your GitLab token is set as an environment variable
	if token == "" {
		log.Fatalf("GITLAB_TOKEN environment variable is not set")
	}

	gitlabClient, err := gitlab.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	projectID := fmt.Sprintf("%s/%s", projectNamespace, projectName)
	mergeRequest, _, err := gitlabClient.MergeRequests.CreateMergeRequest(projectID, &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String(mergeRequestTitle),
		SourceBranch: gitlab.String(sourceBranch),
		TargetBranch: gitlab.String(targetBranch),
		Description:  gitlab.String(mergeRequestDescription),
	})
	if err != nil {
		log.Fatalf("Failed to create MR: %v", err)
	}

	fmt.Printf("Merge Request created: %s\n", mergeRequest.WebURL)
}
