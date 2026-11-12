# Task: Configure RBAC for Developer Role

You are tasked with setting up Role-Based Access Control (RBAC) in the `development` namespace.

1. Ensure the namespace `development` exists.
2. Create a `Role` named `developer-role` in namespace `development` with the following permissions:
   - Resources: `pods`, `services`
   - Verbs: `get`, `list`, `watch`, `create`, `update`, `delete`
3. Create a `RoleBinding` named `developer-binding` in namespace `development`:
   - Role: `developer-role`
   - User: `developer`
4. Verify that the user `developer` can list pods in `development`, but **cannot** list pods in the `kube-system` namespace.
