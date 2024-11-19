package work

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	"log"
	"time"

	utils "github.com/shn27/Test/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mongodb"
)

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

func GetMongoDBClient() (*mongodb.Client, error) {
	kbClient, err := utils.GetKBClient()
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
func GetMongoDBClientUsingAppRef() (*mongodb.Client, error) {
	kbClient, err := utils.GetKBClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}
	ref := kmapi.ObjectReference{
		Name:      "mongodb",
		Namespace: "monitoring",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "MongoDB",
	}

	obj, err := utils.GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &api.MongoDB{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}

	kubeDBClient, err := mongodb.NewKubeDBClientBuilder(kbClient, db).
		GetMongoClient()
	if err != nil {
		return nil, fmt.Errorf("failed to build kubedb postgres client : %v", err)
	}
	return kubeDBClient, nil
}

func IsMongoDBSharded() (bool, error) {
	kbClient, err := utils.GetKBClient()
	if err != nil {
		return false, fmt.Errorf("failed to get k8s client: %w", err)
	}
	ref := kmapi.ObjectReference{
		Name:      "mongodb",
		Namespace: "monitoring",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "MongoDB",
	}

	obj, err := utils.GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return false, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &api.MongoDB{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return false, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}
	if db.Spec.ShardTopology != nil {
		fmt.Println("==============kaka=============sharded")
		return true, nil
	}
	return false, nil
}

func GetDBClientLocalHost() {
	fmt.Println("==================================client local==============")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:OswyvKqPX9iqpZY1@127.0.0.1:27017/?directConnection=true&replicaSet=shard0")) // not base64
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

	replicaStatus := bson.M{}
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"replSetGetStatus", 1}}).Decode(&replicaStatus)
	if err != nil {
		log.Fatalf("Failed to get replication status: %v", err)
	}
}
