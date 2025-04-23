#!/bin/bash

# Script to to point k3s to the docker registry running on the host
# This is used to avoid hitting dockerhub rate limits on E2E runners
ip_addr=$1

mkdir -p /etc/k3s/
echo "mirrors:
  docker.io:
    endpoint:
      - \"http://$ip_addr:5000\"" >> /etc/k3s/registries.yaml
