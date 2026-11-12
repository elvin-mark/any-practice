# Task: Create a Custom StorageClass

Create a custom StorageClass with the following specifications:

- Name: `fast-retain`
- Provisioner: `rancher.io/local-path`
- Reclaim Policy: `Retain`
- Volume Binding Mode: `WaitForFirstConsumer`
- Allow Volume Expansion: `true`

Verify that the StorageClass is registered and visible via `kubectl get sc`.
