package work

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mongodb"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func getK8sObject(ctrlClient client.Client) (*api.MongoDB, error) {
	obj := &api.MongoDB{}

	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "kubedb.com",
		Kind:    "MongoDB",
		Version: "v1",
	})

	if err := ctrlClient.Get(context.TODO(), client.ObjectKey{
		Name:      "mongodb",
		Namespace: "monitoring",
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

func GetDBClient() {
	for {
		time.Sleep(30 * time.Second)
		kbClient, err := getKBClient()
		if err != nil {
			fmt.Println("failed to get k8s client", err)
			continue
		}

		//db, err := getK8sObject(kbClient)
		//
		//if err != nil {
		//	fmt.Println("failed to get mongo object: %w", err)
		//}
		//
		//if db.Spec.AuthSecret == nil {
		//	fmt.Printf("auth secret is missing for %s/%s\n", db.Namespace, db.Name)
		//}

		db := &api.MongoDB{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "mongodb",
				Namespace: "monitoring",
			},
			//Spec: api.MongoDBSpec{
			//	AuthSecret: &api.SecretReference{
			//		LocalObjectReference: core.LocalObjectReference{
			//			Name: "mongodb-auth",
			//		},
			//	},
			//},
		}

		kubeDBClient, err := mongodb.NewKubeDBClientBuilder(kbClient, db).
			GetMongoClient()
		if err != nil {
			fmt.Println("failed to get kube db client: %w", err)
			continue
		}

		fmt.Println("database client created")

		// Run db.currentOp() using RunCommand
		var result bson.M
		err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
			{Key: "currentOp", Value: 1},
		}).Decode(&result)

		if err != nil {
			log.Fatal(err)
		}

		// Check for operations in progress
		inProg, ok := result["inprog"].(bson.A)
		if !ok {
			fmt.Println("No operations in progress.")
			continue
		}

		// Iterate through each operation
		for _, op := range inProg {
			// Each operation is a map of string keys to interface{} values
			opDetails, ok := op.(bson.M)
			if !ok {
				continue
			}

			// Check if the operation has been running for more than 20 seconds
			if secsRunning, _ := opDetails["secs_running"].(int32); secsRunning >= 20 {
				logSlowQuery(opDetails)
			} else {
				fmt.Println("Slow Query ok ", secsRunning)
			}

		}
	}
}

func logSlowQuery(opDetails bson.M) {
	fmt.Printf("Slow Query Detected:\n")
	fmt.Printf("  Operation ID: %v\n", opDetails["opid"])
	fmt.Printf("  Connection ID: %v\n", opDetails["connectionId"])
	fmt.Printf("  Namespace: %v\n", opDetails["ns"])
	fmt.Printf("  Command: %v\n", opDetails["command"])
	fmt.Printf("  Running Time: %v seconds\n", opDetails["secs_running"])
	fmt.Printf("  Current Operation Time: %v\n", opDetails["currentOpTime"])
	fmt.Printf("  Client: %v\n", opDetails["client"])
	fmt.Println()
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
