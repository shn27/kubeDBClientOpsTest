# set up
`kc create -f hack/samples`

# Run test
set vm's k3s kubeconfig in .kube/config

`kc port-forward -n demo svc/es 9200`
```
elasticsearchClient, err := elasticsearch.NewKubeDBClientBuilder(kbClient, db).
WithContext(context.Background()).
WithURL("http://127.0.0.1:9200").
GetElasticClient()
```

# Run in cluster
1. make build
2. make push
3. make deploy

# Uninstall
1. make clean