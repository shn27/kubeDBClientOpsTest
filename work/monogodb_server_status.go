package work

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"k8s.io/klog/v2"
	"kubedb.dev/db-client-go/mongodb"
)

func OpLatencies() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}

	// Run db.currentOp() using RunCommand
	//for {
	time.Sleep(10 * time.Second)
	var result bson.M
	err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
		{Key: "serverStatus", Value: 1},
	}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}
	if opLatencies, ok := result["opLatencies"].(bson.M); ok {
		fmt.Println("Operation Latencies:")
		for op, data := range opLatencies {
			fmt.Printf("  %s: %v\n", op, data)
		}
	} else {
		fmt.Println("No opLatencies data found in serverStatus output.")
	}
	//}
}

func Network() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}
	// Run db.currentOp() using RunCommand
	//	for {
	time.Sleep(10 * time.Second)
	var result bson.M
	err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
		{Key: "serverStatus", Value: 1},
	}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}
	// Extract network information
	if networkStats, ok := result["network"].(bson.M); ok {
		fmt.Println("Network Statistics:")
		for key, value := range networkStats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	} else {
		fmt.Println("No network data found in serverStatus output.")
	}
	//	}
}

func Connections() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}
	// Run db.currentOp() using RunCommand
	//for {
	time.Sleep(10 * time.Second)
	var result bson.M
	err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
		{Key: "serverStatus", Value: 1},
	}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}
	// Extract connections information
	if connectionsStats, ok := result["connections"].(bson.M); ok {
		fmt.Println("Connections Statistics:")
		for key, value := range connectionsStats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	} else {
		fmt.Println("No connections data found in serverStatus output.")
	}
	//}
}

func wiredTiger() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}
	// Run db.currentOp() using RunCommand
	//for {
	time.Sleep(10 * time.Second)
	var result bson.M
	err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
		{Key: "serverStatus", Value: 1},
	}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}
	// Extract wiredTiger statistics
	if wiredTigerStats, ok := result["wiredTiger"].(bson.M); ok {
		fmt.Println("WiredTiger Statistics:")
		for key, value := range wiredTigerStats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	} else {
		fmt.Println("No wiredTiger data found in serverStatus output.")
	}
	//	}
}

func metrics_network() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}
	// Run db.currentOp() using RunCommand
	//for {
	time.Sleep(10 * time.Second)
	var result bson.M
	err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
		{Key: "serverStatus", Value: 1},
	}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}
	// Extract metrics.network information
	if metrics, ok := result["metrics"].(bson.M); ok {
		if networkStats, ok := metrics["network"].(bson.M); ok {
			fmt.Println("Network Metrics:")
			for key, value := range networkStats {
				fmt.Printf("  %s: %v\n", key, value)
			}
		} else {
			fmt.Println("No network metrics data found in serverStatus output.")
		}
	} else {
		fmt.Println("No metrics data found in serverStatus output.")
	}
	//}
}

func Metrics_cursor() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}
	// Run db.currentOp() using RunCommand
	//	for {
	time.Sleep(10 * time.Second)
	var result bson.M
	err = kubeDBClient.Database("admin").RunCommand(context.TODO(), bson.D{
		{Key: "serverStatus", Value: 1},
	}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}
	// Extract metrics.cursor information
	if metrics, ok := result["metrics"].(bson.M); ok {
		if cursorStats, ok := metrics["cursor"].(bson.M); ok {
			fmt.Println("Cursor Metrics:")
			for key, value := range cursorStats {
				fmt.Printf("  %s: %v\n", key, value)
			}
		} else {
			fmt.Println("No cursor metrics data found in serverStatus output.")
		}
	} else {
		fmt.Println("No metrics data found in serverStatus output.")
	}
	//	}
}

func DbCurrentOp() {
	kubeDBClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %s\n", err.Error())
		return
	}

	// Run db.currentOp() using RunCommand
	//	for {
	time.Sleep(10 * time.Second)
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
		return
	}
	// Iterate through each operation
	for _, op := range inProg {
		// Each operation is a map of string keys to interface{} values
		opDetails, ok := op.(bson.M)
		if !ok {
			return
		}
		// Check if the operation has been running for more than 20 seconds
		if secsRunning, _ := opDetails["secs_running"].(int32); secsRunning >= 0 {
			logSlowQuery(opDetails)
		} else {
			fmt.Println("Slow Query ok ", secsRunning)
		}
	}
	//	}
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

