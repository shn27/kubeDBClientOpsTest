package work

import (
	"bytes"
	"html/template"
	"os"

	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Render the subtable for the Output column
func renderSubTable(outputs []DiagnosticOutput) string {
	var subBuf bytes.Buffer
	subtable := tablewriter.NewWriter(&subBuf)
	subtable.SetHeader([]string{"Message", "Data"})
	subtable.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})

	// Add each DiagnosticOutput entry as a row in the subtable
	for _, output := range outputs {
		subtable.Append([]string{output.Message, string(output.Data)})
	}
	subtable.Render()
	return subBuf.String()
}

// Render the main table with subtables
func renderMainTable(results []DiagnosticResult) string {
	var buf bytes.Buffer

	// Define the HTML template
	const tpl = `
	<table border="1">
		<tr>
			<th>Check Type</th>
			<th>Timestamp</th>
			<th>Output</th>
		</tr>
		{{- range . }}
		<tr>
			<td>{{ .CheckType }}</td>
			<td>{{ .Timestamp }}</td>
			<td><pre>{{ .SubTable }}</pre></td>
		</tr>
		{{- end }}
	</table>
	`

	// Prepare data by converting each result's output into a subtable string
	var tableData []struct {
		CheckType, Timestamp, SubTable string
	}

	for _, result := range results {
		tableData = append(tableData, struct {
			CheckType, Timestamp, SubTable string
		}{
			CheckType: result.CheckType,
			Timestamp: result.Timestamp.String(),
			SubTable:  renderSubTable(result.Outputs),
		})
	}

	// Execute template with embedded subtables
	t := template.Must(template.New("mainTable").Parse(tpl))
	if err := t.Execute(&buf, tableData); err != nil {
		panic(err)
	}

	return buf.String()
}

func Table() {
	// Sample data
	var data []byte
	for i := 0; i < 5; i++ {
		data = append([]byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry."), data...)
	}

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

	// Render the main table with embedded subtables
	htmlOutput := renderMainTable(results)
	os.WriteFile("output.html", []byte(htmlOutput), 0644)

	// Optional: Write the HTML output to a file
	// ioutil.WriteFile("output.html", []byte(htmlOutput), 0644)
}
