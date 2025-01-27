package work_ElasticSearch

import (
	"context"
	"fmt"
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
		//WithPod("").
		GetElasticClient()
	if err != nil {
		fmt.Println("failed to get kube db client: %w", err)
		return nil, err
	}
	fmt.Println("========================hey done su")
	return elasticsearchClient, nil
}
