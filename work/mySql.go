package work

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mysql"
)

func getMysqlClient() (*mysql.Client, error) {
	kbClient, err := getKBClient()
	if err != nil {
		fmt.Println("failed to get k8s client", err)
		return nil, err
	}
	db := &api.MySQL{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql",
			Namespace: "monitoring",
		},
	}
	kubeDBClient, err := mysql.NewKubeDBClientBuilder(kbClient, db).
		WithPod("mysql-0").
		//WithCred("root:lK!U7bOqp1SdlUOQ").
		GetMySQLClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	return kubeDBClient, nil
}
