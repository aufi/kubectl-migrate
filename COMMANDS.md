# kubectl-migrate Commands Reference

Complete reference for all kubectl-migrate commands. All commands can be invoked with the `kubectl migrate` prefix.

## Table of Contents

- [export](#export) - Export namespace resources
- [transform](#transform) - Transform exported resources
- [apply](#apply) - Apply resources to cluster
- [transfer-pvc](#transfer-pvc) - Transfer PersistentVolumeClaims
- [plugin-manager](#plugin-manager) - Manage plugins
- [skopeo-sync-gen](#skopeo-sync-gen) - Generate Skopeo sync config
- [convert](#convert) - Convert deprecated resources
- [tunnel-api](#tunnel-api) - Set up API tunnel
- [runfn](#runfn) - Run KRM functions
- [version](#version) - Show version information

---

## export

Export discovers and exports all resources from a specified namespace.

### Usage

```bash
kubectl migrate export [namespace] [flags]
```

### Examples

```bash
# Export a namespace to a directory
kubectl migrate export myapp --export-dir ./export

# Export with specific kubeconfig
kubectl migrate export myapp --kubeconfig ~/.kube/prod-config

# Export from specific context
kubectl migrate export myapp --context prod-cluster

# Export including cluster-scoped RBAC
kubectl migrate export myapp --cluster-scoped-rbac

# Export only specific resource types
kubectl migrate export myapp --included-resources deployments,services
```

### Key Flags

- `--export-dir, -e` - Directory to export resources (default: "export")
- `--kubeconfig` - Path to kubeconfig file
- `--context` - Kubeconfig context to use
- `--namespace, -n` - Namespace to export
- `--cluster-scoped-rbac, -c` - Include cluster-scoped RBAC resources
- `--included-resources` - Comma-separated list of resource types to include
- `--excluded-resources` - Comma-separated list of resource types to exclude

---

## transform

Generate and apply JSONPatch transformations to exported resources.

### Usage

```bash
kubectl migrate transform [flags]
```

### Examples

```bash
# Transform exported resources
kubectl migrate transform --export-dir ./export --transform-dir ./transformed

# Use specific plugins
kubectl migrate transform --export-dir ./export --plugin-dir ./my-plugins

# List available transform plugins
kubectl migrate transform listplugins

# Apply optional transformations
kubectl migrate transform optionals --export-dir ./export
```

### Key Flags

- `--export-dir, -e` - Directory containing exported resources
- `--transform-dir, -t` - Directory to write transformed resources
- `--plugin-dir, -p` - Directory containing transform plugins
- `--flags-file, -f` - YAML file with transformation flags

### Subcommands

- `listplugins` - List available transformation plugins
- `optionals` - Apply optional transformations

---

## apply

Apply transformed resources to a target cluster.

### Usage

```bash
kubectl migrate apply [flags]
```

### Examples

```bash
# Apply resources from export directory
kubectl migrate apply --export-dir ./export

# Apply to specific namespace
kubectl migrate apply --export-dir ./export --namespace target-ns

# Apply only cluster-scoped resources
kubectl migrate apply --export-dir ./export --skip-namespaced

# Apply with dry-run
kubectl migrate apply --export-dir ./export --dry-run
```

### Key Flags

- `--export-dir, -e` - Directory containing resources to apply
- `--namespace, -n` - Target namespace
- `--kubeconfig` - Path to kubeconfig file
- `--context` - Kubeconfig context to use
- `--skip-namespaced` - Skip namespaced resources
- `--skip-cluster-scoped` - Skip cluster-scoped resources
- `--dry-run` - Perform dry run without actually applying

---

## transfer-pvc

Transfer PersistentVolumeClaims between clusters.

### Usage

```bash
kubectl migrate transfer-pvc [flags]
```

### Examples

```bash
# Transfer a PVC between clusters
kubectl migrate transfer-pvc \
  --source-context prod-cluster \
  --dest-context staging-cluster \
  --pvc-name myapp-data \
  --pvc-namespace myapp

# Transfer with custom endpoint
kubectl migrate transfer-pvc \
  --source-context source \
  --dest-context dest \
  --pvc-name data \
  --pvc-namespace default \
  --endpoint https://custom-endpoint
```

### Key Flags

- `--source-context` - Source cluster context
- `--dest-context` - Destination cluster context
- `--pvc-name` - Name of PVC to transfer
- `--pvc-namespace` - Namespace of PVC
- `--endpoint` - Custom transfer endpoint
- `--source-path` - Path in source PVC
- `--dest-path` - Path in destination PVC

---

## plugin-manager

Manage kubectl-migrate plugins for custom transformations.

### Usage

```bash
kubectl migrate plugin-manager [subcommand]
```

### Subcommands

#### list

List all installed plugins.

```bash
kubectl migrate plugin-manager list
```

#### add

Add a new plugin.

```bash
kubectl migrate plugin-manager add <path-to-plugin>
kubectl migrate plugin-manager add ./my-transform-plugin
```

#### remove

Remove an installed plugin.

```bash
kubectl migrate plugin-manager remove <plugin-name>
kubectl migrate plugin-manager remove my-transform
```

---

## skopeo-sync-gen

Generate Skopeo sync configuration for container images found in exported resources.

### Usage

```bash
kubectl migrate skopeo-sync-gen [flags]
```

### Examples

```bash
# Generate sync config from exported resources
kubectl migrate skopeo-sync-gen --export-dir ./export --output skopeo-sync.yaml

# Generate with custom registries
kubectl migrate skopeo-sync-gen \
  --export-dir ./export \
  --source-registry old-registry.io \
  --dest-registry new-registry.io
```

### Key Flags

- `--export-dir, -e` - Directory containing exported resources
- `--output, -o` - Output file for sync configuration
- `--source-registry` - Source container registry
- `--dest-registry` - Destination container registry

---

## convert

Convert deprecated resources to their current API versions.

### Usage

```bash
kubectl migrate convert [flags]
```

### Examples

```bash
# Convert resources in a directory
kubectl migrate convert --input-dir ./export --output-dir ./converted

# Convert specific files
kubectl migrate convert --filename deployment.yaml --output converted.yaml
```

### Key Flags

- `--input-dir` - Directory containing resources to convert
- `--output-dir` - Directory for converted resources
- `--filename, -f` - Specific file to convert
- `--output, -o` - Output format (yaml, json)

---

## tunnel-api

Set up an OpenVPN tunnel to access a source cluster from a destination cluster.

### Usage

```bash
kubectl migrate tunnel-api [flags]
```

### Examples

```bash
# Set up tunnel between clusters
kubectl migrate tunnel-api \
  --source-context onprem-cluster \
  --dest-context cloud-cluster
```

### Key Flags

- `--source-context` - Source cluster context
- `--dest-context` - Destination cluster context
- `--tunnel-config` - Path to tunnel configuration

---

## runfn

Transform resources by executing Kustomize Resource Model (KRM) functions.

### Usage

```bash
kubectl migrate runfn [flags]
```

### Examples

```bash
# Run KRM function on resources
kubectl migrate runfn --export-dir ./export --function-path ./my-function.yaml

# Run with specific function image
kubectl migrate runfn \
  --export-dir ./export \
  --image gcr.io/kpt-fn/set-namespace:v0.4
```

### Key Flags

- `--export-dir, -e` - Directory containing resources
- `--function-path` - Path to KRM function configuration
- `--image` - Container image of KRM function
- `--network` - Enable network access for function

---

## version

Display version information for kubectl-migrate and crane-lib.

### Usage

```bash
kubectl migrate version
```

### Example Output

```
crane:
	Version: v0.1.0
crane-lib:
	Version: v0.0.10
```

---

## Global Flags

These flags are available for all commands:

- `--debug` - Enable debug logging
- `--flags-file, -f` - Path to YAML file containing flag values
- `--help, -h` - Show help for command

---

## Environment Variables

kubectl-migrate respects standard kubectl environment variables:

- `KUBECONFIG` - Path to kubeconfig file
- `KUBECTL_CONTEXT` - Default context to use

---

## Command Mapping from Crane

| Crane Command | kubectl-migrate Command | Description |
|---------------|-------------------------|-------------|
| `crane export` | `kubectl migrate export` | Export namespace resources |
| `crane transform` | `kubectl migrate transform` | Transform exported resources |
| `crane apply` | `kubectl migrate apply` | Apply resources to cluster |
| `crane transfer-pvc` | `kubectl migrate transfer-pvc` | Transfer PVCs |
| `crane plugin-manager` | `kubectl migrate plugin-manager` | Manage plugins |
| `crane skopeo-sync-gen` | `kubectl migrate skopeo-sync-gen` | Generate Skopeo config |
| `crane convert` | `kubectl migrate convert` | Convert resources |
| `crane tunnel-api` | `kubectl migrate tunnel-api` | Set up API tunnel |
| `crane runfn` | `kubectl migrate runfn` | Run KRM functions |
| `crane version` | `kubectl migrate version` | Show version |

---

## Getting Help

For detailed help on any command:

```bash
kubectl migrate <command> --help
```

For general help:

```bash
kubectl migrate --help
```

For issues or questions:
- [GitHub Issues](https://github.com/konveyor/kubectl-migrate/issues)
- [Konveyor Community](https://github.com/konveyor/community)
