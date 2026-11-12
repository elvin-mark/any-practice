# Task: Expose Deployment via NodePort Service

A deployment named `front-web` runs in the `web-apps` namespace.

1. Ensure the namespace `web-apps` exists and contains the deployment `front-web`.
2. Expose the `front-web` deployment as a Service with:
   - Name: `front-web-svc`
   - Namespace: `web-apps`
   - Service Type: `NodePort`
   - Service Port: `80`
   - Target Port: `80`
   - Node Port: `30080`
3. Verify that the endpoints for `front-web-svc` match the ready pod IPs and that accessing nodePort 30080 responds successfully.
