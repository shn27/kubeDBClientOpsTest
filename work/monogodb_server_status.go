package work

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func OpLatencies() {
	kubeDBClient, err := GetDBClient()
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
	kubeDBClient, err := GetDBClient()
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
	kubeDBClient, err := GetDBClient()
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
	kubeDBClient, err := GetDBClient()
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
	kubeDBClient, err := GetDBClient()
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
	kubeDBClient, err := GetDBClient()
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
	kubeDBClient, err := GetDBClient()
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
