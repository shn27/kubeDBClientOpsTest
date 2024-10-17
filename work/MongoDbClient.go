package work

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mgoptions "go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/mongodb"
	"log"
	"strings"
	"time"
)

func GetDBClient() {

	for {
		db := &api.MongoDB{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "mongodb",
			},
		}

		kbClient, err := getKBClient()
		if err != nil {
			fmt.Println("========", err)
			continue
		}
		client, err := mongodb.NewKubeDBClientBuilder(kbClient, db).
			WithPod("mongodb-0").
			WithContext(context.Background()).
			WithCred("root:T7a5*UYKY5W1MgtL").
			GetMongoClient()
		if err != nil {
			fmt.Println("============", err)
			continue
		}
		fmt.Println("Nice and attractive database client created", client)
		wait.RealTimer(time.Minute)
	}

	//	fmt.Println(client.ListDatabases(context.Background(), nil, nil))
	//res := make(map[string]interface{})
	//err = client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "isMaster", Value: "1"}}).Decode(&res)
	//if err != nil {
	//	fmt.Printf("error checking isPrimary. error: %v\n", err)
	//}
}
func GetDBClientLocalHost() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:Xwv88O(VRDt_RooZ@127.0.0.1:27017/?directConnection=true")) // not base64
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Minute)
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

func getMongoDBClientOpts(db *api.MongoDB) {
	url := ""
	repSetConfig := ""
	authDatabase := ""
	if db.Spec.ShardTopology != nil {
		// Shard
		url = strings.Join(db.MongosHosts(), ",")
	} else {
		// Standalone or ReplicaSet
		url = strings.Join(db.Hosts(), ",")
	}
	cred := "root:T7a5*UYKY5W1MgtL"

	var clientOpts *mgoptions.ClientOptions
	uri := fmt.Sprintf("mongodb://%s@%s/admin?%vtls=true&tlsCAFile=%v&tlsCertificateKeyFile=%v%v", cred, url, repSetConfig, authDatabase)
	clientOpts = mgoptions.Client().ApplyURI(uri)
	fmt.Println("===========================kaka ==", clientOpts.AppName)
}
