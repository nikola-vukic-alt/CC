#!/bin/bash

# Apply Volumes and Claims
kubectl apply -f volumes/central-host-pv.yaml
kubectl apply -f volumes/ns-host-pv.yaml
kubectl apply -f volumes/bg-host-pv.yaml
kubectl apply -f volumes/nis-host-pv.yaml
echo

# Apply Deployments
kubectl apply -f deployments/central-deployment.yaml
kubectl apply -f deployments/ns-deployment.yaml 
kubectl apply -f deployments/bg-deployment.yaml 
kubectl apply -f deployments/nis-deployment.yaml 
echo

# Apply Services
kubectl apply -f services/central-service.yaml
kubectl apply -f services/ns-service.yaml
kubectl apply -f services/bg-service.yaml
kubectl apply -f services/nis-service.yaml
echo

# Apply Ing
kubectl apply -f ingress/library-ingress.yaml
echo 

echo "All resources applied."