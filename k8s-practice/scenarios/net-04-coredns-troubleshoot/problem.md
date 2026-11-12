# Task: Troubleshoot Cluster DNS Resolution Failure

Applications inside the cluster report that internal domain names (such as `kubernetes.default.svc.cluster.local`) fail to resolve.

1. Test DNS resolution using the provided `dns-client` pod in namespace `default`:
   `kubectl exec dns-client -n default -- nslookup kubernetes.default`
2. Investigate cluster DNS infrastructure components in namespace `kube-system`.
3. Identify and fix the root cause preventing DNS resolution.
4. Verify that `kubectl exec dns-client -n default -- nslookup kubernetes.default` now succeeds and resolves the cluster IP of the Kubernetes API service.
