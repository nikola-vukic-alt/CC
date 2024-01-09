#!/bin/bash

# Start the minikube
minikube start

# Apply Configs
kubectl apply -f k8s/config/central-config.yaml
kubectl apply -f k8s/config/ns-config.yaml
kubectl apply -f k8s/config/bg-config.yaml
kubectl apply -f k8s/config/nis-config.yaml
echo 

# Apply Volumes and Claims
kubectl apply -f k8s/volumes/central-host-pv.yaml
kubectl apply -f k8s/volumes/ns-host-pv.yaml
kubectl apply -f k8s/volumes/bg-host-pv.yaml
kubectl apply -f k8s/volumes/nis-host-pv.yaml
echo

# Apply Deployments
kubectl apply -f k8s/deployments/central-deployment.yaml
kubectl apply -f k8s/deployments/ns-deployment.yaml 
kubectl apply -f k8s/deployments/bg-deployment.yaml 
kubectl apply -f k8s/deployments/nis-deployment.yaml 
echo

# Apply Services
kubectl apply -f k8s/services/central-service.yaml
kubectl apply -f k8s/services/ns-service.yaml
kubectl apply -f k8s/services/bg-service.yaml
kubectl apply -f k8s/services/nis-service.yaml
echo

# Apply Ingress
kubectl apply -f k8s/ingress/library-ingress.yaml
echo 

echo "All resources applied."