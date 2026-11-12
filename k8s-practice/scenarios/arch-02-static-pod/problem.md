# Task: Create a Static Pod on Control Plane Node

Create a Static Pod on the control-plane node of the cluster with the following specifications:

- Pod Name: `my-static-web` (note: Kubernetes appends the node name suffix to static pods, e.g. `my-static-web-<node-name>`)
- Namespace: `default`
- Container Image: `nginx:1.25-alpine`
- Restart Policy: `Always`
- Target Node: The control-plane node (access using `docker exec -it <control-plane-node> bash` or `sh`)
- Manifest Path: Must be placed under the kubelet's static pod directory (`/etc/kubernetes/manifests/my-static-web.yaml`)

Ensure the static pod starts successfully and enters `Running` state.
