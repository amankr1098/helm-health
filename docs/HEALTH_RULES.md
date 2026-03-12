# Health Rules

This document defines how **helm-health** determines the health of each Kubernetes resource in a Helm release.

Every resource is assigned one of three statuses:

| Status | Icon | Meaning |
|--------|------|---------|
| **Healthy** | ✅ | Resource is operating as expected |
| **Unhealthy** | ❌ | Resource has a problem that needs attention |
| **Unknown** | ❓ | Status could not be determined (API fetch error) |

If a resource kind has no dedicated rules (ConfigMap, Secret, ServiceAccount, RBAC resources, etc.), it is reported as **Healthy** as long as it exists in the release manifest.

---

## Deployment

| Condition | Result |
|-----------|--------|
| `ReadyReplicas >= DesiredReplicas` and `DesiredReplicas > 0` | ✅ Healthy — *"3/3 replicas ready"* |
| `DesiredReplicas == 0` (scaled to zero) | ✅ Healthy — *"scaled to 0"* |
| `ReadyReplicas < DesiredReplicas` | ❌ Unhealthy — reports available vs desired count |
| Unable to fetch Deployment from API | ❓ Unknown |

**Pod diagnostics:** When a Deployment is unhealthy, the plugin fetches the Deployment's pods and inspects each container. Reported details include:

- Container waiting state and reason (`CrashLoopBackOff`, `ImagePullBackOff`, etc.)
- Container terminated state, exit code, and reason
- Time since last restart

---

## StatefulSet

| Condition | Result |
|-----------|--------|
| `ReadyReplicas >= DesiredReplicas` and `DesiredReplicas > 0` | ✅ Healthy — *"2/2 replicas ready"* |
| `DesiredReplicas == 0` (scaled to zero) | ✅ Healthy — *"scaled to 0"* |
| `ReadyReplicas < DesiredReplicas` | ❌ Unhealthy — reports ready vs desired count |
| Unable to fetch StatefulSet from API | ❓ Unknown |

Health data includes `desired`, `ready`, `current`, and `updated` replica counts.

---

## DaemonSet

| Condition | Result |
|-----------|--------|
| `NumberReady >= DesiredNumberScheduled` and `DesiredNumberScheduled > 0` | ✅ Healthy — *"5/5 nodes ready"* |
| `DesiredNumberScheduled == 0` | ✅ Healthy — *"no nodes scheduled"* |
| `NumberReady < DesiredNumberScheduled` | ❌ Unhealthy — reports ready vs desired count |
| `NumberMisscheduled > 0` (while unhealthy) | ❌ Unhealthy — additionally reports misscheduled node count |
| Unable to fetch DaemonSet from API | ❓ Unknown |

Health data includes `desired`, `ready`, `available`, and `misscheduled` counts.

---

## Pod

| Condition | Result |
|-----------|--------|
| `Phase == Running` AND all containers `Ready == true` | ✅ Healthy — *"3/3 containers ready"* |
| `Phase == Succeeded` | ✅ Healthy — *"Completed"* |
| Any other phase (`Pending`, `Failed`, `Unknown`) or any container not ready | ❌ Unhealthy — reports phase and container issues |
| Unable to fetch Pod from API | ❓ Unknown |

**Container-level issues reported:**

- Container in `Waiting` state → reports reason (e.g. `CrashLoopBackOff`, `ImagePullBackOff`) and message
- Container in `Terminated` state with non-zero exit code → reports reason, message, and exit code
- Time since last restart (if last termination state exists)

---

## Job

| Condition | Result |
|-----------|--------|
| `Succeeded >= Completions` | ✅ Healthy — *"1/1 completed"* |
| `Active > 0` AND `Failed == 0` (still running) | ✅ Healthy — *"1 active, 0/1 completed"* |
| `Failed > 0` | ❌ Unhealthy — reports failed, succeeded, and desired counts |
| No succeeded, active, or failed pods | ❓ Unknown |
| Unable to fetch Job from API | ❓ Unknown |

Health data includes `desired` completions, `succeeded`, `failed`, and `active` counts.

---

## Service

| Condition | Result |
|-----------|--------|
| Ready endpoints exist (and LB assigned if type is `LoadBalancer`) | ✅ Healthy — *"3 endpoints"* or *"3 endpoints, LoadBalancer ready"* |
| Type `LoadBalancer` but no ingress IP/hostname assigned | ❌ Unhealthy — *"LoadBalancer ingress not yet assigned"* |
| No ready endpoints and service is not `ExternalName` or headless (`ClusterIP: None`) | ❌ Unhealthy — *"No endpoints available"* |
| `ReadyEndpoints < TotalEndpoints` | ❌ Unhealthy — *"Only 2 of 3 expected endpoints available"* |
| Unable to fetch Service or Endpoints from API | ❓ Unknown |

Health data includes service `type`, `clusterIP`, and endpoint counts (`ready`, `notReady`, `total`).

> **Note:** `ExternalName` services and headless services (`ClusterIP: None`) with no endpoints are not flagged as unhealthy.

---

## PersistentVolumeClaim (PVC)

| Condition | Result |
|-----------|--------|
| `Phase == Bound` | ✅ Healthy — *"Bound"* |
| Any other phase (`Pending`, `Lost`, etc.) | ❌ Unhealthy — reports current phase |
| Unable to fetch PVC from API | ❓ Unknown |

---

## Ingress

| Condition | Result |
|-----------|--------|
| LoadBalancer IP/hostname assigned AND all backend services have ready endpoints | ✅ Healthy — *"3/3 backends ready, LB assigned"* |
| Some backends healthy but not all | ❌ Unhealthy — *"2/3 backends ready"* |
| No LoadBalancer ingress assigned | ❌ Unhealthy — *"No LoadBalancer ingress IP/hostname assigned"* |
| Backend service has no ready endpoints | ❌ Unhealthy — reports which backend has no endpoints |
| Backend service endpoints could not be fetched | ❌ Unhealthy — reports which backend failed |
| Unable to fetch Ingress from API | ❓ Unknown |

Health data includes LoadBalancer details, ingress class, TLS hosts, and per-backend status (service name, host, path, ready endpoint count).

> **Note:** Both rule-based backends and the default backend (if configured) are checked.

---

## NetworkPolicy

| Condition | Result |
|-----------|--------|
| Policy exists and at least one pod matches the pod selector | ✅ Healthy — *"5 matching pods, types: [Ingress Egress]"* |
| Policy exists but zero pods match the pod selector | ❌ Unhealthy — *"no matching pods"* |
| Unable to fetch NetworkPolicy from API | ❓ Unknown |

Health data includes the pod selector, matching pod count, policy types, and parsed ingress/egress rules (ports, pod/namespace selectors, IP blocks).

---

## Aggregated Release Health

The overall release status is derived from all individual resource statuses:

| Condition | Overall Status |
|-----------|----------------|
| All resources are Healthy | ✅ Healthy |
| Any resource is Unhealthy | ❌ Unhealthy |

When the overall status is **Unhealthy**, `helm-health` exits with a non-zero exit code.

> **Note:** If the Helm release status itself is not `deployed`, the release is immediately reported as **Unhealthy** without checking individual resources.
