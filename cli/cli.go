package main

import (
	"bufio"
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

	client := &http.Client{Timeout: 6 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)

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
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	args := strings.Split(text, " ")
	if len(args) > 2 {
		return "", "", fmt.Errorf("input should look like <owner>/<repoName> or <owner> <repoName>")
	}
	if len(args) == 1 {
		args = strings.Split(args[0], "/")
	}
	if len(args) != 2 {
		return "", "", fmt.Errorf("input should look like <owner>/<repoName> or <owner> <repoName>")
	}
	return args[0], args[1], nil

}

func (r *RepoInfo) printInfo() {
	fmt.Printf("%s:\n", r.Name)
	fmt.Printf("	Description: %s\n", r.Description)
	fmt.Printf("	Stars      : %d\n", r.StarsCount)
	fmt.Printf("	Forks      : %d\n", r.ForksCount)
	fmt.Printf("	Created at : %s\n", r.CreatedAt.Format("02.01.2006"))
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
	repo_info.printInfo()
}
