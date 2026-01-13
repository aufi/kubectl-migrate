# Krew Plugin Deployment Guide

This guide explains how to publish kubectl-migrate to the krew plugin index.

## Prerequisites

- Plugin successfully built and tested
- GitHub releases created with platform-specific binaries
- SHA256 checksums generated for all release artifacts

## Steps to Publish to Krew

### 1. Create a GitHub Release

```bash
# Tag your release
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

This will trigger the GitHub Actions release workflow which:
- Builds binaries for all platforms
- Creates .tar.gz and .zip archives
- Generates SHA256 checksums
- Creates a GitHub release with all artifacts

### 2. Update the Krew Manifest

After the release is created, update `kubectl-migrate.yaml`:

```bash
# Download the checksums file from the GitHub release
wget https://github.com/konveyor/kubectl-migrate/releases/download/v0.1.0/checksums.txt

# Update the sha256 values in kubectl-migrate.yaml with the checksums
# For each platform (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64)
```

Example from checksums.txt:
```
abc123...  kubectl-migrate-linux-amd64.tar.gz
def456...  kubectl-migrate-darwin-amd64.tar.gz
```

Update in kubectl-migrate.yaml:
```yaml
- selector:
    matchLabels:
      os: linux
      arch: amd64
  uri: https://github.com/konveyor/kubectl-migrate/releases/download/v0.1.0/kubectl-migrate-linux-amd64.tar.gz
  sha256: "abc123..."
  bin: kubectl-migrate
```

### 3. Test the Plugin Locally

Before submitting to krew-index, test locally:

```bash
# Install krew if not already installed
kubectl krew install krew

# Test installing your plugin locally
kubectl krew install --manifest=kubectl-migrate.yaml
kubectl migrate version

# Uninstall after testing
kubectl krew uninstall migrate
```

### 4. Submit to Krew Index

Fork and clone the krew-index repository:

```bash
git clone https://github.com/kubernetes-sigs/krew-index.git
cd krew-index

# Create a new branch
git checkout -b add-kubectl-migrate

# Copy your plugin manifest
cp /path/to/kubectl-migrate.yaml plugins/migrate.yaml

# Commit and push
git add plugins/migrate.yaml
git commit -m "Add kubectl-migrate plugin"
git push origin add-kubectl-migrate
```

### 5. Create Pull Request

1. Go to https://github.com/kubernetes-sigs/krew-index
2. Create a pull request from your fork
3. Fill in the PR template with:
   - Plugin name: migrate
   - Short description
   - Link to plugin repository
   - Confirmation that you've tested the plugin

### 6. Wait for Review

The krew maintainers will review your PR. They will check:
- Manifest format is correct
- All platform binaries are available
- SHA256 checksums match
- Plugin installs and works correctly

## Updating the Plugin

For subsequent releases:

```bash
# 1. Create new release
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# 2. Wait for GitHub Actions to build and create release

# 3. Update kubectl-migrate.yaml with new version and checksums

# 4. Submit PR to krew-index
cd krew-index
git checkout main
git pull upstream main
git checkout -b update-kubectl-migrate-v0.2.0

# Update plugins/migrate.yaml with your new version
cp /path/to/kubectl-migrate.yaml plugins/migrate.yaml

git add plugins/migrate.yaml
git commit -m "Update kubectl-migrate to v0.2.0"
git push origin update-kubectl-migrate-v0.2.0

# Create PR
```

## Validation Checklist

Before submitting to krew, ensure:

- [ ] All platform binaries build successfully
- [ ] Binary names match the manifest (kubectl-migrate)
- [ ] Archives extract correctly
- [ ] SHA256 checksums are correct
- [ ] URIs point to actual release artifacts
- [ ] Version in manifest matches git tag
- [ ] Plugin installs via `kubectl krew install --manifest`
- [ ] All commands work after installation
- [ ] Uninstall works cleanly

## Testing Different Platforms

### Linux (amd64)
```bash
# Download and test
wget https://github.com/konveyor/kubectl-migrate/releases/download/v0.1.0/kubectl-migrate-linux-amd64.tar.gz
tar xzf kubectl-migrate-linux-amd64.tar.gz
./kubectl-migrate version
```

### macOS (arm64)
```bash
# Download and test
curl -LO https://github.com/konveyor/kubectl-migrate/releases/download/v0.1.0/kubectl-migrate-darwin-arm64.tar.gz
tar xzf kubectl-migrate-darwin-arm64.tar.gz
./kubectl-migrate version
```

### Windows (amd64)
```powershell
# Download and test
Invoke-WebRequest -Uri "https://github.com/konveyor/kubectl-migrate/releases/download/v0.1.0/kubectl-migrate-windows-amd64.zip" -OutFile "kubectl-migrate.zip"
Expand-Archive kubectl-migrate.zip
.\kubectl-migrate\kubectl-migrate.exe version
```

## Troubleshooting

### Plugin not installing
- Verify all URLs are accessible
- Check SHA256 checksums match exactly
- Ensure binary has execute permissions in the archive

### Wrong binary name
- Binary must be named `kubectl-migrate` (or `kubectl-migrate.exe` on Windows)
- Krew automatically creates the `kubectl-` prefix for plugin invocation

### Architecture not supported
- Ensure you've built for all major platforms
- Minimum recommended: linux/amd64, darwin/amd64, darwin/arm64

## Resources

- [Krew Developer Guide](https://krew.sigs.k8s.io/docs/developer-guide/)
- [Krew Plugin Naming Guide](https://krew.sigs.k8s.io/docs/developer-guide/plugin-naming/)
- [krew-index Repository](https://github.com/kubernetes-sigs/krew-index)
