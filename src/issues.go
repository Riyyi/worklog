package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type IssueData struct {
	totalIssues []issueResponse
}

// Constructor
func MakeIssueData() IssueData {
	return IssueData {}
}

func (self *IssueData) GenerateIssuesFile(path string) {
	// Safety check
	verify(filepath.Base(path) != "worklog.org", "protected from overwriting worklog file!")

	// Get all issues of the active sprint, and its subtasks
	self.fetchIssues(0)

	// Store the issues
	outputFile, err := os.Create(path)
	assert(err)
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	self.writeIssues(writer)

	// Clear the results
	self.totalIssues = self.totalIssues[:0]
}

// -----------------------------------------

type issuesResponse struct {
	StartAt    int              `json:"startAt"`
	MaxResults int              `json:"maxResults"`
	Total      int              `json:"total"`
	Issues     []issueResponse  `json:"issues"`
}

type issueResponse struct {
	Key    string          `json:"key"`
	Fields fieldsResponse  `json:"fields"`
}

type fieldsResponse struct {
	Summary  string            `json:"summary"`
	SubTasks []subTaskResponse `json:"subtasks"`
}

type subTaskResponse struct {
	Key    string                 `json:"key"`
	Fields subTaskFieldsResponse  `json:"fields"`
}

type subTaskFieldsResponse struct {
	Summary string `json:"summary"`
}

// -----------------------------------------

func (self *IssueData) fetchIssues(startAt int) error {
	var url string = baseUrl + "/rest/api/2/search"
	var maxResults int = 100
    data := map[string]interface{}{
		"fields": []string{ "key", "summary", "subtasks" },
		"jql": `project = "` + projectName + `" AND sprint IN openSprints() AND issuetype != "Sub-task" ORDER BY created ASC`,
		"maxResults": maxResults,
		"startAt": startAt,
    }

	body, err := Request(url, data, 200) // "OK"
	if err != nil { return err }

	var result issuesResponse
	err = json.Unmarshal(body, &result)
	if err != nil { fmt.Println("NOPE!"); return err }

	// Add fetched issues
	self.totalIssues = append(self.totalIssues, result.Issues...)

	// Pagination, if more results
	if startAt + maxResults < result.Total  {
		self.fetchIssues(startAt + maxResults)
	}

	return nil
}

func (self* IssueData) writeIssues(writer *bufio.Writer) {
	// Issues
	writer.WriteString(".\n")
	var count int = len(self.totalIssues)
	for i, issue := range(self.totalIssues) {
		if i == count - 1 {
			writer.WriteString("└")
		} else {
			writer.WriteString("├")
		}
		writer.WriteString("── ")
		writer.WriteString(issue.Key)
		writer.WriteString("        ")
		writer.WriteString(issue.Fields.Summary)
		writer.WriteString("\n")

		// Subtasks
		var subtaskCount int = len(issue.Fields.SubTasks)
		for j, subtask := range(issue.Fields.SubTasks) {
			// Last issue
			if i == count - 1 {
				writer.WriteString("    ")
			} else {
				writer.WriteString("│   ")
			}
			// Last subtask
			if j == subtaskCount - 1 {
				writer.WriteString("└")
			} else {
				writer.WriteString("├")
			}
			writer.WriteString("── ")
			writer.WriteString(subtask.Key)
			writer.WriteString("    ")
			writer.WriteString(subtask.Fields.Summary)
			writer.WriteString("\n")
		}
	}
}
