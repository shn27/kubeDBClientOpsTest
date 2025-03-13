package work_redis

import (
	"context"
	"fmt"
	utils "github.com/shn27/Test/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/redis"
	"strings"
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

	ans, _ := redisClient.ClusterNodes(context.Background()).Result()
	fmt.Println(ans)

	return redisClient, nil
}

// Function to check for split-brain scenario
func checkSplitBrain(client *redis.ClusterClient) error {
	ctx := context.Background()
	clusterNodes, err := client.ClusterNodes(ctx).Result()
	if err != nil {
		return fmt.Errorf("error fetching cluster nodes: %v", err)
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
