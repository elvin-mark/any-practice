# Task: Volume Backup & Restoration Pipeline

Docker volumes store critical database files and user assets. Because volumes reside inside `/var/lib/docker/volumes/`, disaster recovery strategies require backing up volume contents to portable compressed archives and restoring them into fresh recovery volumes.

### Objectives:
1. An existing volume named `prod-backup-source` has been provisioned containing configuration and database dumps.
2. Back up the contents of `prod-backup-source` into an archive file at `/tmp/prod-backup.tar.gz` using an ephemeral Docker container.
3. Create a new volume named `prod-backup-restored` labeled `practice=docker`.
4. Restore the contents of `/tmp/prod-backup.tar.gz` into `prod-backup-restored` using another ephemeral container.
5. All original files (`config.json`, `data.db`, and `checksum.sha256`) must exist in `prod-backup-restored` and their SHA256 hashes must match the source volume.
