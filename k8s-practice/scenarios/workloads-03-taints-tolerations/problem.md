# Task: Configure Taints and Pod Tolerations

1. Identify the worker node in the cluster.
2. Add a taint to the worker node:
   - Key: `tier`
   - Value: `critical`
   - Effect: `NoSchedule`
3. Create a Pod named `toleration-pod` in the `default` namespace:
   - Image: `nginx:1.25-alpine`
   - Configure a toleration to allow scheduling on nodes with the `tier=critical:NoSchedule` taint.
4. Verify that the pod is successfully scheduled on the tainted node and running.