func checkOpenCursors(client *mongodb.Client, data *bytes.Buffer) error {
	result, err := getServerStatus(context.TODO(), client)
	if err != nil {
		fmt.Println("============lalalla")
		log.Fatal(err)
	}

	// Access the cursor metrics from serverStatus
	if metrics, ok := result["metrics"].(bson.M); ok {
		if cursorMetrics, ok := metrics["cursor"].(bson.M); ok {
			if openCursors, ok := cursorMetrics["open"].(primitive.M); ok {
				if total, ok := openCursors["total"].(int64); ok {
					data.WriteString(fmt.Sprintf("Number of Open Cursors: %d\n", total))
					if total >= 0 { // Example threshold
						data.WriteString("Warning: High number of open cursors detected. Ensure that cursors are closed after use.\n")
						fmt.Println("Warning: High number of open cursors detected. Ensure that cursors are closed after use.\n")
					}
				}
			}

			fmt.Println("==================64     ", cursorMetrics["timedOut"].(int64))
			// Retrieve cursor timeouts
			if timedOutCursors, ok := cursorMetrics["timedOut"].(int64); ok {
				data.WriteString(fmt.Sprintf("Number of Timed Out Cursors: %d\n", timedOutCursors))
				fmt.Printf("Number of Timed Out Cursors: %d\n\n", timedOutCursors)
				if timedOutCursors >= 0 {
					data.WriteString("Note: Some cursors have timed out. Ensure that your application properly closes cursors to prevent this.\n")
					fmt.Println("Note: Some cursors have timed out. Ensure that your application properly closes cursors to prevent this.\n")
				}
			}
		} else {
			data.WriteString("Cursor metrics not found in serverStatus output.\n")
			fmt.Println("Cursor metrics not found in serverStatus output.\n")
		}
	} else {
		data.WriteString("Metrics section not found in serverStatus output.\n")
		fmt.Println("Metrics section not found in serverStatus output.\n")
	}
	return nil
}
func checkConnectionMetrics(client *mongodb.Client, data *bytes.Buffer) error {
	result, err := getServerStatus(context.TODO(), client)
	if err != nil {
		fmt.Println("============lalalla")
		log.Fatal(err)
	}

	if connections, ok := result["connections"].(bson.M); ok {
		fmt.Println("===========", connections["current"].(int32), connections["available"].(int32), connections["totalCreated"].(int32))
		if current, ok := connections["current"].(int32); ok {
			data.WriteString(fmt.Sprintf("Current Connections: %d\n", current))
			fmt.Printf("Current Connections: %d\n", current)
		}
		if available, ok := connections["available"].(int32); ok {
			data.WriteString(fmt.Sprintf("Available Connections: %d\n", available))
			fmt.Printf("Available Connections: %d\n", available)
		}
		if totalCreated, ok := connections["totalCreated"].(int32); ok {
			data.WriteString(fmt.Sprintf("Total Connections Created: %d\n", totalCreated))
			fmt.Printf("Total Connections Created: %d\n", totalCreated)
		}
	} else {
		data.WriteString("Connection metrics not found in serverStatus output.\n")
		fmt.Println("Connection metrics not found in serverStatus output.\n")
	}
	return nil
}

// Function to group active connections by IP address
func analyzeActiveConnectionsByIP(client *mongodb.Client, data *bytes.Buffer) error {
	// Running the "currentOp" command
	var result bson.M
	err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "currentOp", Value: 1}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	// Access the "inprog" field which contains the list of active operations
	inprog, ok := result["inprog"].(bson.A)
	if !ok {
		log.Fatal("inprog field is not an array")
	}

	ipConnections := make(map[string]int)
	totalConnections := 0

	// Iterate over each operation in "inprog"
	for _, op := range inprog {
		opDoc, ok := op.(bson.M)
		if !ok {
			continue
		}

		// Extract client IP
		if clientInfo, ok := opDoc["client"].(string); ok {
			ipAddress := "Internal"
			if parts := strings.Split(clientInfo, ":"); len(parts) > 0 {
				ipAddress = parts[0]
			}
			// Increment count for this IP
			ipConnections[ipAddress]++
			// Increment total connection count
			totalConnections++
		}
	}

	// Output the summarized connection counts
	data.WriteString(fmt.Sprintf("Total Active Connections: %d\n", totalConnections))
	fmt.Printf("Total Active Connections: %d\n", totalConnections)
	for ip, count := range ipConnections {
		data.WriteString(fmt.Sprintf("  IP: %s - Active Connections: %d\n", ip, count))
		fmt.Printf("  IP: %s - Active Connections: %d\n", ip, count)
	}
	return nil
}

