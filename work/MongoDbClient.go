package work

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	kmapi "kmodules.xyz/client-go/api/v1"

	"go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mongodb"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetK8sObject(
	gvk schema.GroupVersionKind,
	ref kmapi.ObjectReference,
) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Kind:    gvk.Kind,
		Version: gvk.Version,
	})

	if err := r.ctrlClient.Get(context.TODO(), client.ObjectKey{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

/*

fmt.Println("==========================")
fmt.Println(gvk.Group)
fmt.Println(gvk.Kind)
fmt.Println(gvk.Version)
fmt.Println(ref.Name)
fmt.Println(ref.Namespace)

kubedb.com
MongoDB
v1
mongodb
monitoring
*/

func GetDBClient() (*mongodb.Client, error) {
	kbClient, err := getKBClient()
	if err != nil {
		fmt.Println("failed to get k8s client", err)
		return nil, err
	}
	db := &api.MongoDB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb",
			Namespace: "monitoring",
		},
	}
	kubeDBClient, err := mongodb.NewKubeDBClientBuilder(kbClient, db).
		WithPod("mongodb-0").
		WithCred("root:lK!U7bOqp1SdlUOQ").
		GetMongoClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	return kubeDBClient, nil
}

func GetDBClientLocalHost() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:Xwv88O(VRDt_RooZ@127.0.0.1:27017/?directConnection=true")) // not base64
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Minute)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("pinged")

	res, _ := client.ListDatabases(context.Background(), "")
	fmt.Println(res)

}
