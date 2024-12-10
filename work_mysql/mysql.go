package work_mysql

import (
	"fmt"
	"strings"
)

type QueryStats struct {
	Time         string  // Timestamp from the # Time line
	QueryTime    float64 // Time the query took to execute
	LockTime     float64 // Time the query was locked
	RowsSent     int     // Number of rows sent
	RowsExamined int     // Number of rows examined
	Query        string  // Full query text
}

func analyzeSlowLogReverse(slowLog string) ([]QueryStats, error) {
	lines := strings.Split(slowLog, "\n")

	var stats []QueryStats
	var currentBlock strings.Builder
	queryCount := 0

	// Process lines in reverse order
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])

		// Detect a new query block starting with "# Time:"
		if strings.HasPrefix(line, "# Time:") {
			if currentBlock.Len() > 0 {
				// Parse the current block and prepend the result
				queryStat, err := parseQueryBlock(currentBlock.String())
				if err == nil {
					stats = append([]QueryStats{queryStat}, stats...)
					queryCount++
				}
				currentBlock.Reset()
			}
		}

		// Add the current line to the block
		currentBlock.WriteString(line + "\n")

		// Stop if we have enough queries
		if queryCount >= 3 {
			break
		}
	}

	// Process the last block if it exists
	if queryCount < 3 && currentBlock.Len() > 0 {
		queryStat, err := parseQueryBlock(currentBlock.String())
		if err == nil {
			stats = append([]QueryStats{queryStat}, stats...)
		}
	}
	return stats, nil
}

func parseQueryBlock(block string) (QueryStats, error) {
	var stat QueryStats
	lines := strings.Split(block, "\n")
	var queryBuilder strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse # Time line
		if strings.HasPrefix(line, "# Time:") {
			stat.Time = strings.TrimPrefix(line, "# Time: ")
		} else if strings.HasPrefix(line, "# Query_time:") {
			// Parse query metadata
			_, err := fmt.Sscanf(line, "# Query_time: %f  Lock_time: %f Rows_sent: %d  Rows_examined: %d",
				&stat.QueryTime, &stat.LockTime, &stat.RowsSent, &stat.RowsExamined)
			if err != nil {
				return stat, fmt.Errorf("failed to parse query metadata: %w", err)
			}
		} else if strings.HasPrefix(line, "SET timestamp=") || strings.HasPrefix(line, "# User@Host:") {
			// Skip these lines as they are not part of the query
			continue
		} else if line != "" {
			// Append lines that are part of the query
			queryBuilder.WriteString(line + " ")
		}
	}

	// Assign the constructed query
	stat.Query = strings.TrimSpace(queryBuilder.String())
	if stat.Query == "" {
		return stat, fmt.Errorf("query text is empty")
	}

	return stat, nil
}

func printSlowQuery(slowQuery string) {
	stats, err := analyzeSlowLogReverse(slowQuery)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Last 3 Slow Queries:")
	for i, stat := range stats {
		fmt.Printf("Query #%d:\n", i+1)
		fmt.Printf("  Time: %s \n", stat.Time)
		fmt.Printf("  Query Time: %.2f sec\n", stat.QueryTime)
		fmt.Printf("  Lock Time: %.2f sec\n", stat.LockTime)
		fmt.Printf("  Rows Sent: %d\n", stat.RowsSent)
		fmt.Printf("  Rows Examined: %d\n", stat.RowsExamined)
		fmt.Printf("  Query: %s\n", stat.Query)
		fmt.Println()
	}
}
