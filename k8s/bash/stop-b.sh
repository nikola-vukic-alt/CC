#!/bin/bash

# Stop Ingress
kubectl delete -f ../ingress/library-ingress.yaml
echo

# Stop Services
kubectl delete -f ../services/central-service.yaml
kubectl delete -f ../services/ns-service.yaml
kubectl delete -f ../services/bg-service.yaml
kubectl delete -f ../services/nis-service.yaml
echo

# Stop Deployments
kubectl delete -f ../deployments/central-deployment.yaml
kubectl delete -f ../deployments/ns-deployment.yaml
kubectl delete -f ../deployments/bg-deployment.yaml
kubectl delete -f ../deployments/nis-deployment.yaml
echo

# Stop Persistent Volumes and Claims
kubectl delete -f ../volumes/central-host-pv.yaml
kubectl delete -f ../volumes/bg-host-pv.yaml
kubectl delete -f ../volumes/ns-host-pv.yaml
kubectl delete -f ../volumes/nis-host-pv.yaml
echo

echo "All resources stopped."
