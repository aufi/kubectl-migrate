# E2E Tests for kubectl-migrate

This directory contains end-to-end (E2E) tests for the `kubectl-migrate` CLI tool. These tests validate the functionality of CLI commands against a real Kubernetes cluster.

## Overview

The E2E test framework provides:

- **CLI Execution**: Test CLI commands and validate their output
- **Cluster Interaction**: Deploy resources and verify cluster state
- **Assertions**: Rich assertion helpers for validating results
- **Test Isolation**: Each test runs in its own namespace with automatic cleanup

## Directory Structure

```
test/e2e/
├── framework/           # Test framework components
│   ├── cli.go          # CLI command execution
│   ├── cluster.go      # Kubernetes cluster management
│   ├── assertions.go   # Assertion helpers
│   └── suite.go        # Test suite helpers
├── export_test.go      # Export command tests
└── README.md           # This file
```

## Prerequisites

1. **Go 1.21+** installed
2. **kubectl** installed and configured
3. **Running Kubernetes cluster** (Kind, Minikube, or any other cluster)
4. **Built binary** of kubectl-migrate

## Running Tests

### Build the Binary

```bash
make build
```

### Run All E2E Tests

```bash
make test-e2e
```

### Run Specific Test Suite

```bash
# Export command tests only
make test-e2e-export

# Quick tests (short mode)
make test-e2e-quick
```

### Run Individual Tests

```bash
# Run specific test function
E2E_BINARY=./bin/kubectl-migrate go test -v ./test/e2e/export_test.go ./test/e2e/framework/*.go -run TestExportCommand/basic

# Run with verbose output
E2E_BINARY=./bin/kubectl-migrate go test -v ./test/e2e/... -test.v

# Run in short mode (skips long-running tests)
E2E_BINARY=./bin/kubectl-migrate go test -v ./test/e2e/... -short
```

## Writing Tests

### Basic Test Structure

```go
package e2e

import (
    "testing"
    "github.com/konveyor-ecosystem/kubectl-migrate/test/e2e/framework"
)

func TestMyCommand(t *testing.T) {
    // Create test suite
    suite := framework.NewTestSuite(t)
    defer suite.Cleanup()

    // Skip if running in short mode
    suite.SkipIfShort("This test requires a cluster")

    // Create test namespace
    ns := suite.CreateTestNamespace("my-test")

    // Run CLI command
    result := suite.RunCLI("export", "--namespace", ns, "--export-dir", suite.TempDir)

    // Assert results
    suite.Assert.AssertCommandSuccess(result)
    suite.Assert.AssertDirExists(suite.TempDir)
}
```

### Test Framework Components

#### CLI Executor

Execute CLI commands and capture output:

```go
// Basic execution
result := suite.RunCLI("export", "--namespace", "default")

// Check results
if result.Success() {
    fmt.Println("Command succeeded")
}

// Access output
fmt.Println(result.Stdout)
fmt.Println(result.Stderr)
fmt.Println(result.ExitCode)
```

#### Cluster Manager

Interact with Kubernetes cluster:

```go
// Wait for deployment
err := suite.Cluster.WaitForDeployment("default", "my-app", 120*time.Second)

// Get pods
pods, err := suite.Cluster.GetPods("default", map[string]string{"app": "nginx"})

// Check resource exists
exists, err := suite.Cluster.ResourceExists(gvr, "default", "my-resource")
```

#### Assertions

Validate test results:

```go
// Command assertions
suite.Assert.AssertCommandSuccess(result)
suite.Assert.AssertOutputContains(result, "expected text")

// File assertions
suite.Assert.AssertFileExists("/path/to/file")
suite.Assert.AssertYAMLFileValid("/path/to/file.yaml")

// Cluster assertions
suite.Assert.AssertResourceExists(cluster, gvr, "namespace", "name")
suite.Assert.AssertResourceCount(cluster, gvr, "namespace", "app=nginx", 3)
```

### Test Patterns

#### Testing Export Command

