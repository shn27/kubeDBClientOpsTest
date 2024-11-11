package work

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mysql"
)

func PrimaryServiceDNS(db *api.MySQL) string {
	return fmt.Sprintf("%v.%v.svc", db.ServiceName(), db.Namespace)
}

func getMysqlClient() (*mysql.Client, error) {
	kbClient, err := getKBClient()
	if err != nil {
		fmt.Println("failed to get k8s client", err)
		return nil, err
	}
	ref := kmapi.ObjectReference{
		Name:      "mysql",
		Namespace: "monitoring",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1alpha2",
		Group:   "kubedb.com",
		Kind:    "MySQL",
	}

	obj, err := GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &api.MySQL{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}
	kubeDBClient, err := mysql.NewKubeDBClientBuilder(kbClient, db).
		WithPod("mysql-0").
		WithURL(PrimaryServiceDNS(db)).
		GetMySQLClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	return kubeDBClient, nil
}
