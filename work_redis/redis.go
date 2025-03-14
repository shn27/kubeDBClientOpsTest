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
	fmt.Println("clusterInfo:", clusterInfo)
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
		log.Printf("❌ Redis node %s is unreachable: %v\n", node, err)
		return false
	}
	return pong == "PONG"
}

// Function to check for split-brain scenario
func checkClusterNodeInfo(client *redis.Client) error {
	ctx := context.Background()
	clusterNodes, err := client.ClusterNodes(ctx).Result()
	if err != nil {
		return fmt.Errorf("error fetching cluster nclusterNodes = {string} \"cbba139773c3602d042ffa72ce9203b66062e27d 10.42.0.90:6379@16379 master - 0 1741931362000 2 connected 5461-10922\\n67fc5c4e342db9bb48d5d293092173de6e36b9a2 10.42.0.98:6379@16379 master - 0 1741931363000 3 connected 10923-16383\\n9649dee04d62a6743648347b2bc3214776c48e22 10.42.0.100:6379@16379 slave 3f1ddf3fa187cc1cc99962dd1acfac74648950ad 0 1741931364085 1 connected\\n14217248cf9ab80312267de67abc99bc2e9920f6 10.42.0.102:6379@16379 slave cbba139773c3602d042ffa72ce9203b66062e27d 0 1741931363583 2 connected\\n3f1ddf3fa187cc1cc99962dd1acfac74648950ad 10.42.0.86:6379@16379 myself,master - 0 0 1 connected 0-5460\\n1bf7ab8412f8f9a664414ecaff23458ad8dc0b54 10.42.0.103:6379@16379 slave 67fc5c4e342db9bb48d5d293092173de6e36b9a2 0 1741931363583 3 connected\\n\"odes: %v", err)
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
			slaveState := strings.TrimSpace(strings.Split(slaveInfo[1], "=")[1])

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
			offsetStr := strings.Split(slaveInfo[2], "=")[1]
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
