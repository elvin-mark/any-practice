# Task: Troubleshoot Service Missing Endpoints

The frontend application cannot reach the authentication service `auth-api` in the `security` namespace.

1. Inspect the Service `auth-api` and its endpoints in namespace `security`.
2. Inspect the backing pods running in namespace `security`.
3. Identify why no endpoints are registered for `auth-api`.
4. Fix the issue without recreating the deployment.
5. Verify that `auth-api` now has active endpoints corresponding to the running pods.
