package work

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/postgres"
)

func GetPostgresClient() (*postgres.Client, error) {
	kbClient, err := getKBClient()
	if err != nil {
		fmt.Println("failed to get k8s client", err)
		return nil, err
	}
	ref := kmapi.ObjectReference{
		Name:      "postgres",
		Namespace: "monitoring",
	}
	gvk := metav1.GroupVersionKind{
		Version: "v1alpha2",
		Group:   "kubedb.com",
		Kind:    "Postgres",
	}

	obj, err := getK8sObject(kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := dbapi.Postgres{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}

	kubeDBClient, err := postgres.NewKubeDBClientBuilder(kbClient, db).
		WithContext(context.Background()).
		GetPostgresClient()
	if err != nil {
		return nil, err
	}

	return kubeDBClient, nil
}
