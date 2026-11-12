# Task: Deploy a Multi-Container Pod with Shared Volume

Create a Pod named `multi-app` in the `default` namespace containing two containers:

1. **Primary Container**:
   - Name: `producer`
   - Image: `busybox:1.36`
   - Command: `sh -c 'while true; do date >> /var/log/app.log; sleep 2; done'`
   - Volume Mount: mount volume `shared-logs` at `/var/log`

2. **Sidecar Container**:
   - Name: `consumer`
   - Image: `busybox:1.36`
   - Command: `sh -c 'tail -n+1 -F /var/log/app.log'`
   - Volume Mount: mount volume `shared-logs` at `/var/log`

3. **Volume**:
   - Name: `shared-logs`
   - Type: `emptyDir: {}`

Both containers must reach the `Running` state and the sidecar must be able to read output from the shared file.
