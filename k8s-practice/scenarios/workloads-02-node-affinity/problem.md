# Task: Schedule Workload Using Node Affinity

Configure node labeling and schedule a deployment using hard Node Affinity:

1. Label one of your worker nodes with the label `disktype=ssd`.
2. Create a Deployment named `ssd-app` in the `default` namespace with:
   - Replicas: `2`
   - Image: `nginx:1.25-alpine`
   - Pod label: `app=ssd-app`
   - Node Affinity: Must enforce `requiredDuringSchedulingIgnoredDuringExecution` requiring nodes to match the key `disktype` with value `ssd`.
3. Verify that both pods are scheduled strictly on the node with `disktype=ssd`.
