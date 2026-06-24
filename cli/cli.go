package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type RepoInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StarsCount  int       `json:"stargazers_count"`
	ForksCount  int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func getRepoInfo(owner, repoName string) (*RepoInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "my-cli-tool")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var repo_info RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&repo_info); err != nil {
		return nil, err
	}

	return &repo_info, nil
}

func parseInput() (string, string, error) {
	var args []string = os.Args[1:]
	if len(args) > 2 {
		return "", "", fmt.Errorf("input arguments should look like <owner>/<repoName> or <owner> <repoName>")
	}
	if len(args) == 1 {
		args = strings.Split(args[0], "/")
	}
	if len(args) != 2 {
		return "", "", fmt.Errorf("input arguments should look like <owner>/<repoName> or <owner> <repoName>")
	}
	return args[0], args[1], nil

}

func (r *RepoInfo) String() string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("%s:\n", r.Name))
	str.WriteString(fmt.Sprintf("	Description: %s\n", r.Description))
	str.WriteString(fmt.Sprintf("	Stars      : %d\n", r.StarsCount))
	str.WriteString(fmt.Sprintf("	Forks      : %d\n", r.ForksCount))
	str.WriteString(fmt.Sprintf("	Created at : %s", r.CreatedAt.Format("02.01.2006")))
	return str.String()
}

func main() {
	owner, repoName, err := parseInput()
	for err != nil {
		fmt.Println(err)
		owner, repoName, err = parseInput()
	}
	repo_info, err := getRepoInfo(owner, repoName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(repo_info)
}
