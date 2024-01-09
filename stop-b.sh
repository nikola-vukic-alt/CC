#!/bin/bash

# Stop Ingress
kubectl delete -f k8s/ingress/library-ingress.yaml
echo

# Stop Services
kubectl delete -f k8s/services/central-service.yaml
kubectl delete -f k8s/services/ns-service.yaml
kubectl delete -f k8s/services/bg-service.yaml
kubectl delete -f k8s/services/nis-service.yaml
echo

# Stop Deployments
kubectl delete -f k8s/deployments/central-deployment.yaml
kubectl delete -f k8s/deployments/ns-deployment.yaml
kubectl delete -f k8s/deployments/bg-deployment.yaml
kubectl delete -f k8s/deployments/nis-deployment.yaml
echo

# Stop Persistent Volumes and Claims
kubectl delete -f k8s/volumes/central-host-pv.yaml
kubectl delete -f k8s/volumes/bg-host-pv.yaml
kubectl delete -f k8s/volumes/ns-host-pv.yaml
kubectl delete -f k8s/volumes/nis-host-pv.yaml
echo

echo "All resources stopped."
