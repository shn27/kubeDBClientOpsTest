package work_redis

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	utils "github.com/shn27/Test/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

func getRedisClient() (*redis.Client, error) {
	kbClient, err := utils.GetKBClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}
	ref := kmapi.ObjectReference{
		Name:      "redis",
		Namespace: "demo",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "Redis",
	}
	obj, err := utils.GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &api.Redis{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}

	redisClient, err := redis.NewKubeDBClientBuilder(kbClient, db).
		WithURL("127.0.0.1:6379").
		//WithPod("redis-0").
		GetRedisClient(context.Background())

	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}

	result, err := redisClient.Info(context.Background(), "replication").Result()
	fmt.Println(result)

	return redisClient, nil
}

func checkClusterInfo(client *redis.Client) error {
	clusterInfo, err := client.ClusterInfo(context.Background()).Result()
	if err != nil {
		return err
	}
	lines := strings.Split(clusterInfo, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cluster_state:") {
			state := strings.TrimPrefix(line, "cluster_state:")
			fmt.Println("📊 Cluster State:", strings.TrimSpace(state))
			if strings.TrimSpace(state) != "ok" {
				fmt.Println("⚠️  Warning: Cluster is not in a healthy state!")
			}
		}
		if strings.HasPrefix(line, "cluster_slots_fail:") {
			failures := strings.TrimPrefix(line, "cluster_slots_fail:")
			if strings.TrimSpace(failures) != "0" {
				fmt.Println("⚠️  Warning: Some slots have failures!")
			}
		}
	}
	return nil
}

func checkNetwork(client *redis.Client) error {
	// Check connectivity for each Redis node

	redisClient, err := getRedisClient()
	if err != nil {
		return err
	}

	fmt.Println("🔍 Checking Redis connectivity...")
	nodes, err := getRedisNodes(redisClient)
	for _, node := range nodes {
		if checkRedisConnectivity(redisClient, node) {
			fmt.Printf("✅ Redis node is reachable: %s\n", node)
		} else {
			fmt.Printf("❌ Redis node is unreachable: %s\n", node)
		}
	}
	return nil
}

// Get Redis nodes dynamically using CLUSTER NODES
func getRedisNodes(client *redis.Client) ([]string, error) {
	nodes := []string{}
	result, err := client.ClusterNodes(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster nodes: %v", err)
	}

	for _, line := range strings.Split(result, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 2 {
			address := fields[1] // Format: ip:port@bus-port
			// Extract only ip:port (before '@')
			nodes = append(nodes, strings.Split(address, "@")[0])
		}
	}
	return nodes, nil
}

// Check if a Redis node is reachable
func checkRedisConnectivity(client *redis.Client, node string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis node %s is unreachable: %v\n", node, err)
		return false
	}
	return pong == "PONG"
}

// Function to check for split-brain scenario
func checkClusterNodeInfo(client *redis.Client) error {
	ctx := context.Background()
	clusterNodes, err := client.ClusterNodes(ctx).Result()
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// Count master nodes
	masterCount := 0
	for _, line := range strings.Split(clusterNodes, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 2 && strings.Contains(fields[2], "master") {
			masterCount++
		}
	}

	if masterCount > 1 {
		return fmt.Errorf("detected multiple masters (%d), possible split-brain scenario", masterCount)
	}

	fmt.Println("Cluster is healthy, only one master found.")
	return nil
}

func checkResourceLimits() {
	v, _ := mem.VirtualMemory()
	cpuUsage, _ := cpu.Percent(0, false)

	fmt.Printf("🔍 Memory Usage: %.2f%%\n", v.UsedPercent)
	fmt.Printf("🔍 CPU Usage: %.2f%%\n", cpuUsage[0])

	if v.UsedPercent > 90 {
		fmt.Println("⚠️ High memory usage detected (>90%)!")
	}
	if cpuUsage[0] > 90 {
		fmt.Println("⚠️ High CPU usage detected (>90%)!")
	}
}

// Function to check replica synchronization
func checkReplicaSync(client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Run INFO replication command
	result, err := client.Info(ctx, "replication").Result()
	if err != nil {
		log.Fatalf("❌ Failed to get replication info: %v\n", err)
	}

	// Parse the response
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Identify and check slave status
		if strings.HasPrefix(line, "slave") {
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				continue
			}

			slaveInfo := strings.Split(parts[1], ",")
			slaveAddr := slaveInfo[0]
			slaveState := strings.TrimSpace(strings.Split(slaveInfo[2], "=")[1])

			if slaveState != "online" {
				fmt.Printf("❌ Disconnected Replica: %s (State: %s)\n", slaveAddr, slaveState)
			} else {
				fmt.Printf("✅ Replica synchronized: %s (State: online)\n", slaveAddr)
			}
		}
	}
}

// Function to check master-slave offset difference
func checkReplicationLag(client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Run INFO replication
	result, err := client.Info(ctx, "replication").Result()
	if err != nil {
		log.Fatalf("❌ Failed to get replication info: %v\n", err)
	}

	var masterOffset int64
	slaveOffsets := make(map[string]int64)

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Get master offset
		if strings.HasPrefix(line, "master_repl_offset") {
			parts := strings.Split(line, ":")
			masterOffset, _ = strconv.ParseInt(parts[1], 10, 64)
		}

		// Get slave offsets
		if strings.HasPrefix(line, "slave") {
			parts := strings.Split(line, ":")
			slaveInfo := strings.Split(parts[1], ",")
			slaveAddr := slaveInfo[0]
			offsetStr := strings.Split(slaveInfo[3], "=")[1]
			slaveOffset, _ := strconv.ParseInt(offsetStr, 10, 64)

			slaveOffsets[slaveAddr] = slaveOffset
		}
	}

	// Compare slave offsets with master offset
	for slave, offset := range slaveOffsets {
		lag := masterOffset - offset
		fmt.Printf("🔍 Slave: %s | Offset Lag: %d\n", slave, lag)
		if lag > 100 { // Threshold can be adjusted
			fmt.Printf("⚠️ Warning: Slave %s has high replication lag!\n", slave)
		}
	}
}

// Fetch and analyze Redis connections
func analyzeConnections(client *redis.Client) {

	// Run CLIENT LIST command
	result, err := client.ClientList(context.Background()).Result()
	if err != nil {
		log.Fatalf("❌ Failed to fetch client connections: %v\n", err)
	}

	// Parse CLIENT LIST output
	connections := strings.Split(result, "\n")
	connectionCount := len(connections) - 1 // Last entry is empty

	fmt.Printf("🔍 Active Redis Connections: %d\n", connectionCount)

	// Track connections by IP
	connectionMap := make(map[string]int)

	for _, conn := range connections {
		fields := strings.Fields(conn)
		if len(fields) > 1 {
			for _, field := range fields {
				if strings.HasPrefix(field, "addr=") {
					ipPort := strings.TrimPrefix(field, "addr=")
					ip := strings.Split(ipPort, ":")[0] // Extract IP (ignore port)
					connectionMap[ip]++
					break
				}
			}
		}
	}

	// Print connection stats
	fmt.Println("📊 Connection Breakdown by IP:")
	for ip, count := range connectionMap {
		fmt.Printf("   - %s: %d connections\n", ip, count)
		if count > 50 { // Threshold (adjust as needed)
			fmt.Printf("⚠️  Warning: High number of connections from %s!\n", ip)
		}
	}
}
