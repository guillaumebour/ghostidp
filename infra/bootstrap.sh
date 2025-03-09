#!/usr/bin/env bash

# This script bootstraps a local Kubernetes Cluster
# with nginx-gateway-fabric installed.

echo "Bootstrapping a local dev cluster..."

# Start a local Kind cluster
ctlptl apply -f kind.yaml

# Install nginx-gateway-fabric

## Gateway API resources
kubectl kustomize "https://github.com/nginx/nginx-gateway-fabric/config/crd/gateway-api/standard?ref=v1.6.1" | kubectl apply -f -

## Nginx Gateway Fabric helm Chart
helm install ngf oci://ghcr.io/nginx/charts/nginx-gateway-fabric --create-namespace -n nginx-gateway --set service.create=false

## Add a NodePort svc so that ports 80/443 on the host target nginx-gateway-fabric
kubectl apply -f nodeport-config.yaml

# Create a Gateway with a default http -> https redirect
kubectl apply -f gateway-tls-cert-secret.yaml  # Gateway certificates for TLS Termination
kubectl apply -f gateway.yaml                  # Gateway
kubectl apply -f http_https_redirect.yaml      # HTTP -> HTTPS Redirect

echo "The local dev cluster is ready!"
