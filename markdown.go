package main

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func GenerateMarkDown(result []DiagnosticResult) {
	data := [][]string{
		[]string{"A", "The Good", "500"},
		[]string{"B", "The Very very Bad Man", "288"},
		[]string{"C", "The Ugly", "120"},
		[]string{"D", "The Gopher", "800"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Sign", "Rating"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output

	table1 := tablewriter.NewWriter(os.Stdout)
	table1.SetHeader([]string{"Name", "Result", "Time", "Link"})

	for _, v := range result {
		x := v.Timestamp.String()
		table1.Append([]string{v.CheckType, "ok", x, "abc"})
	}
	table1.Render()
}
