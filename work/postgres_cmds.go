package work

import (
	work_postgres "github.com/shn27/Test/work_postgres"
	"github.com/spf13/cobra"
)

var pgCmdTest = &cobra.Command{
	Use:   "pgCmdTest",
	Short: "This is a simple CLI application",
	Long:  `Test postgres cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		work_postgres.TestPostgresServerStatus()
	},
}
