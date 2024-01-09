package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/google/go-github/github"
)

const (
	owner    = "poteto"
	repo     = "hiring-without-whiteboards"
	basePath = "README.md"
)

func main() {
	fmt.Println("Hello, world!")

	var file FileData
	file.getContents()
	file.parseFile()
	sort.Strings(file.Locations)

	writeFile(file.Locations, "locations.json")
	writeFile(file.Jobs, "jobs.json")
}

type FileData struct {
	File          string   `json:"file"`
	Jobs          []Job    `json:"jobs"`
	Locations     []string `json:"locations"`
	TempLocations []string `json:"tempLocations"`
}

type Job struct {
	Name      string   `json:"name"`
	Link      string   `json:"link"`
	Locations []string `json:"locations"`
	Process   string   `json:"process"`
}

func (f *FileData) getContents() {
	context := context.Background()
	client := github.NewClient(nil)

	fmt.Println("Getting contents...")

	// fileContent, directoryContent, resp, err := client.Repositories.GetContents(context, owner, repo, basePath, nil)
	fileContent, _, _, err := client.Repositories.GetContents(context, owner, repo, basePath, nil)

	if err != nil {
		fmt.Println("ERROR GETTING CONTENTS", err)
	}

	if fileContent == nil {
		fmt.Println("ERROR: File content is nil")
	}

	content, err := fileContent.GetContent()

	if err != nil {
		fmt.Println("ERROR GETTING FILE CONTENT", err)
	}

	f.File = content
}

func (f *FileData) parseFile() {
	fmt.Println("Outputting file...")

	x := NewSet()
	lines := strings.Split(f.File, "\n")
	lineStart := getStartingLine(lines)

	for i := lineStart + 1; i < len(lines); i++ {
		line := lines[i]

		if len(line) > 0 && line[0:1] == "-" {
			pattern := `^\- \[(.*?)\]\((.*?)\) \| (.*?) \| (.*)$`
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(line)

			if len(matches) > 0 {
				locs := []string{}
				locations := strings.Split(matches[3], "/")

				// Trim spaces and split by semicolon
				// then creates a new array of locations
				for _, loc := range locations {
					trimmedLocation := strings.TrimSpace(loc)
					s := strings.Split(trimmedLocation, ";")

					for _, loc := range s {
						trimmedLocation := strings.TrimSpace(loc)
						locs = append(locs, trimmedLocation)

						if !x.Exists(trimmedLocation) {
							x.Add(trimmedLocation)
							f.Locations = append(f.Locations, trimmedLocation)
						}
					}
				}

				// Create new job
				job := Job{
					Name:      matches[1],
					Link:      matches[2],
					Locations: locs,
					Process:   matches[4],
				}

				f.Jobs = append(f.Jobs, job) // Append job to f.Jobs
			}
		}
	}
}

func getStartingLine(lines []string) int {
	var lineStart int

	// finds the line where the table of contents starts
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if len(line) > 0 && line[0:3] == "---" {
			lineStart = i
			break
		}
	}

	return lineStart
}

func writeFile(data interface{}, fileName string) {
	fmt.Println("Writing file...")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error occurred during file creation. Error: %s", err.Error())
	}

	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error occurred during file writing. Error: %s", err.Error())
	}
}

// func outputJson(lines interface{}) {
// 	fmt.Println("Outputting JSON...")

// 	jsonData, err := json.MarshalIndent(lines, "", "  ")
// 	if err != nil {
// 		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
// 	}

// 	fmt.Println(string(jsonData))
// }
