package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
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
	getContents()
}

func getContents() {
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

	parseFile(content)

	// fmt.Printf("File Contents: %#v\n", fileContent)
	// fmt.Printf("Directory Contents: %#v\n", directoryContent)
	// fmt.Printf("Resp: %#v\n", resp)
}

func (j Job) String() {
	fmt.Printf("%s - %s - %s - %s\n", j.Name, j.Link, j.Locations, j.Process)
}

func getStartingLine(lines []string) int {
	var lineStart int

	// finds the line where the table of contents starts
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if len(line) > 0 && line[0:3] == "---" {
			lineStart = i
		}
	}

	return lineStart
}

func parseLine(line string) *Job {
	job := Job{}

	pattern := `^\- \[(.*?)\]\((.*?)\) \| (.*?) \| (.*)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 0 {
		locations := strings.Split(matches[3], "/")

		for i, loc := range locations {
			locations[i] = strings.TrimSpace(loc)
		}

		job.Name = matches[1]
		job.Link = matches[2]
		job.Locations = locations
		job.Process = matches[4]

		// job.String()

		return &job
	}

	return nil
}

func parseFile(content string) {
	fmt.Println("Outputting file...")

	lines := strings.Split(content, "\n")
	lineStart := getStartingLine(lines)

	var jobs []Job

	for i := lineStart + 1; i < len(lines); i++ {
		line := lines[i]

		if len(line) > 0 && line[0:1] == "-" {
			job := parseLine(line)

			if job != nil {
				jobs = append(jobs, *job)
			}
		}
	}

	locations := []string{"San Francisco", "Remote"}
	writeFile(locations, "locations.json")
	writeFile(jobs, "jobs.json")
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

// if len(line) > 0 && line[0:1] == "-" {
// 	fmt.Println(line)
// 	break
// }
// break

// if line[0:1] == "#" {
// 	fmt.Println(line)
// }
// for _, line := range strings.Split(content, "\n") {
// 	fmt.Println(line)
// }
