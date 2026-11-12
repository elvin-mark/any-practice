#!/usr/bin/env bash
set -e

CP_CONTAINER=$(docker ps --filter "name=control-plane" --format "{{.Names}}" 2>/dev/null | head -n 1 || true)
if [ -n "$CP_CONTAINER" ]; then
  # Inject broken static pod manifest with invalid command
  docker exec "$CP_CONTAINER" sh -c 'cat << "EOF" > /etc/kubernetes/manifests/infra-agent.yaml
apiVersion: v1
kind: Pod
metadata:
  name: infra-agent
  namespace: kube-system
spec:
  restartPolicy: Always
  containers:
  - name: agent
    image: busybox:1.36
    command: ["invalid-command-does-not-exist"]
EOF'
  echo "Injected broken static pod manifest into /etc/kubernetes/manifests/infra-agent.yaml"
fi
