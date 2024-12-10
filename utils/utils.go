package work

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	kmapi "kmodules.xyz/client-go/api/v1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func GetInClusterConfig() (*restclient.Config, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		config, err = restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func GetPod(namespace string, name string) (*corev1.Pod, error) {
	config, err := GetInClusterConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v\n", err)
		return nil, err
	}
	return pod, nil
}

func GetKBClient() (client.Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		config, err = restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	cl, err := NewRuntimeClient(config)
	if err != nil {
		return nil, err
	}

	//pods, err := r.kubeClient.CoreV1().Pods(namespace).List(context.TODO(), listOptions)

	//pods, err := cl

	return cl, nil
}

func NewRuntimeClient(config *restclient.Config) (client.Client, error) {
	scheme := runtime.NewScheme()

	utilruntime.Must(api.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))

	hc, err := restclient.HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(config, hc)
	if err != nil {
		return nil, err
	}

	return client.New(config, client.Options{
		Scheme: scheme,
		Mapper: mapper,
	})
}

func GetK8sClient() (kubernetes.Interface, error) {
	// First, try to get the client from within the cluster
	config, err := restclient.InClusterConfig()
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

func GetResource(resourceName string) {
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

func GetK8sObject(
	gvk schema.GroupVersionKind,
	ref kmapi.ObjectReference,
	kbClient client.Client,
) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Kind:    gvk.Kind,
		Version: gvk.Version,
	})

	if err := kbClient.Get(context.TODO(), client.ObjectKey{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}, obj); err != nil {
		return nil, err
	}
	return obj, nil
}
