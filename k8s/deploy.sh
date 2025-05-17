#!/bin/bash

docker compose build
docker compose push

kubectl delete -f services/
kubectl apply -f services/

watch kubectl get pods
