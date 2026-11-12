# Task: Troubleshoot Worker Node in NotReady Status

A worker node in the Kubernetes cluster is showing `NotReady` status, preventing pods from scheduling onto it.

1. Identify the worker node that is `NotReady` using `kubectl get nodes`.
2. Inspect the node events using `kubectl describe node <worker-node>`.
3. Access the worker node (e.g. `docker exec -it <worker-node> bash` or via SSH).
4. Investigate the `kubelet` system service logs using `journalctl -u kubelet` or `systemctl status kubelet`.
5. Fix the issue causing kubelet failure and start the service.
6. Verify that the node returns to `Ready` status and is able to accept pods.
