package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

//go:embed help.txt
var helpText []byte

func printHelp() {
	fmt.Printf("%s\n", helpText)
}

func getAPIURL(username string) string {
	return fmt.Sprintf("https://api.github.com/users/%s/events", username)
}

func getJSON(username string) []byte {
	url := getAPIURL(username)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Response status: %v", resp.Status)
		fmt.Printf("Could not access user events.\n")
		panic(errors.New("bad HTTP response code"))
	}

	scanner := bufio.NewScanner(resp.Body)
	var buf []byte
	var maxBufSize int = 512000
	scanner.Buffer(buf, maxBufSize)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	data := scanner.Bytes()
	return data
}

type GHResponse struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload struct {
		Action       string `json:"action"`
		Size         int    `json:"size"`
		DistinctSize int    `distinct_size:"distinct_size"`
		RefType      string `json:"ref_type"`
		Commits      []struct {
			Message string `json:"message"`
			URL     string `json:"url"`
		} `json:"commits"`
		PullRequest struct {
			Title string `json:"title"`
		} `json:"pull_request"`
		Issue struct {
			Title string `json:"title"`
		} `json:"issue"`
	} `json:"payload"`
	Repo *struct {
		Name string `json:"name"`
	} `json:"repo"`
	When string `json:"created_at"`
}

// Describes a GHResponse with a short line of text
func (r GHResponse) Describe() string {
	switch r.Type {
	case "PushEvent":
		return fmt.Sprintf("Pushed %d commits to %s", r.Payload.Size, r.Repo.Name)
	case "WatchEvent":
		if r.Payload.Action == "started" {
			return fmt.Sprintf("Is now watching %s", r.Repo.Name)
		} else {
			return fmt.Sprintf("Stopped watching %s", r.Repo.Name) //?
		}
	case "PublicEvent":
		return fmt.Sprintf("Made the repo %s public", r.Repo.Name)
	case "PullRequestEvent":
		if r.Payload.Action == "opened" {
			return fmt.Sprintf("Opened pull request \"%s\" in repo %s", r.Payload.PullRequest.Title, r.Repo.Name)
		} else if r.Payload.Action == "closed" {
			return fmt.Sprintf("Closed pull request \"%s\" in repo %s", r.Payload.PullRequest.Title, r.Repo.Name)
		}
	case "IssuesEvent":
		if r.Payload.Action == "opened" {
			return fmt.Sprintf("Opened issue \"%s\" in repo %s", r.Payload.Issue.Title, r.Repo.Name)
		} else if r.Payload.Action == "closed" {
			return fmt.Sprintf("Closed issue \"%s\" in repo %s", r.Payload.Issue.Title, r.Repo.Name)
		}
	case "IssueCommentEvent":
		return fmt.Sprintf("Commented on issue \"%s\" in repo %s", r.Payload.Issue.Title, r.Repo.Name)
	case "CreateEvent":
		return fmt.Sprintf("Created repo %s", r.Repo.Name)
	default:
		fmt.Printf("This event has an invalid type (the creator of this app didn't properly understand the API response schema. Please report this on github)\n")
		log.Printf("Unknown event type %s\n", r.Type)
	}
	return ""
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	argument := os.Args[1]
	jsonData := getJSON(argument)

	var GitHubEvents []GHResponse
	json.Unmarshal(jsonData, &GitHubEvents)

	fmt.Printf("Recent activity of %s\n", argument)

	for _, ghr := range GitHubEvents {
		if ghr.Describe() != "" {
			fmt.Printf(" - %s\n", ghr.Describe())
		}
	}
}
