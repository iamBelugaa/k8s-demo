# K8s Demo - Go Application with Kubernetes

A Go web application demonstrating cloud-native development patterns with Kubernetes deployment, featuring health checks, database connectivity.


## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) (20.10+)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (1.24+)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) (0.17+)
- [Helm](https://helm.sh/docs/intro/install/) (3.10+)
- [Go](https://golang.org/doc/install) (1.19+) - for local development

## Quick Start

### 1. Create Kind Cluster

First, create a local Kubernetes cluster using Kind:

```bash
# Create the cluster with the provided configuration
kind create cluster --config k8s/kind-cluster.yaml

# Wait for cluster to be ready
kubectl cluster-info --context kind-k8s-demo-cluster
```

### 2. Install NGINX Ingress Controller

Install the NGINX ingress controller for Kind:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

# Wait for ingress controller to be ready
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

### 3. Deploy with Helm

Deploy the application using Helm:

```bash
# Add Bitnami repository for PostgreSQL dependency
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install the application
cd infra/helm
helm install k8s-demo . --namespace default

# Wait for deployment to be ready
kubectl wait --for=condition=available --timeout=300s deployment/k8s-demo-deploy
```

### 4. Verify Installation

Check that all components are running:

```bash
# Check pods status
kubectl get pods

# Check services
kubectl get svc

# Check ingress
kubectl get ingress

# Test the health endpoint
curl http://localhost:80/health
```

### Helm Values

Customize the deployment by modifying `values.yaml`:

```yaml
# Replica count
deploy:
  replicaCount: 5

# Resource limits
deploy:
  resources:
    limits:
      memory: '128Mi'
      cpu: '500m'

# Service configuration
service:
  port: 80
  type: ClusterIP
```

## Development

### Local Development

Run the application locally:

```bash
# Install dependencies
go mod download

# Create .env file for development
cat > .env << EOF
# ==========================================
# DATABASE CONFIGURATION
# ==========================================
DB_TLS=require
DB_NAME=k8s-demo
DB_MAX_IDLE_CONN=5
DB_MAX_OPEN_CONN=20
DB_SCHEME=postgresql
DB_USER=postgresql
DB_HOST=localhost
DB_PASSWORD=postgresql

# ==========================================
# SERVER CONFIGURATION
# ==========================================
SERVER_READ_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=120s
SERVER_WRITE_TIMEOUT=10s
SERVER_SHUTDOWN_TIMEOUT=20s
SERVER_API_HOST=localhost:8080
EOF

# Set development environment
export ENV=DEVELOPMENT

# Run the application
go run cmd/server/main.go
```

## Monitoring and Health Checks

### Health Endpoint

The application provides a health check endpoint:

```bash
curl http://localhost:8080/health
```

Response format:
```json
{
  "success": true,
  "message": "OK"
}
```