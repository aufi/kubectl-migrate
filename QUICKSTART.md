# Quick Start Guide

This guide will help you get started with kubectl-migrate in minutes.

## Prerequisites

- kubectl installed and configured
- Access to source and destination Kubernetes clusters
- Appropriate RBAC permissions on both clusters

## Installation

### Option 1: Using Krew (Recommended when published)

```bash
kubectl krew install migrate
```

### Option 2: Manual Installation

```bash
# Clone the repository
git clone https://github.com/konveyor/kubectl-migrate.git
cd kubectl-migrate

# Build and install
make install
```

### Option 3: Download Pre-built Binary

Download from the [releases page](https://github.com/konveyor/kubectl-migrate/releases) and add to your PATH.

## Verify Installation

```bash
kubectl migrate version
```

## Basic Usage Examples

### Example 1: Simple Namespace Migration

Migrate a namespace from one cluster to another:

```bash
# Step 1: Set your source cluster context
kubectl config use-context source-cluster

# Step 2: Export the namespace
kubectl migrate export myapp --export-dir ./myapp-backup

# Step 3: Switch to destination cluster
kubectl config use-context dest-cluster

# Step 4: Apply to destination
kubectl migrate apply --export-dir ./myapp-backup --namespace myapp
```

### Example 2: Migration with Namespace Rename

```bash
# Export from source
kubectl migrate export production-app \
  --context source-cluster \
  --export-dir ./export

# Apply with new namespace name
kubectl migrate apply \
  --context dest-cluster \
  --export-dir ./export \
  --namespace staging-app
```

### Example 3: Selective Resource Export

```bash
# Export only specific resource types
kubectl migrate export myapp \
  --export-dir ./export \
  --included-resources deployments,services,configmaps
```

### Example 4: Transfer PersistentVolumeClaims

```bash
# Transfer a PVC between clusters
kubectl migrate transfer-pvc \
  --source-context prod-cluster \
  --dest-context staging-cluster \
  --pvc-name myapp-data \
  --pvc-namespace myapp
```

### Example 5: Using Transformations

```bash
# Export resources
kubectl migrate export myapp --export-dir ./export

# Transform resources (e.g., change image registry)
kubectl migrate transform \
  --export-dir ./export \
  --transform-dir ./transformed \
  --plugin-dir ./my-transforms

# Apply transformed resources
kubectl migrate apply \
  --export-dir ./transformed \
  --namespace myapp
```

## Common Workflows

### Disaster Recovery

Back up a namespace for disaster recovery:

```bash
# Regular backup
kubectl migrate export critical-app \
  --export-dir ./backups/$(date +%Y%m%d) \
  --context production

# Restore when needed
kubectl migrate apply \
  --export-dir ./backups/20260113 \
  --namespace critical-app \
  --context production
```

### Environment Promotion

Promote from dev to staging:

```bash
# Export from dev
kubectl migrate export myapp \
  --context dev-cluster \
  --export-dir ./promotion

# Transform for staging (if needed)
kubectl migrate transform \
  --export-dir ./promotion \
  --transform-dir ./staging-ready

# Apply to staging
kubectl migrate apply \
  --export-dir ./staging-ready \
  --context staging-cluster \
  --namespace myapp
```

### Cluster Migration

Migrate entire workloads between clusters:

```bash
# List of namespaces to migrate
NAMESPACES="app1 app2 app3"

for ns in $NAMESPACES; do
  echo "Migrating $ns..."

  # Export from source
  kubectl migrate export $ns \
    --context old-cluster \
    --export-dir ./migration/$ns

  # Apply to destination
  kubectl migrate apply \
    --context new-cluster \
    --export-dir ./migration/$ns \
    --namespace $ns
done
```

## Tips and Best Practices

1. **Always test first**: Try migrations in a test environment before production

2. **Use version control**: Store exported manifests in Git for tracking

3. **Review transformations**: Inspect transformed resources before applying

4. **Incremental migration**: Migrate non-critical apps first to validate the process

5. **Check dependencies**: Ensure all dependencies (storage classes, secrets, etc.) exist in destination

6. **Validate after migration**: Use `kubectl get all -n <namespace>` to verify resources

7. **Plan for downtime**: Some migrations may require brief service interruptions

## Troubleshooting

### Export fails with permission errors

```bash
# Check your RBAC permissions
kubectl auth can-i get deployments --namespace myapp
```

### Resources not created in destination

Check if resource definitions are compatible with the destination cluster version.

### PVC transfer fails

Ensure storage backends are compatible and accessible from both clusters.

## Next Steps

- Read the full [README](README.md) for detailed documentation
- Check available commands: `kubectl migrate --help`
- Explore plugin system for custom transformations
- Join the [Konveyor community](https://github.com/konveyor/community)

## Getting Help

- [GitHub Issues](https://github.com/konveyor/kubectl-migrate/issues)
- [Konveyor Community](https://github.com/konveyor/community)
- Command help: `kubectl migrate <command> --help`
