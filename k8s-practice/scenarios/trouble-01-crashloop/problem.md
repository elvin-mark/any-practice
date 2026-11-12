# Task: Troubleshoot Application Stuck in CrashLoopBackOff

A microservice deployment named `order-service` deployed in namespace `finance` is failing to start.

1. Inspect the pods of `order-service` in namespace `finance`.
2. Check container termination reasons and logs to determine why the pods are entering `CrashLoopBackOff`.
3. Repair the `order-service` deployment so that:
   - The container runs without crashing.
   - All 2 replicas achieve the `Ready` and `Running` status.
