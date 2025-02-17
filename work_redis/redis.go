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
)

func getRedisClient() (*redis.ClusterClient, error) {
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
		//WithURL("127.0.0.1:6379").
		//WithPod("redis-0").
		GetRedisClusterClient(context.Background())

	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	return redisClient, nil
}
