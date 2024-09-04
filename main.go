package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

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

	// HTML template for the table
	const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Diagnostic Results</title>
    <style>
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid black; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>Diagnostic Results</h1>
    <pre>{{ . }}</pre>
</body>
</html>
`

	// Create a new template and parse the HTML string into it
	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		fmt.Println("Error creating template:", err)
		return
	}

	// Create a buffer to capture the final HTML output
	var htmlBuf bytes.Buffer

	// Execute the template with the table output
	err = tmpl.Execute(&htmlBuf, tableOutput)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// Save the HTML content to a file
	err = os.WriteFile("output.html", htmlBuf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing to output.html:", err)
	} else {
		fmt.Println("Table output saved to output.html")
	}
}
