# CKA Interactive Mock Simulator & Practice Labs

An interactive, command-line testing framework for the **Certified Kubernetes Administrator (CKA)** exam, built with Go and powered by local **Kind** (Kubernetes in Docker) clusters.

---

## 🎯 Features

- **Full Exam Domain Coverage (100 Points Total)**: 20 comprehensive scenarios accurately weighted according to the Linux Foundation / CNCF CKA exam specification:
  - 🏛 **Cluster Architecture, Installation & Configuration (25%)**
  - 📦 **Workloads & Scheduling (15%)**
  - 🌐 **Services & Networking (20%)**
  - 💾 **Storage (10%)**
  - 🔧 **Troubleshooting (30%)**
- **Dual Practice Modes**:
  - **Drill / Practice Mode**: Target specific topics one-by-one with immediate automated grading, progressive hints, and complete step-by-step solutions.
  - **Timed 2-Hour Mock Exam Simulation**: 20-question randomized exam session with active countdown timer, question navigation (`next`, `prev`, `goto`), automated grading, and final domain-by-domain scorecard (66% passing threshold) exported to Markdown.
- **Dynamic Scenario Discovery**: Scenarios are defined externally in `scenarios/` as clean YAML manifests, markdown prompts, and bash verification scripts (`setup.sh`, `verify.sh`, `solution.md`).
- **Smart Cluster Management**: Automatically configures multi-node Kind clusters with **Calico CNI** (for genuine packet-level NetworkPolicy enforcement) and seamlessly falls back to existing active clusters (like `kind-dev-cluster`).

---

## 🚀 Quickstart

### 1. Build and Install
The CLI is written in Go:
```bash
make build
make install  # Copies binary to ~/.local/bin/cka
```

### 2. Prepare Kind Clusters
To create dedicated clusters with Calico CNI:
```bash
cka cluster init
```
Or check cluster status:
```bash
cka cluster status
```
*(If you already have a running Kind cluster, `cka` will automatically discover and use your active context!)*

### 3. List All Available Scenarios
```bash
cka list
```
Filter by domain:
```bash
cka list --domain Troubleshooting
```

---

## 📖 Practice Workflow

### Step 1: Start a Scenario
```bash
cka practice arch-01-rbac
```
This prepares the cluster state, creates necessary namespaces, switches kubectl context, and displays the task instructions.

### Step 2: Solve the Task
Use standard `kubectl` commands to solve the question in your terminal.

### Step 3: Grade Your Work
```bash
cka verify
```
`cka` executes granular assertion checks against the cluster and displays `[PASS]` / `[FAIL]` results for every requirement.

### Step 4: Need Help?
- **View hints**: `cka hint`
- **View step-by-step solution**: `cka solve`
- **Reset to starting state**: `cka reset`

---

## ⏱ Timed Mock Exam Simulation

Simulate the real CKA 2-hour exam experience:

```bash
# Start a 120-minute timed exam
cka exam start --duration 120

# Check active timer and progress
cka exam status

# Navigate questions
cka exam next
cka exam prev
cka exam goto 5

# Grade current question
cka verify

# Finish exam and get full score report
cka exam finish
```

Upon completion, `cka` generates a terminal summary and saves an exportable Markdown report (e.g. `cka-exam-report-YYYY-MM-DD-HHMMSS.md`) with domain percentages and pass/fail analysis.

---

## 📚 Curriculum Catalog (20 Scenarios)

| ID | Domain | Weight | Difficulty | Scenario Title |
|---|---|---|---|---|
| `arch-01-rbac` | Architecture & Config | 6 pts | Medium | Configure RBAC for Developer Role |
| `arch-02-static-pod` | Architecture & Config | 6 pts | Medium | Create a Static Pod on Control Plane Node |
| `arch-03-etcd-backup` | Architecture & Config | 7 pts | Hard | Take an etcd Snapshot Backup on Control Plane |
| `arch-04-node-maintenance` | Architecture & Config | 6 pts | Medium | Safely Drain and Maintain a Cluster Node |
| `workloads-01-multi-container` | Workloads & Scheduling | 3 pts | Easy | Deploy a Multi-Container Pod with Shared Volume |
| `workloads-02-node-affinity` | Workloads & Scheduling | 3 pts | Medium | Schedule Workload Using Node Affinity |
| `workloads-03-taints-tolerations` | Workloads & Scheduling | 3 pts | Easy | Configure Taints and Pod Tolerations |
| `workloads-04-daemonset` | Workloads & Scheduling | 3 pts | Medium | Deploy a Cluster-Wide DaemonSet with Tolerations |
| `workloads-05-cronjob` | Workloads & Scheduling | 3 pts | Easy | Create a Batch CronJob with Limits |
| `net-01-networkpolicy` | Services & Networking | 7 pts | Hard | Restrict Ingress Traffic with NetworkPolicy |
| `net-02-clusterip-nodeport` | Services & Networking | 4 pts | Easy | Expose Deployment via NodePort Service |
| `net-03-ingress` | Services & Networking | 5 pts | Medium | Configure Ingress with Path-Based Routing |
| `net-04-coredns-troubleshoot` | Services & Networking | 4 pts | Medium | Troubleshoot Cluster DNS Resolution Failure |
| `storage-01-pv-pvc` | Storage | 4 pts | Easy | Create PersistentVolume and PersistentVolumeClaim |
| `storage-02-pod-mount` | Storage | 3 pts | Easy | Mount Persistent Storage into a Pod |
| `storage-03-storageclass` | Storage | 3 pts | Easy | Create a Custom StorageClass |
| `trouble-01-crashloop` | Troubleshooting | 7 pts | Medium | Troubleshoot Application Stuck in CrashLoopBackOff |
| `trouble-02-worker-kubelet` | Troubleshooting | 8 pts | Hard | Troubleshoot Worker Node in NotReady Status |
| `trouble-03-apiserver-manifest` | Troubleshooting | 8 pts | Hard | Troubleshoot Control Plane Static Pod Failure |
| `trouble-04-endpoints` | Troubleshooting | 7 pts | Easy | Troubleshoot Service Missing Endpoints |

---

## 🛠 Adding Custom Scenarios

Each scenario in `scenarios/<id>/` contains:
- `metadata.yaml`: Scenario ID, title, domain, difficulty, weight, target context.
- `problem.md`: Task instructions.
- `setup.sh`: Idempotent script initializing the scenario state.
- `verify.sh`: Grading script printing `[PASS]` or `[FAIL]` checks and exiting `0` or `1`.
- `solution.md`: Step-by-step walkthrough.

You can drop new scenario folders into `scenarios/` at any time without recompiling the binary!
