package work

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/postgres"
)

func GetPostgresClient() (*postgres.Client, error) {
	kbClient, err := getKBClient()
	if err != nil {
		fmt.Println("failed to get k8s client", err)
		return nil, err
	}
	db := &api.Postgres{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb",
			Namespace: "monitoring",
		},
	}

	kubeDBClient, err := postgres.NewKubeDBClientBuilder(kbClient, db).
		WithContext(context.Background()).
		WithPod("postgres-0").
		GetPostgresClient()
	if err != nil {
		return nil, err
	}

	return kubeDBClient, nil
}
