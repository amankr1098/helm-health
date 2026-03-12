# Kubernetes Resource Health Rules

This document defines the rules used to determine the health status of Kubernetes resources.

---

## Deployment

| Condition | Status |
|-----------|--------|
| `ReadyReplicas == DesiredReplicas` AND `AvailableReplicas == DesiredReplicas` | ✅ Healthy |
| `ReadyReplicas < DesiredReplicas` | ⚠️ Degraded |
| `AvailableReplicas == 0` AND `DesiredReplicas > 0` | ❌ Unhealthy |
| `UpdatedReplicas < DesiredReplicas` | 🔄 Updating |
| Condition `Progressing` is `False` with reason `ProgressDeadlineExceeded` | ❌ Failed |

---

## Pod

| Condition | Status |
|-----------|--------|
| `Phase == Running` AND all containers `Ready == true` | ✅ Healthy |
| `Phase == Succeeded` | ✅ Completed |
| `Phase == Pending` | ⏳ Pending |
| `Phase == Failed` | ❌ Failed |
| `Phase == Unknown` | ❓ Unknown |
| Any container `RestartCount > threshold` (e.g., 5) | ⚠️ CrashLooping |
| Container state `Waiting` with reason `CrashLoopBackOff` | ❌ CrashLooping |
| Container state `Waiting` with reason `ImagePullBackOff` | ❌ ImagePullError |
| Container state `Waiting` with reason `ErrImagePull` | ❌ ImagePullError |

---

## StatefulSet

| Condition | Status |
|-----------|--------|
| `ReadyReplicas == DesiredReplicas` | ✅ Healthy |
| `ReadyReplicas < DesiredReplicas` | ⚠️ Degraded |
| `CurrentReplicas != UpdatedReplicas` | 🔄 Updating |
| `ReadyReplicas == 0` AND `DesiredReplicas > 0` | ❌ Unhealthy |

---

## DaemonSet

| Condition | Status |
|-----------|--------|
| `NumberReady == DesiredNumberScheduled` | ✅ Healthy |
| `NumberReady < DesiredNumberScheduled` | ⚠️ Degraded |
| `NumberMisscheduled > 0` | ⚠️ Misscheduled |
| `NumberUnavailable > 0` | ⚠️ Unavailable |
| `NumberReady == 0` AND `DesiredNumberScheduled > 0` | ❌ Unhealthy |

---

## ReplicaSet

| Condition | Status |
|-----------|--------|
| `ReadyReplicas == DesiredReplicas` | ✅ Healthy |
| `ReadyReplicas < DesiredReplicas` | ⚠️ Degraded |
| `ReadyReplicas == 0` AND `DesiredReplicas > 0` | ❌ Unhealthy |

---

## Job

| Condition | Status |
|-----------|--------|
| `Succeeded >= Completions` | ✅ Completed |
| `Active > 0` AND `Failed == 0` | 🔄 Running |
| `Failed > 0` AND `Active > 0` | ⚠️ PartiallyFailed |
| `Failed > BackoffLimit` | ❌ Failed |
| Condition `Complete` is `True` | ✅ Completed |
| Condition `Failed` is `True` | ❌ Failed |

---

## CronJob

| Condition | Status |
|-----------|--------|
| `LastSuccessfulTime` exists AND recent | ✅ Healthy |
| `LastScheduleTime` exists AND no active jobs failing | ✅ Healthy |
| Active jobs count > `ConcurrencyPolicy` limit (if applicable) | ⚠️ Backlog |
| `LastSuccessfulTime` is stale (older than 2x schedule interval) | ⚠️ Stale |
| No `LastScheduleTime` AND should have run | ❌ NotScheduled |

---

## Service

