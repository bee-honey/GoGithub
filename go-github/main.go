package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Struct to parse the commit response
type Commit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commit"`
}

// Fetch commits from the GitHub API
func fetchCommits(repoOwner, repoName string) ([]Commit, error) {
	// GitHub API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching commits: %s", resp.Status)
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// API handler to return commits as JSON
func commitsHandler(w http.ResponseWriter, r *http.Request) {
	repoOwner := "bee-honey"
	repoName := "GoGithub"

	commits, err := fetchCommits(repoOwner, repoName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commits)
}

// Serve static files and run the server
func main() {
	http.HandleFunc("/api/commits", commitsHandler)

	// Serve static files (HTML, JS, CSS)
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
