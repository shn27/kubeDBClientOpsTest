package work

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"kubedb.dev/db-client-go/mongodb"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func replicationLag(client *mongodb.Client) {
	var data bytes.Buffer
	data.WriteString("MongoDB Replication Lag and Network Stats:\n")

	// Get replication status
	replicaStatus := bson.M{}
	err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"replSetGetStatus", 1}}).Decode(&replicaStatus)
	if err != nil {
		log.Fatalf("Failed to get replication status: %v", err)
	}

	// Check replication lag for each secondary node
	// Process each member in the replica set
	if members, ok := replicaStatus["members"].(bson.A); ok {
		for _, member := range members {
			memberInfo := member.(bson.M)
			stateStr, _ := memberInfo["stateStr"].(string)
			memberName, _ := memberInfo["name"].(string)

			// Display basic info about each node
			data.WriteString(fmt.Sprintf("\nNode: %s (State: %s)\n", memberName, stateStr))

			// Check replication lag for secondary nodes
			if stateStr == "SECONDARY" {
				checkReplicationLag(memberInfo, &data)
			}
		}
	} else {
		data.WriteString("No members found in replica set status.\n")
	}

	fmt.Println(data.String())
}

func checkReplicationLag(memberInfo bson.M, data *bytes.Buffer) {
	// Get the last applied operation time (optimeDate) from the secondary node
	if optime, ok := memberInfo["optimeDate"].(primitive.DateTime); ok {
		lastOptime := optime.Time()
		currentTime := time.Now()
		lag := currentTime.Sub(lastOptime).Seconds()

		data.WriteString(fmt.Sprintf("  Last Applied Optime: %v\n", lastOptime))
		data.WriteString(fmt.Sprintf("  Replication Lag: %.2f seconds\n", lag))

		// Provide warning if lag is too high
		if lag > 30 { // Threshold in seconds
			data.WriteString("  Warning: High replication lag detected. Check network or disk performance.\n")
		}
	} else {
		data.WriteString("  No valid optimeDate found; potential replication issue.\n")
	}
}

func checkNetworkStats(client *mongodb.Client, data *bytes.Buffer) error {
	// Get serverStatus
	serverStatus := bson.M{}
	err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return err
	}

	// Extract relevant network data
	if network, ok := serverStatus["network"].(bson.M); ok {
		if bytesIn, ok := network["bytesIn"].(int64); ok {
			data.WriteString(fmt.Sprintf("Bytes In: %d\n", bytesIn))
		}
		if bytesOut, ok := network["bytesOut"].(int64); ok {
			data.WriteString(fmt.Sprintf("Bytes Out: %d\n", bytesOut))
		}
		if numRequests, ok := network["numRequests"].(int64); ok {
			data.WriteString(fmt.Sprintf("Number of Requests: %d\n", numRequests))
		}
	} else {
		data.WriteString("Network stats unavailable\n")
	}

	// Connection stats
	if connections, ok := serverStatus["connections"].(bson.M); ok {
		if current, ok := connections["current"].(int32); ok {
			data.WriteString(fmt.Sprintf("Current Connections: %d\n", current))
		}
		if available, ok := connections["available"].(int32); ok {
			data.WriteString(fmt.Sprintf("Available Connections: %d\n", available))
		}
	} else {
		data.WriteString("Connection stats unavailable\n")
	}
	fmt.Println(data.String())
	return nil
}
func getVmStat() {
	data, err := ioutil.ReadFile("/proc/diskstats")
	if err != nil {
		fmt.Printf("Failed to read /proc/diskstats: %v\n", err)
		return
	}

	var blocksRead, blocksWritten int64
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue
		}
		read, _ := strconv.ParseInt(fields[5], 10, 64)    // Number of sectors read
		written, _ := strconv.ParseInt(fields[9], 10, 64) // Number of sectors written
		blocksRead += read
		blocksWritten += written
	}

	fmt.Printf("Blocks in: %d, Blocks out: %d\n", blocksRead, blocksWritten)
}
