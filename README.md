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
# Install and Wait for ingress controller to be ready
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

### 3. Deploy with Helm

Deploy the application using Helm:

```bash
# Add Bitnami repository for PostgreSQL dependency
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install the application and Wait for deployment to be ready
cd infra/helm
helm install k8s-demo . --namespace default --values ./path-to-yours-values-file
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

Customize the deployment by creating a `values.yaml`:

```yaml
# Service Configuration
service:
  port: 80                    # Service port
  protocol: TCP               # Service protocol
  type: ClusterIP             # Service type (ClusterIP/LoadBalancer)

# Deployment Configuration
deploy:
  replicaCount: 5                           # Number of application replicas
  image: iamnilotpal/k8s-demo:v2            # Container image
  resources:
    limits:
      memory: '128Mi'                       # Memory limit per pod
      cpu: '500m'                           # CPU limit per pod
    requests:
      memory: '128Mi'                       # Memory request per pod
      cpu: '500m'                           # CPU request per pod
  port:
    name: http                              # Port name
    protocol: TCP                           # Port protocol
    containerPort: 8080                     # Container port

# Application Configuration
config:
  server:
    readTimeout: 10s                        # HTTP read timeout
    idleTimeout: 120s                       # HTTP idle timeout
    writeTimeout: 10s                       # HTTP write timeout
    shutdownTimeout: 20s                    # Graceful shutdown timeout
    apiHost: "0.0.0.0:8080"                 # Server bind address
  db:
    tls: disable                            # Database TLS mode
    name: k8s-demo                          # Database name
    maxIdleConn: 5                          # Maximum idle connections
    maxOpenConn: 20                         # Maximum open connections
    scheme: postgres                        # Database scheme
    host: postgresql                        # Database host

# Database Secrets (Base64 encoded)
secrets:
  db:
    user: cG9zdGdyZXNxbA==                 # Database username (postgresql)
    password: cG9zdGdyZXNxbA==             # Database password (postgresql)

# PostgreSQL Configuration
postgresql:
  auth:
    database: k8s-demo                      # PostgreSQL database name
    username: postgresql                    # PostgreSQL username
    password: postgresql                    # PostgreSQL password
  primary:
    persistence:
      enabled: true                         # Enable persistent storage
      size: 1Gi                             # Storage size
  architecture: standalone                  # PostgreSQL architecture

# MetalLB Configuration
metallb:
  namespace: metallb-system                 # MetalLB namespace
  ipAddressPool: 172.18.255.1-172.18.255.25 # IP address range for LoadBalancer
  service:
    port: 80                               # MetalLB service port
    protocol: TCP                          # MetalLB service protocol
    type: LoadBalancer                     # MetalLB service type
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