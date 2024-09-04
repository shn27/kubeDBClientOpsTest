package main

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
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

const maxLinesPerOutput = 1

func main() {
	// Example data
	results := []DiagnosticResult{
		{
			CheckType: "Type1",
			Timestamp: metav1.Now(),
			Outputs: []DiagnosticOutput{
				{Message: "Output 1", Data: []byte("Data 1")},
				{Message: "Output 2", Data: []byte("Data 2")},
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

	// Convert buffer to string
	tableOutput := buf.String()

	// Print the string variable containing the table
	fmt.Println(tableOutput)

	err := os.WriteFile("README.md", []byte(tableOutput), 0644)
	if err != nil {
		fmt.Println("Error writing to README.md:", err)
	} else {
		fmt.Println("Table output saved to README.md")
	}
}
