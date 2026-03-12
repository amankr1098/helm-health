<div align="center">

# ⛑️ Helm Health

**Instantly diagnose the health of your Helm releases.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/amankr1098/helm-health)](https://go.dev/)
[![License](https://img.shields.io/github/license/amankr1098/helm-health)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/amankr1098/helm-health?style=social)](https://github.com/amankr1098/helm-health/stargazers)

A Helm plugin that performs deep health checks on every Kubernetes resource in a release and gives you a clear, color-coded (or JSON) report in seconds.

</div>

---

## Why Helm Health?

`helm status` tells you a release is *deployed* — but it doesn't tell you if the underlying resources are actually running and healthy.

**helm-health** goes deeper:

- Checks **real-time Kubernetes status** of every resource in the release manifest.
- Applies **resource-specific health rules** tailored to each resource kind.
- Returns a **single overall verdict**: Healthy, Degraded, or Unhealthy.
- Supports **text** and **JSON** output for CI/CD pipelines.
- Works as a **native Helm plugin** — no extra tooling required.

---

## Demo

```
$ helm health -r my-app -n production

RELEASE: my-app
STATUS: ✓ Healthy
NAMESPACE: production

RESOURCES:
✓ Deployment/my-app-web (3/3 ready)
✓ Deployment/my-app-worker (2/2 ready)
✓ Service/my-app-web (endpoints ready)
✓ StatefulSet/my-app-redis (1/1 ready)
✓ PersistentVolumeClaim/my-app-data (Bound)
✓ Ingress/my-app-ingress (configured)
✗ Job/my-app-migration (Failed)
  backoff limit exceeded

Health check completed in 1.2s
```

---

## Installation

```bash
helm plugin install https://github.com/amankr1098/helm-health
```

> **Requirements:** Helm 3+ and a valid kubeconfig pointing to the target cluster.

### Upgrade

```bash
helm plugin update health
```

### Uninstall

```bash
helm plugin uninstall health
```

---

## Usage

```bash
# Basic health check
helm health -r <release-name> -n <namespace>

# JSON output (great for CI/CD)
helm health -r <release-name> -n <namespace> -o json
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--release_name` | `-r` | *(required)* | Name of the Helm release to check |
| `--namespace` | `-n` | `default` | Kubernetes namespace of the release |
| `--output` | `-o` | `text` | Output format: `text` or `json` |

> **Tip:** When invoked through `helm health`, the `-n` flag from Helm is automatically forwarded via the `HELM_NAMESPACE` environment variable.

### JSON Output

```json
{
  "release": "my-app",
  "namespace": "production",
  "status": "Healthy",
  "timestamp": "2026-03-12T10:30:00Z",
  "duration": "1.2s",
  "summary": {
    "total": 6,
    "healthy": 5,
    "unhealthy": 1,
    "unknown": 0
  },
  "resources": [
    {
      "kind": "Deployment",
      "name": "my-app-web",
      "namespace": "production",
      "status": "Healthy",
      "health": { "readyReplicas": 3, "desiredReplicas": 3 }
    }
  ]
}
```

---

## Health Rules

Each Kubernetes resource kind has tailored health rules — for example, replica readiness for workloads, binding phase for PVCs, endpoint availability for Services, and more. Resources without specific checks are reported as **Healthy** if they exist in the release.

> See the full decision matrix in [docs/HEALTH_RULES.md](docs/HEALTH_RULES.md).

---

## Use in CI/CD

`helm-health` exits with a **non-zero exit code** when any resource is unhealthy, making it easy to gate deployments:

```yaml
# GitHub Actions example
- name: Check release health
  run: helm health -r my-app -n production -o json
```

```bash
# Shell script with retry
for i in $(seq 1 10); do
  helm health -r my-app -n production && break
  echo "Attempt $i: release not healthy yet, retrying in 10s..."
  sleep 10
done
```

---

## Development

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Helm 4.x](https://helm.sh/docs/intro/install/)
- Access to a Kubernetes cluster (kubeconfig)

### Build from source

```bash
git clone https://github.com/amankr1098/helm-health.git
cd helm-health
go build -o helm-health main.go
```

### Install locally

```bash
helm plugin install .
```

### Project structure

```
helm-health/
├── main.go                  # Entry point
├── cmd/
│   └── root.go              # CLI definition (cobra)
├── internal/
│   ├── data/                # Resource kind constants & priority lists
│   ├── output/              # Result types, text & JSON rendering
│   ├── release/             # Helm release fetching & manifest parsing
│   └── resources/           # Per-resource health check logic
│       ├── deployment.go
│       ├── statefulset.go
│       ├── pod.go
│       ├── job.go
│       ├── services.go
│       ├── pvc.go
│       ├── ingress.go
│       ├── networkpolicy.go
│       ├── demonset.go
│       └── clientset.go
└── docs/
    └── HEALTH_RULES.md      # Full health rule decision matrix
```

---

## Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository.
2. **Create a branch** for your feature or fix: `git checkout -b feat/my-feature`.
3. **Make your changes** and ensure they build: `go build ./...`.
4. **Submit a pull request** with a clear description of the change.

### Ideas for contributions

- Add `--watch` mode for continuous health monitoring
- Add configurable thresholds (e.g., restart count limits)
- Add health checks for more resource kinds
- Helm test integration
- Support for custom resource health plugins

---

## License

This project is open source. See the [LICENSE](LICENSE) file for details.

---

<div align="center">

**If this plugin saved you debugging time, consider giving it a ⭐!**

</div>