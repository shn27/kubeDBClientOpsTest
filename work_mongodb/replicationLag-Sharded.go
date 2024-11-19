package work

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"kubedb.dev/db-client-go/mongodb"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func shard(mongodbClient *mongodb.Client) {
	var data bytes.Buffer
	data.WriteString("MongoDB Sharded Cluster Replication Lag:\n")

	// Access the shards collection in the config database
	shardsColl := mongodbClient.Database("config").Collection("shards")

	// Retrieve all shard information
	cursor, err := shardsColl.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to retrieve shard information: %v", err)
	}
	defer cursor.Close(context.TODO())

	shards := []bson.M{}
	if err := cursor.All(context.TODO(), &shards); err != nil {
		log.Fatalf("Failed to decode shard information: %v", err)
	}

	// Iterate over each shard
	for _, shard := range shards {
		shardID := shard["_id"].(string)
		shardURI := shard["host"].(string)

		data.WriteString(fmt.Sprintf("\nShard ID: %s\n", shardID))
		data.WriteString(fmt.Sprintf("Shard URI: %s\n", shardURI))

		// Connect to individual shard
		shardClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", shardURI)))
		if err != nil {
			data.WriteString(fmt.Sprintf("Failed to connect to shard %s: %v\n", shardID, err))
			continue
		}
		defer shardClient.Disconnect(context.TODO())

		// Retrieve replication status from shard
		replicaStatus := bson.M{}
		err = shardClient.Database("admin").RunCommand(context.TODO(), bson.D{{"replSetGetStatus", 1}}).Decode(&replicaStatus)
		if err != nil {
			data.WriteString(fmt.Sprintf("Failed to get replication status for shard %s: %v\n", shardID, err))
			continue
		}

		// Check replication lag for each secondary member
		if members, ok := replicaStatus["members"].(bson.A); ok {
			for _, member := range members {
				memberInfo := member.(bson.M)
				stateStr := memberInfo["stateStr"].(string)

				if stateStr == "SECONDARY" {
					memberName := memberInfo["name"].(string)
					lastOptime := memberInfo["optime"].(bson.M)["ts"].(primitive.Timestamp).T
					lag := calculateReplicationLag(lastOptime)

					data.WriteString(fmt.Sprintf("  Secondary Node: %s\n", memberName))
					data.WriteString(fmt.Sprintf("  Replication Lag: %.2f seconds\n", lag))
				}
			}
		} else {
			data.WriteString(fmt.Sprintf("No members found in replication status for shard %s\n", shardID))
		}
	}

	// Output the data
	fmt.Println(data.String())
}

// calculateReplicationLag calculates the lag in seconds using the optime timestamp.
func calculateReplicationLag(optime uint32) float64 {
	currentTime := uint32(time.Now().Unix())
	return float64(currentTime - optime)
}

// Extract the primary host from the shard's connection string
func extractPrimaryHost(connectionString string) string {
	parts := strings.Split(connectionString, "/")
	if len(parts) < 2 {
		return ""
	}
	hosts := strings.Split(parts[1], ",")
	return hosts[0]
}

// Check replication lag for the shard's replica set
func checkReplicationLag1(client *mongo.Client, data *bytes.Buffer) error {
	// Get replication status
	replicaStatus := bson.M{}
	err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"replSetGetStatus", 1}}).Decode(&replicaStatus)
	if err != nil {
		return fmt.Errorf("failed to get replication status: %v", err)
	}

	// Check replication lag for each secondary node
	if members, ok := replicaStatus["members"].(bson.A); ok {
		for _, member := range members {
			memberInfo := member.(bson.M)
			if state, _ := memberInfo["stateStr"].(string); state == "SECONDARY" {
				secondaryName := memberInfo["name"].(string)
				lastOptime, ok := memberInfo["optimeDate"].(time.Time)
				if !ok {
					return fmt.Errorf("failed to parse optimeDate for member %s", secondaryName)
				}
				currentTime := time.Now()
				lag := currentTime.Sub(lastOptime).Seconds()

				data.WriteString(fmt.Sprintf("\nSecondary Node: %s\n", secondaryName))
				data.WriteString(fmt.Sprintf("Replication Lag: %.2f seconds\n", lag))
			}
		}
	} else {
		data.WriteString("No members found in replica set status.\n")
	}

	return nil
}
