package work_ElasticSearch

import (
	"fmt"
	utils "github.com/shn27/Test/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/elasticsearch"
	"testing"
)

func GetElasticSearchClient() (*elasticsearch.Client, error) {
	kbClient, err := utils.GetKBClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client: %w", err)
	}
	ref := kmapi.ObjectReference{
		Name:      "elasticsearch",
		Namespace: "demo",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "ElasticSearch",
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

	kubeDBClient, err := elasticsearch.NewKubeDBClientBuilder(kbClient, db).
		//WithPod("").
		GetElasticClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	return kubeDBClient, nil
}

func Test_getELClient(t *testing.T) {
	_, err := GetElasticSearchClient()
	if err != nil {
		return
	}
}
