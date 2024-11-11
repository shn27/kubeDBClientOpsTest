package work_postgres

import (
	"github.com/spf13/cobra"
)

var PgCmdTest = &cobra.Command{
	Use:   "pgCmdTest",
	Short: "This is a simple CLI application",
	Long:  `Test postgres cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		TestPostgresServerStatus()
	},
}
