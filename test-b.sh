#!/bin/bash

CENTRAL_SERVICE_NAME="central-library-service"
NS_SERVICE_NAME="ns-library-service"
BG_SERVICE_NAME="bg-library-service"
NIS_SERVICE_NAME="nis-library-service"

send_request() {
  local SERVICE_NAME=$1
  local PORT=$2

  curl -X POST \
       -H "Content-Type: application/json" \
       -d '{"SSN": "123", "Title": "Anna Karenina", "Author": "L. N. Tolstoy", "ISBN": "420-3-69-148410-0"}' \
       http://localhost:${PORT}/borrow
  echo
  sleep 1

  # Add more requests as needed...

  echo "Requests to ${SERVICE_NAME} completed."
  sleep 1
}

# Deploy the application in Kubernetes
kubectl apply -f path/to/your/kubernetes/yaml/files

# Sleep to allow pods to start (you can adjust the duration)
sleep 10

echo "Trying to borrow a book while not being a member"
send_request ${NS_SERVICE_NAME} 8081

echo "Trying to return the book while not being a member"
# Add similar calls for other services...

# More requests...

# Tear down the Kubernetes deployment
kubectl delete -f path/to/your/kubernetes/yaml/files

# Sleep to allow resources to be deleted (you can adjust the duration)
sleep 10

echo "Testing done..."
