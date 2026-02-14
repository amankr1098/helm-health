# Helm Health Plugin

A Helm plugin for health checking deployed applications.

## Installation

```bash
helm plugin install https://github.com/your-username/helm-health
```

## Usage

```bash
helm health [command]
```

## Development

### Prerequisites

- Go 1.21 or higher
- Helm 3.x

### Building

```bash
go build -o helm-health main.go
```

### Testing locally

```bash
helm plugin install .
helm health
```