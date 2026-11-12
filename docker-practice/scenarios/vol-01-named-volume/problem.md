# Task: Persistent Storage with Named Volumes

Container filesystems are ephemeral: any data written to the container's writable layer is permanently lost when the container is deleted. Docker named volumes isolate data storage from the container lifecycle.

### Objectives:
1. Create a Docker named volume named `db-store` labeled `practice=docker`.
2. Run a container named `datastore-app` using `alpine:3.20` in the background (`sleep 3600`).
3. Mount the `db-store` volume to `/var/lib/appdata` inside the container.
4. Inside `/var/lib/appdata`, create a file named `persistent.txt` containing the exact text `DOCKER_PERSISTENCE_VERIFIED`.
5. The data must persist inside the volume independently of the container.
