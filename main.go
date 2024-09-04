package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DiagnosticResult struct {
	CheckType string
	Timestamp metav1.Time
	Outputs   []DiagnosticOutput
}

type DiagnosticOutput struct {
	Message string
	Data    []byte
}

const maxLinesPerOutput = 10

func main() {
	// Example data
	results := []DiagnosticResult{
		{
			CheckType: "Type1",
			Timestamp: metav1.Now(),
			Outputs: []DiagnosticOutput{
				{Message: "Output 1", Data: []byte("Data 1")},
				{Message: "Output 2", Data: []byte("Data 2")},
				{Message: "Output 3", Data: []byte("Data 3")},
				{Message: "Output 4", Data: []byte("Data 4")},
				{Message: "Output 5", Data: []byte("Data 5")},
				{Message: "Output 6", Data: []byte("Data 6")},
				{Message: "Output 7", Data: []byte("Data 7")},
				{Message: "Output 8", Data: []byte("Data 8")},
				{Message: "Output 9", Data: []byte("Data 9")},
				{Message: "Output 10", Data: []byte("Data 10")},
				{Message: "Output 11", Data: []byte("Data 11")},
				{Message: "Output 12", Data: []byte("Data 12")},
			},
		},
	}

	// Create a buffer to capture the table output
	var buf bytes.Buffer

	// Create the main table and render it to the buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Check Type", "Timestamp", "Outputs"})

	for _, result := range results {
		// Create a subtable for the DiagnosticOutput
		var subBuf bytes.Buffer
		subtable := tablewriter.NewWriter(&subBuf)
		subtable.SetHeader([]string{"Message", "Data"})

		for i, output := range result.Outputs {
			if i >= maxLinesPerOutput {
				break
			}
			subtable.Append([]string{output.Message, string(output.Data)})
		}

		subtable.Render()
		subtableString := subBuf.String()

		// Check if more lines exist and append "..." to indicate more content
		if len(result.Outputs) > maxLinesPerOutput {
			subtableString += "\n... (More data available)"
		}

		// Add the DiagnosticResult to the main table
		table.Append([]string{
			result.CheckType,
			result.Timestamp.Time.Format("2006-01-02 15:04:05"),
			subtableString,
		})
	}

	// Render the main table to the buffer
	table.Render()

	// Save the buffer content to a file
	err := os.WriteFile("output.md", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	} else {
		fmt.Println("Table output saved to output.md")
	}
}
