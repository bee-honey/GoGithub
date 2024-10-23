package main

import (
	"encoding/json"
	"fmt"
	"go-github/entities"
	"log"
	"net/http"
	"strings"
)

// Fetch commits from the GitHub API
func fetchCommits(repoOwner, repoName string) ([]entities.CICommit, error) {
	// GitHub API URL
	// TBD: Need to be able to change the reponame with a dropdown
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching commits: %s", resp.Status)
	}

	var apiCommits []struct {
		Sha    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiCommits); err != nil {
		return nil, err
	}

	var commits []entities.CICommit
	for _, c := range apiCommits {
		commit := entities.CICommit{
			Sha:     c.Sha,
			Message: c.Commit.Message,
			Author:  c.Commit.Author.Name,
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// Fetch workflow runs from GitHub Actions API
func fetchWorkflowRuns(repoOwner, repoName string) (map[string]entities.WorkflowRun, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching workflow runs: %s", resp.Status)
	}

	var result struct {
		WorkflowRuns []entities.WorkflowRun `json:"workflow_runs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Create a map of commit SHA to workflow run for quick lookup
	workflowMap := make(map[string]entities.WorkflowRun)
	for _, run := range result.WorkflowRuns {
		workflowMap[run.HeadSha] = run
	}

	return workflowMap, nil
}

func commitsHandler(w http.ResponseWriter, r *http.Request) {
	repoOwner := "bee-honey"
	repoName := "AbbtKeerthy"

	//Commits info
	commits, err := fetchCommits(repoOwner, repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Workflow info
	workflowMap, err := fetchWorkflowRuns(repoOwner, repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Both workflows and commits needs to match here
	for i, commit := range commits {
		if workflowRun, found := workflowMap[commit.Sha]; found {
			// Set the workflow status on the commit
			commits[i].Status = strings.Title(workflowRun.Status)
			commits[i].Conclusion = strings.Title(workflowRun.Conclusion)
		} else {
			commits[i].Status = "No Workflow Run"
			commits[i].Conclusion = "No"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commits)
}

// Serve static files and run the server
func main() {
	http.HandleFunc("/api/commits", commitsHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// port := os.Getenv("PORT")
	var port = "8085"
	// if port == "" {
	// 	port = "8085"
	// }

	fmt.Println("Server started at port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
