package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getKBClient() (client.Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	cl, err := client.New(config, client.Options{
		Scheme: runtime.NewScheme(),
	})
	if err != nil {
		return nil, err
	}
	return cl, nil
}
func GetK8sClient() (kubernetes.Interface, error) {
	// First, try to get the client from within the cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		// If not running inside the cluster, fallback to kubeconfig file
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Use the config to create a new Kubernetes client
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k8sClient, nil
}

func getResource(resourceName string) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}
	gvr := schema.GroupVersionResource{
		Group:    "monitoring.appscode.com",
		Version:  "v1alpha1",
		Resource: "runbooks",
	}
	resource, err := dynamicClient.Resource(gvr).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error getting resource: %v", err)
	}
	fmt.Println(resource.Object)
}
