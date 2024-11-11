package work

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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

const maxLinesPerOutput = 5

func GetMarkdown() string {
	// Example data

	var data []byte
	for i := 0; i < 5; i++ {
		data = append([]byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry."), data...)
	}

	//fmt.Println(string(data))

	results := []DiagnosticResult{
		{
			CheckType: "InspectLogs",
			Timestamp: metav1.Now(),
			Outputs: []DiagnosticOutput{
				{Message: "Log 1", Data: data},
				{Message: "Log 2", Data: data},
				{Message: "Log 3", Data: data},
				{Message: "Log 4", Data: data},
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

		//subtable.SetRowLine(true)
		//subtable.SetAutoMergeCells(false)
		//subtable.SetAutoWrapText(false)

		subtable := tablewriter.NewWriter(&subBuf)
		subtable.SetHeader([]string{"Message", "Data"})
		subtable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})

		for i, output := range result.Outputs {
			if i >= maxLinesPerOutput {
				break
			}
			subtable.Append([]string{output.Message, string(output.Data)})
		}

		subtable.Render()

		subtableString := mdToHTML(subBuf.Bytes())

		// Check if more lines exist and append "..." to indicate more content
		//if len(result.Outputs) > maxLinesPerOutput {
		//	subtableString += "\n... (More data available)"
		//}

		// Add the DiagnosticResult to the main table
		table.Append([]string{
			result.CheckType,
			result.Timestamp.Time.Format("2006-01-02 15:04:05"),
			string(subtableString),
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
		return ""
	}

	// Create a buffer to capture the final HTML output
	var htmlBuf bytes.Buffer

	// Execute the template with the table output
	err = tmpl.Execute(&htmlBuf, tableOutput)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return ""
	}

	// Save the HTML content to a file
	err = os.WriteFile("output.html", mdToHTML(htmlBuf.Bytes()), 0644)
	if err != nil {
		fmt.Println("Error writing to output.html:", err)
	} else {
		fmt.Println("Table output saved to output.html")
	}
	return ""
}

func TestTableWriter() {
	result := DiagnosticResult{
		CheckType: "Network Check",
		Timestamp: metav1.Time{Time: time.Now()},
		Outputs: []DiagnosticOutput{
			{Message: "Ping Success", Data: []byte("Response Time: 32ms")},
			{Message: "Traceroute", Data: []byte("Hops: 15")},
		},
	}

	// Create a buffer for the subtable
	var subBuf bytes.Buffer
	subtable := tablewriter.NewWriter(&subBuf)
	subtable.SetHeader([]string{"Message", "Data"})
	subtable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})

	// Add each DiagnosticOutput entry as a row in the subtable
	for _, output := range result.Outputs {
		subtable.Append([]string{output.Message, string(output.Data)})
	}

	// Render the subtable and capture it as a single string
	subtable.Render()
	subtableStr := mdToHTML(subBuf.Bytes())

	// Now create the main table
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Check Type", "Timestamp", "Output"})

	// Append a single row to the main table with the subtable in the Output column
	table.Append([]string{
		result.CheckType,
		result.Timestamp.String(),
		string(subtableStr),
	})

	table.Render()
	fmt.Println(buf.String())
	tableOutput := buf.String()
	// Print the main table with embedded subtable

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

	// If needed, save to an HTML file
	// Example: ioutil.WriteFile("output.html", buf.Bytes(), 0644)
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
