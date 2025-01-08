package work_mysql

import "github.com/spf13/cobra"

var MySqlCmdTest = &cobra.Command{
	Use:   "mysql",
	Short: "mysql test",
	Long:  "mysql test long",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := GetMysqlClient()
		if err != nil {
			return
		}
	},
}
