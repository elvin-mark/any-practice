# Task: Mount Persistent Storage into a Pod

A PersistentVolumeClaim named `storage-claim` is already bound in namespace `default`.

1. Create a Pod named `storage-pod` in namespace `default`:
   - Image: `busybox:1.36`
   - Command: `sh -c 'echo "cka-practice-success" > /opt/store/status.txt; sleep 3600'`
   - Volume: Mount the existing PVC `storage-claim` to mount path `/opt/store`.
2. Ensure the pod reaches `Running` state and that `/opt/store/status.txt` exists with the expected text.
