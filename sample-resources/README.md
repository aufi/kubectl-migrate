# Sample Resources for Container Migration

This directory contains sample Kubernetes applications for testing and developing container migration functionality with kubectl-migrate.

## Available Applications

| Application | Type | Description |
|------------|------|-------------|
| hello-world | Stateless | Apache HTTPD web server with static "Hello World" page |

## Prerequisites

- Kubernetes cluster (local or remote)
- `kubectl` configured to communicate with your cluster
- Proper permissions to create/delete resources in your cluster

## Quick Start

### Deploy All Applications

To deploy all sample applications to your Kubernetes cluster:

```bash
make resources-deploy
```

This will execute all `deploy.sh` scripts in the subdirectories.

### Destroy All Applications

To remove all sample applications from your Kubernetes cluster:

```bash
make resources-destroy
```

This will execute all `destroy.sh` scripts in the subdirectories.

## Individual Application Management

Each application subdirectory contains:
- `manifest.yaml` - Kubernetes resource definitions
- `deploy.sh` - Script to deploy the application
- `destroy.sh` - Script to remove the application
- `README.md` - Application-specific documentation

To manage applications individually, navigate to the specific application directory and run:

```bash
cd <application-directory>
./deploy.sh    # Deploy the application
./destroy.sh   # Remove the application
```

## Adding New Sample Applications

When adding new sample applications to this directory:

1. Create a new subdirectory with a descriptive name
2. Include `manifest.yaml`, `deploy.sh`, `destroy.sh`, and `README.md`
3. Make scripts executable: `chmod +x deploy.sh destroy.sh`
4. Update the applications table in this README
5. The Makefile will automatically pick up new deploy/destroy scripts

## Usage with kubectl-migrate

These sample applications are designed to test various migration scenarios:

- Stateless applications
- Stateful applications with persistent storage
- Applications with complex networking requirements
- Multi-tier applications

Use these applications to verify migration functionality, test backup/restore operations, and validate cross-cluster migrations.