| Condition | Status |
|-----------|--------|
| Has matching Endpoints with `Ready` addresses | ✅ Healthy |
| Has Endpoints but all are `NotReady` | ⚠️ NoReadyEndpoints |
| No matching Endpoints found | ❌ NoEndpoints |
| Type `LoadBalancer` AND `Ingress` IP/Hostname assigned | ✅ Healthy |
| Type `LoadBalancer` AND no `Ingress` assigned | ⏳ Pending |

---

## PersistentVolumeClaim (PVC)

| Condition | Status |
|-----------|--------|
| `Phase == Bound` | ✅ Healthy |
| `Phase == Pending` | ⏳ Pending |
| `Phase == Lost` | ❌ Lost |

---

## PersistentVolume (PV)

| Condition | Status |
|-----------|--------|
| `Phase == Bound` | ✅ Healthy |
| `Phase == Available` | ✅ Available |
| `Phase == Released` | ⚠️ Released |
| `Phase == Failed` | ❌ Failed |

---

## Ingress

| Condition | Status |
|-----------|--------|
| All backend services have ready endpoints | ✅ Healthy |
| LoadBalancer `Ingress` IP/Hostname assigned | ✅ Configured |
| Some backend services have no ready endpoints | ⚠️ PartiallyHealthy |
| No LoadBalancer `Ingress` assigned | ⏳ Pending |
| No ready endpoints on any backend | ❌ Unhealthy |

---

## NetworkPolicy

| Condition | Status |
|-----------|--------|
| Policy exists and matches at least one pod | ✅ Healthy |
| Policy exists but no pods match the selector | ⚠️ NoMatchingPods |
| Failed to fetch policy | ❓ Unknown |

---

## HorizontalPodAutoscaler (HPA)

| Condition | Status |
|-----------|--------|
| `CurrentReplicas == DesiredReplicas` | ✅ Stable |
| `CurrentReplicas < DesiredReplicas` | 🔼 ScalingUp |
| `CurrentReplicas > DesiredReplicas` | 🔽 ScalingDown |
| Condition `ScalingActive` is `False` | ⚠️ ScalingDisabled |
| Condition `AbleToScale` is `False` | ❌ Unable |

---

## Node

| Condition | Status |
|-----------|--------|
| Condition `Ready` is `True` | ✅ Healthy |
| Condition `Ready` is `False` | ❌ NotReady |
| Condition `Ready` is `Unknown` | ❓ Unknown |
| Condition `MemoryPressure` is `True` | ⚠️ MemoryPressure |
| Condition `DiskPressure` is `True` | ⚠️ DiskPressure |
| Condition `PIDPressure` is `True` | ⚠️ PIDPressure |
| Condition `NetworkUnavailable` is `True` | ❌ NetworkUnavailable |

---

## Health Status Definitions

| Status | Icon | Meaning |
|--------|------|---------|
| **Healthy** | ✅ | Resource is fully operational |
| **Degraded** | ⚠️ | Resource is partially working |
| **Unhealthy** | ❌ | Resource has failed |
| **Pending** | ⏳ | Resource is waiting |
| **Updating** | 🔄 | Resource is being updated |
| **Completed** | ✅ | Job/Task finished successfully |
| **Failed** | ❌ | Job/Task failed |
| **Unknown** | ❓ | Status cannot be determined |

---

## Aggregated Release Health

For a Helm release containing multiple resources:

| Condition | Overall Status |
|-----------|----------------|
| All resources are Healthy/Completed | ✅ Healthy |
| Any resource is Unhealthy/Failed | ❌ Unhealthy |
| Any resource is Degraded (but none Failed) | ⚠️ Degraded |
| Any resource is Pending/Updating (but none Failed/Degraded) | 🔄 Progressing |

---

## Notes

1. **Thresholds are configurable** - Restart count thresholds, staleness intervals, etc. should be configurable.
2. **Events matter** - Recent warning events can indicate issues not visible in status.
3. **Age matters** - A Pending pod for 1 minute is normal; Pending for 30 minutes is a problem.
4. **Context matters** - A Job with Failed pods might still be healthy if retries are expected.
