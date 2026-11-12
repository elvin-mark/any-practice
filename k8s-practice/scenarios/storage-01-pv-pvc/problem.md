# Task: Create PersistentVolume and PersistentVolumeClaim

Configure persistent storage for workloads:

1. Create a `PersistentVolume` named `data-pv`:
   - Storage Capacity: `1Gi`
   - Access Mode: `ReadWriteOnce`
   - Volume Type: `hostPath` at `/mnt/data`
   - StorageClass: `manual`
2. Create a `PersistentVolumeClaim` named `data-pvc` in the `default` namespace:
   - Storage Request: `1Gi`
   - Access Mode: `ReadWriteOnce`
   - StorageClass: `manual`
3. Verify that the PVC successfully binds to the PV and shows `Status: Bound`.
