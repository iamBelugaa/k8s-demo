# K8s Demo - Cloud-Native Go Application

A Go web application with observability stack including Prometheus, Grafana,
Jaeger, and PostgreSQL database.

## Local Setup and Deployment

### Step 1: Create Kubernetes Cluster

```bash
kind create cluster --config infra/k8s/kind-cluster.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

### Step 2: Prepare Helm Dependencies and Install Required CRDs

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus-operator prometheus-community/kube-prometheus-stack \
  --set prometheus.enabled=false \
  --set alertmanager.enabled=false \
  --set grafana.enabled=false \
  --set kubeStateMetrics.enabled=false \
  --set nodeExporter.enabled=false \
  --set prometheusOperator.enabled=true
```

### Step 3: Deploy the Complete Stack

```bash
cd infra/helm
helm install k8s-demo . --values values.yaml
kubectl get pods --watch
```

### Step 4: Configure Local Access

```bash
echo "127.0.0.1 k8s-demo.com" | sudo tee -a /etc/hosts
echo "127.0.0.1 jaeger-demo.com" | sudo tee -a /etc/hosts
echo "127.0.0.1 grafana-demo.com" | sudo tee -a /etc/hosts
echo "127.0.0.1 prometheus-demo.com" | sudo tee -a /etc/hosts
```

After adding these entries, you can access the services through your web browser
using the domain names instead of IP addresses and ports.
