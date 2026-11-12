# Task: Safely Drain and Maintain a Cluster Node

A worker node in the cluster requires hardware maintenance.

1. Identify the worker node in the cluster (e.g. node without control-plane role).
2. Drain the worker node safely:
   - Evict all workload pods.
   - Ignore DaemonSets.
   - Allow eviction of pods using `emptyDir` volumes.
3. Verify that the node is marked as `SchedulingDisabled`.
4. Ensure a sample deployment named `maintenance-test` in namespace `default` continues running with all ready replicas scheduled away from the drained node.
5. Once maintenance is approved, mark the node schedulable again using `kubectl uncordon`.
