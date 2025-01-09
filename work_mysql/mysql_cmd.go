package work_mysql

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var MySqlCmdTest = &cobra.Command{
	Use:   "mysql",
	Short: "mysql test",
	Long:  "mysql test long",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			time.Sleep(time.Second * 5)
			client, err := GetMysqlClient()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = mysqlQueryTest(client)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("==============mysql test success================")
		}
	},
}
