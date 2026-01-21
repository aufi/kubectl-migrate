#!/bin/bash

echo "# Using kubectl context: $(kubectl config current-context)"

# Delete all resources
kubectl delete -k ./
