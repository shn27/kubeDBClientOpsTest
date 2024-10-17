# Define variables
IMAGE_NAME=my-app-image
IMAGE_TAG=latest
REGISTRY_URL=shn27# Replace with your Docker registry (Docker Hub or private registry)
DOCKERFILE_PATH=./Dockerfile
K8S_DEPLOYMENT_FILE=k8s/deployment.yaml
KIND_CLUSTER_NAME=kind
NAMESPACE=default

# Build the Docker image
.PHONY: build
build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f $(DOCKERFILE_PATH) .

# Tag the Docker image for the registry
.PHONY: tag
tag:
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)

# Push the Docker image to the registry
.PHONY: push
push: tag
	docker push $(REGISTRY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)

# Deploy Docker image to Kind cluster
.PHONY: deploy
deploy: build
deploy:tag
deploy:push
	# Apply the Kubernetes deployment and service YAML
	kubectl apply -f $(K8S_DEPLOYMENT_FILE) -n $(NAMESPACE)



# Clean up the environment (optional)
.PHONY: clean
clean:
	kubectl delete -f $(K8S_DEPLOYMENT_FILE) -n $(NAMESPACE)

# Help command to display available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build   - Build the Docker image"
	@echo "  make push    - Tag and push the Docker image to the registry"
	@echo "  make deploy  - Build and deploy the Docker image to the Kind cluster"
	@echo "  make clean   - Remove the deployed resources from the Kubernetes cluster"
