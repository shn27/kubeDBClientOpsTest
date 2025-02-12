package work_ElasticSearch

import (
	"bytes"
	"context"
	"fmt"
	_ "github.com/elastic/go-elasticsearch/v8/esapi"
	utils "github.com/shn27/Test/utils"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/elasticsearch"
)

var ElasticSearch = &cobra.Command{
	Use:   "elasticSearch",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("elasticSearch called")
		_, err := GetElasticSearchClient()
		if err != nil {
			return
		}
	},
}
var indentL1 = "  "

func GetElasticSearchClient() (*elasticsearch.Client, error) {
	fmt.Println("GetElasticSearchClient")
	kbClient, err := utils.GetKBClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}
	ref := kmapi.ObjectReference{
		Name:      "es",
		Namespace: "demo",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "Elasticsearch",
	}
	obj, err := utils.GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &api.Elasticsearch{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}

	elasticsearchClient, err := elasticsearch.NewKubeDBClientBuilder(kbClient, db).
		WithContext(context.Background()).
		WithURL("http://127.0.0.1:9200").
		//WithPod("").
		GetElasticClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	return elasticsearchClient, nil
}

func Test() error {
	elasticSearchClient, err := GetElasticSearchClient()
	if err != nil {
		return err
	}
	nodeStats, err := elasticSearchClient.NodesStats()
	if err != nil {
		return err
	}
	fmt.Printf("node stats: %+v\n", nodeStats)

	return nil
}

func getGCData(data *bytes.Buffer, nodes map[string]interface{}) {
	data.WriteString("Analyzing Garbage Collection (GC)...\n")
	for _, nodeData := range nodes {
		node, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}
		nodeName := node["name"]
		jvm, ok := node["jvm"].(map[string]interface{})
		if !ok {
			continue
		} else {
			gc, ok := jvm["gc"].(map[string]interface{})
			if !ok {
				continue
			}
			collectors, ok := gc["collectors"].(map[string]interface{})
			if !ok {
				continue
			}
			data.WriteString(fmt.Sprintf("Node Name: %s\n", nodeName))
			for key, collectorData := range collectors {
				collector, ok := collectorData.(map[string]interface{})
				if !ok {
					continue
				}
				collectionCount := collector["collection_count"]
				data.WriteString(indentL1)
				data.WriteString(fmt.Sprintf("GC Collector: %s, Collection Count: %s, Collection Time: %s\n", key, collectionCount, collector["collection_time"]))
			}
		}
	}
	data.WriteString("\n")
}

func getThreadPoolData(data *bytes.Buffer, nodes map[string]interface{}) {
	data.WriteString("Reviewing Thread Pool Activity...\n")
	for _, nodeData := range nodes {
		node, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}
		nodeName := node["name"]
		threadPool, ok := node["thread_pool"].(map[string]interface{})
		if !ok {
			continue
		} else {
			data.WriteString(fmt.Sprintf("Node Name: %s\n", nodeName))
			for key, threadPoolData := range threadPool {
				thread, ok := threadPoolData.(map[string]interface{})
				if !ok {
					continue
				} else {
					if key == "write" || key == "search" || thread["rejected"].(float64) > 0 {
						data.WriteString(indentL1)
						data.WriteString(fmt.Sprintf("ThreadPool Threads: %s, Active: %v, rejected: %v, Completed: %v\n", key, thread["active"], thread["rejected"], thread["completed"]))
					}
				}
			}
		}
	}
	data.WriteString("\n")
}

func getHeapUsageData(data *bytes.Buffer, nodes map[string]interface{}) {
	data.WriteString("Analyzing Heap Usage...\n")
	for _, nodeData := range nodes {
		node, ok := nodeData.(map[string]interface{})
		if !ok {
			continue
		}
		nodeName := node["name"]
		jvm, ok := node["jvm"].(map[string]interface{})
		if !ok {
			continue
		} else {
			mem, ok := jvm["mem"].(map[string]interface{})
			if !ok {
				continue
			}
			heapUsedPercent := mem["heap_used_percent"]
			data.WriteString(fmt.Sprintf("Node Name: %s Heap Usage: %s\n", nodeName, heapUsedPercent))
		}
	}
	data.WriteString("\n")
}
