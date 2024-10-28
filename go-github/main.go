package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-github/entities"
	"go-github/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	repoOwner    = "bee-honey"
	repoName     = "CICDAbt"
	workflowFile = "release.yml"
)

func getAPIBaseURL(repoOwner, repoName, endpoint string) string {
	const apiBaseURLPattern = "https://api.github.com/repos/%s/%s"
	baseURL := fmt.Sprintf(apiBaseURLPattern, repoOwner, repoName)
	return fmt.Sprintf("%s/%s", baseURL, endpoint)
}

// Fetch commits from the GitHub API
func fetchCommits(repoOwner, repoName string) ([]entities.CICommit, error) {
	// GitHub API URL
	// TBD: Need to be able to change the reponame with a dropdown
	url := getAPIBaseURL(repoOwner, repoName, "commits")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching commits: %s", resp.Status)
	}

	var apiCommits []entities.GitCommit

	if err := json.NewDecoder(resp.Body).Decode(&apiCommits); err != nil {
		return nil, err
	}

	var commits []entities.CICommit
	for _, c := range apiCommits {
		readableDate := utils.FormatToReadableDate(c.Commit.Author.Date)
		commit := entities.CICommit{
			Sha:     c.Sha,
			Message: c.Commit.Message,
			Author:  c.Commit.Author.Name,
			Date:    readableDate,
		}
		commits = append(commits, commit)
	}

	return commits, nil
}

// Fetch workflow runs from GitHub Actions API
func fetchWorkflowRuns(repoOwner, repoName string) (map[string]entities.WorkflowRun, error) {

	url := getAPIBaseURL(repoOwner, repoName, "actions/runs")

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

func triggerRelease(w http.ResponseWriter, r *http.Request) {

	githubToken := os.Getenv("GITHUB_TOKEN") //@myself, make sure this is safe if deployed to AWS

	if githubToken == "" {
		http.Error(w, "GITHUB_TOKEN not set", http.StatusInternalServerError)
		return
	}

	// API to trigger the workflow
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches", repoOwner, repoName, workflowFile)
	payload := map[string]string{
		"ref": "main", //For now lets only worry about the main, hotfixes etc to be delt later on
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode JSON payload", http.StatusInternalServerError)
		return
	}

	// trigger the workflow
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to execute request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("GitHub API returned status %d", resp.StatusCode)
		http.Error(w, "GitHub API request failed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Workflow triggered successfully!")
}

func commitsHandler(w http.ResponseWriter, r *http.Request) {

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
	http.HandleFunc("/api/release", triggerRelease)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	var port = "8085"

	fmt.Println("Server started at port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
