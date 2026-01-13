# kubectl-migrate

A kubectl plugin for migrating Kubernetes workloads and their state between clusters. This plugin integrates all features from the [crane migration tool](https://github.com/migtools/crane) and provides them through the familiar `kubectl migrate` command interface.

## Overview

kubectl-migrate is designed to help users migrate workloads between Kubernetes clusters in a safe, transparent, and composable way. It follows a pipeline-based approach to migration:

1. **Export** - Discover and export resources from source cluster
2. **Transform** - Apply transformations to exported manifests
3. **Apply** - Deploy transformed resources to target cluster

## Features

- 📦 **Export Resources** - Discover and export all resources from specified namespaces
- 🔄 **Transform Manifests** - Generate and apply JSONPatch transformations
- 🚀 **Apply Resources** - Deploy redeployable YAML to target clusters
- 💾 **PVC Transfer** - Migrate PersistentVolumeClaims between clusters
- 🔌 **Plugin System** - Extend functionality with custom plugins
- 🖼️ **Image Sync** - Generate Skopeo sync configurations for container images
- 🌐 **API Tunnel** - Tunnel API requests for migration scenarios
- 🔧 **Convert Resources** - Convert between resource formats

## Installation

### Via Krew (Recommended)

Once published to krew index:

```bash
kubectl krew install migrate
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/konveyor/kubectl-migrate.git
cd kubectl-migrate
```

2. Build and install:
```bash
make install
```

This will build the binary and install it to `$GOPATH/bin/kubectl-migrate`.

### From Release

Download the appropriate binary for your platform from the [releases page](https://github.com/konveyor/kubectl-migrate/releases) and place it in your `$PATH`.

## Usage

All commands are accessed via `kubectl migrate` followed by the specific subcommand.

### Basic Migration Workflow

```bash
# 1. Export resources from source cluster namespace
kubectl migrate export --namespace myapp --export-dir ./export

# 2. Transform the exported resources (optional)
kubectl migrate transform --export-dir ./export --transform-dir ./transform

# 3. Apply to target cluster
kubectl migrate apply --export-dir ./export --namespace myapp-migrated
```

## Available Commands

### Export

Export discovers and exports all resources from a specified namespace.

```bash
kubectl migrate export [namespace] [flags]

# Examples:
kubectl migrate export myapp --export-dir ./myapp-export
kubectl migrate export myapp --kubeconfig ./source-kubeconfig
```

**Key Flags:**
- `--export-dir` - Directory to export resources to
- `--kubeconfig` - Path to kubeconfig for source cluster
- `--context` - Context to use from kubeconfig

### Transform

Generate and apply JSONPatch transformations to exported resources.

```bash
kubectl migrate transform [flags]

# Examples:
kubectl migrate transform --export-dir ./export --transform-dir ./transform
kubectl migrate transform --plugin-dir ./plugins
```

**Key Flags:**
- `--export-dir` - Directory containing exported resources
- `--transform-dir` - Directory to write transformed resources
- `--plugin-dir` - Directory containing transform plugins

### Apply

Apply transformed resources to target cluster.

```bash
kubectl migrate apply [flags]

# Examples:
kubectl migrate apply --export-dir ./export --namespace target-ns
kubectl migrate apply --export-dir ./export --skip-namespaced
```

**Key Flags:**
- `--export-dir` - Directory containing resources to apply
- `--namespace` - Target namespace
- `--skip-namespaced` - Skip namespaced resources
- `--kubeconfig` - Path to kubeconfig for target cluster

### Transfer PVC

Transfer PersistentVolumeClaims between clusters.

```bash
kubectl migrate transfer-pvc [flags]

# Examples:
kubectl migrate transfer-pvc --source-context source --dest-context dest \
  --pvc-name my-pvc --pvc-namespace myapp
```

**Key Flags:**
- `--source-context` - Source cluster context
- `--dest-context` - Destination cluster context
- `--pvc-name` - Name of PVC to transfer
- `--pvc-namespace` - Namespace of PVC

### Plugin Manager

Manage kubectl-migrate plugins.

```bash
kubectl migrate plugin-manager [subcommand]

# Subcommands:
kubectl migrate plugin-manager list              # List installed plugins
kubectl migrate plugin-manager add <path>        # Add a plugin
kubectl migrate plugin-manager remove <name>     # Remove a plugin
```

### Skopeo Sync Gen

Generate Skopeo sync configuration for container images.

```bash
kubectl migrate skopeo-sync-gen [flags]

# Example:
kubectl migrate skopeo-sync-gen --export-dir ./export --output skopeo-sync.yaml
```

### Convert

Convert resources between different formats.

```bash
kubectl migrate convert [flags]

# Example:
kubectl migrate convert --input-dir ./export --output-dir ./converted
```

### Tunnel API

Tunnel API requests for specific migration scenarios.

```bash
kubectl migrate tunnel-api [flags]
```

### Version

Display version information.

```bash
kubectl migrate version
```

## Command Mapping from Crane

All `crane` commands are now available as `kubectl migrate` commands:

| Crane Command | kubectl-migrate Command |
|---------------|-------------------------|
| `crane export` | `kubectl migrate export` |
| `crane transform` | `kubectl migrate transform` |
| `crane apply` | `kubectl migrate apply` |
| `crane transfer-pvc` | `kubectl migrate transfer-pvc` |
| `crane plugin-manager` | `kubectl migrate plugin-manager` |
| `crane skopeo-sync-gen` | `kubectl migrate skopeo-sync-gen` |
| `crane convert` | `kubectl migrate convert` |
| `crane tunnel-api` | `kubectl migrate tunnel-api` |
| `crane version` | `kubectl migrate version` |

## Configuration

kubectl-migrate uses standard Kubernetes configuration:

- Kubeconfig files for cluster authentication
- Context switching for multi-cluster operations
- Standard kubectl flags like `--namespace`, `--context`, etc.

## Examples

### Migrate an application to a new cluster

```bash
# Set up contexts for source and target clusters
export SOURCE_CONTEXT=prod-cluster
export TARGET_CONTEXT=staging-cluster

# Export from source
kubectl migrate export myapp \
  --context $SOURCE_CONTEXT \
  --export-dir ./myapp-export

# Apply to target (with optional namespace change)
kubectl migrate apply \
  --context $TARGET_CONTEXT \
  --export-dir ./myapp-export \
  --namespace myapp-staging
```

### Migrate with PVC transfer

```bash
# Export application
kubectl migrate export myapp --export-dir ./myapp-export

# Transfer PVCs
kubectl migrate transfer-pvc \
  --source-context prod-cluster \
  --dest-context staging-cluster \
  --pvc-name myapp-data \
  --pvc-namespace myapp

# Apply to target
kubectl migrate apply \
  --context staging-cluster \
  --export-dir ./myapp-export
```

### Use plugins for transformation

```bash
# List available plugins
kubectl migrate plugin-manager list

# Export and transform with plugins
kubectl migrate export myapp --export-dir ./export
kubectl migrate transform \
  --export-dir ./export \
  --transform-dir ./transformed \
  --plugin-dir ./my-plugins

# Apply transformed resources
kubectl migrate apply --export-dir ./transformed
```

## Development

### Prerequisites

- Go 1.24 or later
- Access to Kubernetes clusters for testing
- kubectl installed

### Building from Source

```bash
# Clone repository
git clone https://github.com/konveyor/kubectl-migrate.git
cd kubectl-migrate

# Download dependencies
make deps

# Build
make build

# Run tests
make test

# Install locally
make install
```

### Running Tests

```bash
make test
```

### Cross-Platform Builds

```bash
make build-all
```

This creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Contributing

Contributions are welcome! Please see the [Konveyor Community](https://github.com/konveyor/community) for guidelines.

## Code of Conduct

This project follows the Konveyor [Code of Conduct](https://github.com/konveyor/community/blob/main/CODE_OF_CONDUCT.md).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [crane](https://github.com/migtools/crane) - Original migration tool
- [crane-lib](https://github.com/konveyor/crane-lib) - Shared library for crane tools
- [Konveyor](https://www.konveyor.io/) - Application modernization and migration toolkit

## Support

For questions and support:
- Open an issue on [GitHub](https://github.com/konveyor/kubectl-migrate/issues)
- Join the [Konveyor community](https://github.com/konveyor/community)

## Status

This project is currently in active development. APIs and features may change.
