# Task: Deploy a Cluster-Wide DaemonSet with Tolerations

You are requested to deploy a node monitoring daemon across all nodes in the cluster.

1. Create a `DaemonSet` named `node-monitor` in the `default` namespace.
2. The DaemonSet specification must satisfy:
   - Pod label: `app=node-monitor`
   - Container name: `monitor`
   - Container image: `busybox:1.36`
   - Command: `sh -c 'while true; do echo monitoring; sleep 30; done'`
   - Tolerations: Must contain toleration(s) allowing pods to run on **all** nodes, including the control-plane node (which carries `node-role.kubernetes.io/control-plane:NoSchedule`).
3. Ensure one pod is scheduled and running on each node in the cluster.
