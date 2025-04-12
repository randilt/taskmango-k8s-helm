# TaskMango Helm Chart

This Helm chart deploys the TaskMango microservices application on Kubernetes.

## Prerequisites

- Kubernetes cluster (Minikube recommended for local development)
- Helm 3 installed
- Docker images for api-service and auth-service built and pushed to a registry

## Installation

1. Build and push your Docker images:
   ```bash
   docker build -t your-repo/api-service:latest -f api-service/Dockerfile .
   docker build -t your-repo/auth-service:latest -f auth-service/Dockerfile .
   docker push your-repo/api-service:latest
   docker push your-repo/auth-service:latest
   ```
