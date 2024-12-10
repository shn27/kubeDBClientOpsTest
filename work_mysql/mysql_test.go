package work_mysql

import (
	"fmt"
	utils "github.com/shn27/Test/utils"
	exec_util "kmodules.xyz/client-go/tools/exec"
	"kubedb.dev/apimachinery/apis/kubedb"
	"testing"
)

func Test(t *testing.T) {
	cmd := "cat /var/lib/mysql/mysql-0-slow.log"
	command := exec_util.Command("bash", "-c", cmd)
	container := exec_util.Container(kubedb.MySQLContainerName)
	options := []func(options *exec_util.Options){
		command,
		container,
	}
	config, err := utils.GetInClusterConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	pod, err := utils.GetPod("demo", "mysql-0")
	if err != nil {
		fmt.Println(err)
	}
	slowQueryLog, err := exec_util.ExecIntoPod(config, pod, options...)
	if err != nil {
		fmt.Println(err)
		return
	}
	printSlowQuery(slowQueryLog)
}
