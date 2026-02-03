# Contributing to E2E Tests

Thank you for contributing to kubectl-migrate E2E tests!

## Quick Start

1. **Setup your environment**
   ```bash
   # Ensure you have a running Kubernetes cluster
   kubectl cluster-info

   # Build the binary
   make build
   ```

2. **Run existing tests**
   ```bash
   # Run all E2E tests
   make test-e2e

   # Run only export tests
   make test-e2e-export
   ```

3. **Write a new test**
   - Add your test to an existing `*_test.go` file or create a new one
   - Follow the patterns shown in `export_test.go`
   - Use the test framework from `test/e2e/framework/`

## Test Framework Usage

### Using Sample Resources

**Always use existing sample resources** from `sample-resources/` instead of creating new test applications:

```go
// Good - uses existing sample-resources/hello-world
func deployTestApp(t *testing.T, suite *framework.TestSuite, namespace string) {
    cmd := exec.Command("kubectl", "apply", "-f", "sample-resources/hello-world/manifest.yaml", "-n", namespace)
    output, err := cmd.CombinedOutput()
    // ... handle deployment
}

// Bad - creates duplicate resources
func deployTestApp(...) {
    manifest := `apiVersion: apps/v1...` // Don't do this!
}
```

### Available Sample Resources

- `sample-resources/hello-world/` - Simple Apache deployment
  - Deployment: `apache-hello`
  - Service: `apache-hello`
  - Labels: `app=hello-world`

- `sample-resources/wordpress/` - WordPress with MySQL
  - Deployments: `wordpress`, `wordpress-mysql`
  - Services: `wordpress`, `wordpress-mysql`
  - Has validation scripts

### Test Structure Template

```go
func TestMyFeature(t *testing.T) {
    suite := framework.NewTestSuite(t)
    defer suite.Cleanup()

    suite.SkipIfShort("Reason for skipping in short mode")

    ns := suite.CreateTestNamespace("my-feature")

    // Deploy using sample resources
    deployTestApp(t, suite, ns)

    // Wait for ready
    err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
    suite.Assert.AssertNoError(err)

    // Run your test
    result := suite.RunCLI("your-command", "--namespace", ns)
    suite.Assert.AssertCommandSuccess(result)
}
```

## Testing Guidelines

### Do's

✅ Use `suite.CreateTestNamespace()` for test isolation
✅ Use `defer suite.Cleanup()` for automatic cleanup
✅ Use existing sample-resources applications
✅ Use assertion helpers (`suite.Assert.*`)
✅ Log important steps with `suite.LogInfo()`
✅ Wait for resources to be ready before testing
✅ Handle both success and failure cases

### Don'ts

❌ Don't create inline YAML manifests
❌ Don't forget cleanup
❌ Don't skip `suite.SkipIfShort()` for cluster tests
❌ Don't use raw error checking (use assertions)
❌ Don't hardcode namespaces
❌ Don't assume resources are immediately ready

## Running Tests Locally

### With Kind

```bash
# Create cluster
kind create cluster --name e2e-test

# Run tests
make test-e2e

# Cleanup
kind delete cluster --name e2e-test
```

### With Existing Cluster

```bash
# Just run tests
make test-e2e

# Run specific test
E2E_BINARY=./bin/kubectl-migrate go test -v ./test/e2e/export_test.go ./test/e2e/framework/*.go -run TestExportCommand/basic
```

## Adding a New Command Test

1. Create `test/e2e/mycommand_test.go`
2. Import the framework:
   ```go
   import "github.com/konveyor-ecosystem/kubectl-migrate/test/e2e/framework"
   ```
3. Write tests following export_test.go patterns
4. Add Makefile target:
   ```makefile
   test-e2e-mycommand: build
       @E2E_BINARY=$(BUILD_DIR)/$(BINARY_NAME) $(GOTEST) -v -timeout 15m ./test/e2e/mycommand_test.go ./test/e2e/framework/*.go
   ```
5. Update `.github/workflows/e2e-tests.yml` if needed

## Debugging Tests

### View Test Output

```bash
# Verbose output
go test -v ./test/e2e/... -test.v

# Run specific test
go test -v ./test/e2e/export_test.go ./test/e2e/framework/*.go -run TestExportHelp
```

### Check Resources After Test

```bash
# List all namespaces (look for *-e2e-* namespaces)
kubectl get ns

# Check resources in test namespace
kubectl get all -n export-basic-e2e-12345
```

### Common Issues

**"namespace not found"**
- Test might be running too fast
- Add `suite.Cluster.WaitForDeployment()` calls

**"deployment not ready"**
- Increase timeout in `WaitForDeployment()`
- Check cluster has enough resources

**"binary not found"**
- Run `make build` first
- Set `E2E_BINARY` environment variable

## Code Review Checklist

Before submitting a PR with E2E tests:

- [ ] Tests pass locally
- [ ] Uses existing sample-resources
- [ ] Includes proper cleanup
- [ ] Has meaningful assertions
- [ ] Includes test documentation
- [ ] Follows existing patterns
- [ ] Works in short mode (if applicable)
- [ ] Updated README if adding new patterns

## Questions?

- Check `test/e2e/README.md` for usage examples
- Review `export_test.go` for patterns
- Look at framework code in `test/e2e/framework/`
- Open an issue for questions