```go
func TestExportBasic(t *testing.T) {
    suite := framework.NewTestSuite(t)
    defer suite.Cleanup()

    // Create namespace and deploy app
    ns := suite.CreateTestNamespace("export-test")
    deployTestApp(t, suite, ns)

    // Wait for deployment
    err := suite.Cluster.WaitForDeployment(ns, "test-app", 120*time.Second)
    suite.Assert.AssertNoError(err)

    // Export resources
    exportDir := suite.CreateTempDir("export")
    result := suite.RunCLI("export", "--namespace", ns, "--export-dir", exportDir)

    // Validate
    suite.Assert.AssertCommandSuccess(result)
    suite.Assert.AssertMinYAMLFiles(exportDir, 1)
}
```

#### Testing with Label Selectors

```go
result := suite.RunCLI(
    "export",
    "--namespace", ns,
    "--label-selector", "tier=frontend",
    "--export-dir", exportDir,
)
```

#### Testing Error Cases

```go
result := suite.RunCLI("export", "--namespace", "nonexistent")
suite.Assert.AssertCommandFails(result)
suite.Assert.AssertStderrContains(result, "not found")
```

## Test Organization

### Test Naming Convention

- `TestCommandName` - Main test function
- `testCommandScenario` - Helper test functions
- Use descriptive names that explain what is being tested

### Test Structure

Each test file should:

1. Import the framework package
2. Create a test suite in each test
3. Register cleanup with `defer suite.Cleanup()`
4. Use assertion helpers instead of raw `if` checks
5. Log important steps for debugging

### Cleanup

Always clean up resources:

```go
func TestMyTest(t *testing.T) {
    suite := framework.NewTestSuite(t)
    defer suite.Cleanup()  // Automatic cleanup

    // Test creates namespace automatically cleaned up
    ns := suite.CreateTestNamespace("test")

    // Register custom cleanup if needed
    suite.RegisterCleanup(func() {
        // Custom cleanup logic
    })
}
```

## CI/CD Integration

### GitHub Actions

E2E tests run automatically on:

- Pull requests that modify code
- Pushes to main branch
- Manual workflow dispatch

See `.github/workflows/e2e-tests.yml` for configuration.

### Local Testing with Kind

```bash
# Create Kind cluster
kind create cluster --name e2e-test

# Run tests
make test-e2e

# Cleanup
kind delete cluster --name e2e-test
```

## Environment Variables

- `E2E_BINARY` - Path to kubectl-migrate binary (default: `./bin/kubectl-migrate`)
- `KUBECONFIG` - Path to kubeconfig file (uses default if not set)

## Troubleshooting

### Tests Fail with "cluster not found"

Ensure you have a running Kubernetes cluster:

```bash
kubectl cluster-info
```

### Tests Timeout

Increase timeout for slow environments:

```bash
go test -v ./test/e2e/... -timeout 60m
```

### Binary Not Found

Ensure the binary is built:

```bash
make build
ls -la ./bin/kubectl-migrate
```

### View Verbose Test Output

```bash
go test -v ./test/e2e/... -test.v
```

### Debug Specific Test

```bash
# Run single test with verbose output
E2E_BINARY=./bin/kubectl-migrate go test -v ./test/e2e/export_test.go ./test/e2e/framework/*.go -run TestExportCommand/basic -test.v
```

## Best Practices

1. **Isolate Tests**: Each test should create its own namespace
2. **Clean Up**: Always use `defer suite.Cleanup()`
3. **Short Mode**: Mark long tests to skip in short mode with `suite.SkipIfShort()`
4. **Descriptive Names**: Use clear test and assertion names
5. **Log Progress**: Use `suite.LogInfo()` for important steps
6. **Check Errors**: Use assertion helpers instead of manual error checking
7. **Wait for Ready**: Always wait for resources to be ready before testing

## Contributing

When adding new E2E tests:

1. Follow existing patterns in `export_test.go`
2. Use the framework helpers instead of raw Kubernetes API calls
3. Add test documentation in comments
4. Ensure tests pass locally before submitting PR
5. Update this README if adding new patterns or features

## Future Enhancements

Planned improvements:

- [ ] Tests for `convert` command
- [ ] Tests for `apply` command
- [ ] Tests for `transform` command
- [ ] Tests for `transfer-pvc` command
- [ ] Multi-cluster tests
- [ ] Performance benchmarks
- [ ] Test fixtures for common scenarios

## Support

For issues or questions:

- Open an issue in the GitHub repository
- Check existing tests for examples
- Review framework code in `test/e2e/framework/`
