# Task: Troubleshoot Control Plane Static Pod Failure

A critical cluster management agent named `infra-agent` on the control-plane node fails to start.

1. Connect to the control plane node (e.g. `docker exec -it <control-plane-node> bash`).
2. Inspect the static pod manifest files located in `/etc/kubernetes/manifests/`.
3. Locate `/etc/kubernetes/manifests/infra-agent.yaml` and identify the configuration error preventing kubelet from running the container.
4. Correct the manifest error.
5. Verify using `kubectl get pods -n kube-system` that `infra-agent-<control-plane-node>` enters the `Running` state.