var (
	highTrafficThreshold    = 1000000                // Example: 1 MB for network traffic
	highRequestThreshold    = 5000                   // Example: 5,000 requests
	cacheLimit              = 5 * 1024 * 1024 * 1024 // Example: 5 GB cache
	highDirtyBytesThreshold = 100 * 1024 * 1024      // Example: 100 MB dirty bytes
	highLogWriteThreshold   = 500 * 1024 * 1024      // Example: 500 MB log writes
	sampleInterval          = 10 * time.Second       // Sample interval for rolling average
	sampleCount             = 5                      // Number of samples for rolling average
)

func analyzeNetworkPerformance(result bson.M) {
	if metrics, ok := result["metrics"].(bson.M); ok {
		if network, ok := metrics["network"].(bson.M); ok {
			bytesIn, _ := network["bytesIn"].(int)
			bytesOut, _ := network["bytesOut"].(int)
			numRequests, _ := network["numRequests"].(int)

			fmt.Println("Network Performance:")
			fmt.Printf("  bytesIn: %d\n  bytesOut: %d\n  numRequests: %d\n", bytesIn, bytesOut, numRequests)

			if bytesIn > highTrafficThreshold || bytesOut > highTrafficThreshold {
				log.Println("  Warning: High network traffic detected.")
			}

			if numRequests > highRequestThreshold {
				log.Println("  Warning: High number of requests detected.")
			}
		} else {
			log.Println("No network metrics data found.")
		}
	}
}

func analyzeDiskPerformance(result bson.M) {
	if wiredTiger, ok := result["wiredTiger"].(bson.M); ok {
		if cache, ok := wiredTiger["cache"].(bson.M); ok {
			bytesInCache, _ := cache["bytes currently in cache"].(int)
			trackedDirtyBytes, _ := cache["tracked dirty bytes in cache"].(int)
			fmt.Println("Disk Performance:")
			fmt.Printf("  bytes currently in cache: %d\n  tracked dirty bytes in cache: %d\n", bytesInCache, trackedDirtyBytes)

			if bytesInCache >= cacheLimit {
				log.Println("  Warning: Cache is full; may impact disk I/O performance.")
			}

			if trackedDirtyBytes > highDirtyBytesThreshold {
				log.Println("  Warning: High dirty bytes in cache; potential disk write latency.")
			}
		}

		if logMetrics, ok := wiredTiger["log"].(bson.M); ok {
			logBytesWritten, _ := logMetrics["total log bytes written"].(int)
			fmt.Printf("  total log bytes written: %d\n", logBytesWritten)

			if logBytesWritten > highLogWriteThreshold {
				log.Println("  Warning: High log write volume; potential disk I/O bottleneck.")
			}
		}
	} else {
		log.Println("No WiredTiger data found.")
	}
}

func CallAllMongoMethod() {
	for {
		//time.Sleep(10 * time.Second)
		metrics_network()
		//time.Sleep(10 * time.Second)
		Metrics_cursor()
		//time.Sleep(10 * time.Second)
		DbCurrentOp()
		//time.Sleep(10 * time.Second)
		Connections()
		//time.Sleep(10 * time.Second)
		OpLatencies()
		//time.Sleep(10 * time.Second)
		Network()
		//time.Sleep(10 * time.Second)
		wiredTiger()
		time.Sleep(10 * time.Second)
	}
}

func getServerStatus(ctx context.Context, client *mongodb.Client) (bson.M, error) {
	var result bson.M
	err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}}).Decode(&result)
	return result, err
}

func Ans() {
	fmt.Println("=================")
	mongodbClient, err := GetMongoDBClient()
	if err != nil {
		fmt.Printf("get db client error: %v\n", err)
		return
	}
	_ = mongodbClient
	klog.Info("=======NEW LOG=======")
}
