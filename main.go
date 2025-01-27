package main

import (
	"github.com/shn27/Test/work_mysql"
)

func main() {
	err := work_mysql.MySqlCmdTest.Execute()
	if err != nil {
		return
	}
}
