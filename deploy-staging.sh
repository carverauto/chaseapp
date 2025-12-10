#!/bin/bash

# Staging Deployment Script for ChaseApp Go API
# This script deploys the new Go monolithic API to staging environment

set -e

# Configuration
PROJECT_ID="chaseapp-8459b"
CLUSTER_NAME="chaseapp-staging"
REGION="us-central1"
NAMESPACE="chaseapp-staging"
REGISTRY="ghcr.io/carverauto/chaseapp-api"
IMAGE_TAG="${1:-staging-$(date +%Y%m%d-%H%M%S)}"

echo "🚀 Deploying ChaseApp Go API to Staging"
echo "Project: $PROJECT_ID"
echo "Cluster: $CLUSTER_NAME"
echo "Namespace: $NAMESPACE"
echo "Image: $REGISTRY:$IMAGE_TAG"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."

    # Check gcloud CLI
    if ! command -v gcloud &> /dev/null; then
        error "gcloud CLI is not installed"
        exit 1
    fi

    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        error "kubectl is not installed"
        exit 1
    fi

    # Check if we're authenticated
    if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
        error "Not authenticated with gcloud. Run: gcloud auth login"
        exit 1
    fi

    # Check if we have access to the project
    if ! gcloud projects describe "$PROJECT_ID" &> /dev/null; then
        error "Cannot access project $PROJECT_ID"
        exit 1
    fi

    log "Prerequisites check passed"
}

# Build and push Docker image
build_and_push() {
    log "Building Docker image..."

    cd api/

    # Build the Docker image
    docker build -t "$REGISTRY:$IMAGE_TAG" .
    docker tag "$REGISTRY:$IMAGE_TAG" "$REGISTRY:staging-latest"

    log "Pushing Docker image..."
    docker push "$REGISTRY:$IMAGE_TAG"
    docker push "$REGISTRY:staging-latest"

    cd ..
    log "Docker image built and pushed successfully"
}

# Deploy to Kubernetes
deploy_to_k8s() {
    log "Deploying to Kubernetes..."

    # Get cluster credentials
    gcloud container clusters get-credentials "$CLUSTER_NAME" --region "$REGION" --project "$PROJECT_ID"

    # Create namespace if it doesn't exist
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    # Set current context to the namespace
    kubectl config set-context --current --namespace="$NAMESPACE"

    # Update kustomize with the new image tag
    cd k8s/

    # Create staging overlay if it doesn't exist
    if [ ! -d "staging" ]; then
        log "Creating staging overlay..."
        mkdir -p staging

        cat > staging/kustomization.yaml << EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: chaseapp-staging
bases:
  - ../base
images:
  - name: ghcr.io/carverauto/chaseapp-api
    newTag: $IMAGE_TAG
patchesStrategicMerge:
  - api-deployment-staging.yaml
EOF

        cat > staging/api-deployment-staging.yaml << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: api
          resources:
            requests:
              cpu: 50m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 256Mi
          env:
            - name: ENVIRONMENT
              value: "staging"
EOF
    fi

    # Apply the configuration
    kubectl apply -k staging/

    cd ..
    log "Kubernetes deployment applied"
}

# Wait for deployment to be ready
wait_for_deployment() {
    log "Waiting for deployment to be ready..."

    kubectl rollout status deployment/api --timeout=300s

    # Check if pods are running
    if ! kubectl get pods -l app.kubernetes.io/component=api | grep -q "Running"; then
        error "Deployment failed - pods are not running"
        kubectl get pods -l app.kubernetes.io/component=api
        kubectl logs deployment/api --tail=50
        exit 1
    fi

    log "Deployment is ready"
}

# Run health checks
run_health_checks() {
    log "Running health checks..."

    # Get service URL
    SERVICE_URL=$(kubectl get ingress api-ingress -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo "")

    if [ -n "$SERVICE_URL" ]; then
        log "Testing health endpoint: https://$SERVICE_URL/health"

        # Wait a moment for the ingress to update
        sleep 10

        # Test health endpoint
        if curl -f -s "https://$SERVICE_URL/health" > /dev/null; then
            log "✅ Health check passed"
        else
            warn "⚠️  Health check failed - may still be propagating"
        fi

        # Test readiness endpoint
        if curl -f -s "https://$SERVICE_URL/ready" > /dev/null; then
            log "✅ Readiness check passed"
        else
            warn "⚠️  Readiness check failed - may still be propagating"
        fi

        log "🌐 Service URL: https://$SERVICE_URL"
    else
        warn "Could not determine service URL (ingress may not be ready)"
    fi
}

# Show deployment status
show_status() {
    log "Deployment Status:"
    kubectl get pods -l app.kubernetes.io/component=api
    kubectl get services
    kubectl get ingress

    log "Recent logs:"
    kubectl logs deployment/api --tail=20
}

# Cleanup function
cleanup() {
    if [ $? -ne 0 ]; then
        error "Deployment failed. Cleaning up..."
        # Add any cleanup logic here if needed
    fi
}

# Main execution
main() {
    trap cleanup ERR

    log "Starting deployment to staging..."

    check_prerequisites
    build_and_push
    deploy_to_k8s
    wait_for_deployment
    run_health_checks
    show_status

    log "🎉 Deployment completed successfully!"
    log "Image tag: $IMAGE_TAG"
    log "Namespace: $NAMESPACE"

    if [ -n "$SERVICE_URL" ]; then
        log "API URL: https://$SERVICE_URL"
    fi
}

# Run main function
main "$@"