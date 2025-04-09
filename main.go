package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
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
	log.Printf("Downloading from %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

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

type Response struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload *struct {
		Action       string `json:"action"`
		Size         int    `json:"size"`
		DistinctSize int    `distinct_size:"distinct_size"`
		RefType      string `json:"ref_type"`
		Commits      []struct {
			Message string `json:"message"`
			URL     string `json:"url"`
		} `json:"commits"`
	} `json:"payload"`
	Repo *struct {
		Name string `json:"name"`
	} `json:"repo"`
	When string `json:"created_at"`
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	argument := os.Args[1]
	jsonData := getJSON(argument)

	jsonFile, err := os.Create("github.json")
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()
	jsonFile.Write(jsonData)

	var GitHubEvents []Response
	json.Unmarshal(jsonData, &GitHubEvents)
	newData, _ := json.MarshalIndent(GitHubEvents, "", " ")

	newFile, err := os.Create("ngithub.json")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	newFile.Write(newData)
}
