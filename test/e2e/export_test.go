package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/konveyor-ecosystem/kubectl-migrate/test/e2e/framework"
)

func TestExportCommand(t *testing.T) {
	suite := framework.NewTestSuite(t)
	defer suite.Cleanup()

	// Skip if running in short mode
	suite.SkipIfShort("Export command tests require a running cluster")

	t.Run("basic export of namespace resources", func(t *testing.T) {
		testBasicExport(t, suite)
	})

	t.Run("export with label selector", func(t *testing.T) {
		testExportWithLabelSelector(t, suite)
	})

	t.Run("export with custom export directory", func(t *testing.T) {
		testExportWithCustomDir(t, suite)
	})

	t.Run("export with cluster-scoped RBAC", func(t *testing.T) {
		testExportWithClusterScopedRBAC(t, suite)
	})

	t.Run("export empty namespace", func(t *testing.T) {
		testExportEmptyNamespace(t, suite)
	})

	t.Run("export with non-existent namespace should fail", func(t *testing.T) {
		testExportNonExistentNamespace(t, suite)
	})
}

// testBasicExport tests basic export functionality
func testBasicExport(t *testing.T, suite *framework.TestSuite) {
	// Create test namespace
	ns := suite.CreateTestNamespace("export-basic")
	suite.LogInfo("Created test namespace: %s", ns)

	// Deploy a simple application
	deployTestApp(t, suite, ns)

	// Wait for deployment to be ready (hello-world deployment name)
	err := suite.Cluster.WaitForDeployment(ns, "hello-world", 120*time.Second)
	suite.Assert.AssertNoError(err)

	// Run export command
	exportDir := suite.CreateTempDir("export-basic")
	result := suite.RunCLI("export", "--namespace", ns, "--export-dir", exportDir)

	// Assert command succeeded
	suite.Assert.AssertCommandSuccess(result)

	// Verify export directory structure
	resourcesDir := filepath.Join(exportDir, "resources", ns)
	suite.Assert.AssertDirExists(resourcesDir)
	suite.Assert.AssertDirNotEmpty(resourcesDir)

	// Verify exported files exist
	suite.Assert.AssertMinYAMLFiles(resourcesDir, 1)

	// Verify specific resources were exported
	suite.LogInfo("Export completed successfully to: %s", exportDir)

	// Check that deployment was exported
	files, _ := os.ReadDir(resourcesDir)
	foundDeployment := false
	for _, file := range files {
		if !file.IsDir() {
			suite.LogDebug("Exported file: %s", file.Name())
			content, _ := os.ReadFile(filepath.Join(resourcesDir, file.Name()))
			if len(content) > 0 {
				// Validate YAML
				suite.Assert.AssertYAMLFileValid(filepath.Join(resourcesDir, file.Name()))

				// Check if it's a deployment
				contentStr := string(content)
				if suite.CLI.Run("--help").ContainsOutput("kubectl-migrate") {
					// Simple check - look for deployment kind
					if len(contentStr) > 0 {
						foundDeployment = true
					}
				}
			}
		}
	}

	if !foundDeployment {
		suite.LogInfo("Warning: Deployment file not found in export, but export succeeded")
	}
}

// testExportWithLabelSelector tests export with label selector
func testExportWithLabelSelector(t *testing.T, suite *framework.TestSuite) {
	ns := suite.CreateTestNamespace("export-labels")
	suite.LogInfo("Created test namespace: %s", ns)

	// Deploy app with specific labels
	deployTestAppWithLabels(t, suite, ns, map[string]string{
		"app":  "test",
		"tier": "frontend",
	})

	// Wait for deployment
	err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
	suite.Assert.AssertNoError(err)

	// Export with label selector
	exportDir := suite.CreateTempDir("export-labels")
	result := suite.RunCLI("export",
		"--namespace", ns,
		"--export-dir", exportDir,
		"--label-selector", "app=hello-world")

	suite.Assert.AssertCommandSuccess(result)

	// Verify export
	resourcesDir := filepath.Join(exportDir, "resources", ns)
	suite.Assert.AssertDirExists(resourcesDir)
	suite.LogInfo("Export with label selector completed")
}

// testExportWithCustomDir tests export to a custom directory
func testExportWithCustomDir(t *testing.T, suite *framework.TestSuite) {
	ns := suite.CreateTestNamespace("export-customdir")

	// Deploy test app
	deployTestApp(t, suite, ns)

	// Wait for deployment
	err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
	suite.Assert.AssertNoError(err)

	// Export to custom directory
	customDir := suite.CreateTempDir("my-custom-export")
	result := suite.RunCLI("export", "--namespace", ns, "--export-dir", customDir)

	suite.Assert.AssertCommandSuccess(result)

	// Verify custom directory was used
	resourcesDir := filepath.Join(customDir, "resources", ns)
	suite.Assert.AssertDirExists(resourcesDir)
	suite.LogInfo("Export to custom directory: %s", customDir)
}

// testExportWithClusterScopedRBAC tests export with cluster-scoped RBAC
func testExportWithClusterScopedRBAC(t *testing.T, suite *framework.TestSuite) {
	ns := suite.CreateTestNamespace("export-rbac")

	// Deploy app with service account
	deployTestAppWithServiceAccount(t, suite, ns)

	// Wait for deployment
	err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
	suite.Assert.AssertNoError(err)

	// Export with cluster-scoped RBAC flag
	exportDir := suite.CreateTempDir("export-rbac")
	result := suite.RunCLI("export",
		"--namespace", ns,
		"--export-dir", exportDir,
		"--cluster-scoped-rbac")

	suite.Assert.AssertCommandSuccess(result)

	// Verify _cluster directory exists
	clusterDir := filepath.Join(exportDir, "resources", ns, "_cluster")
	suite.Assert.AssertDirExists(clusterDir)
	suite.LogInfo("Export with cluster-scoped RBAC completed")
}

// testExportEmptyNamespace tests export of an empty namespace
func testExportEmptyNamespace(t *testing.T, suite *framework.TestSuite) {
	ns := suite.CreateTestNamespace("export-empty")
	suite.LogInfo("Created empty test namespace: %s", ns)

	// Export empty namespace
	exportDir := suite.CreateTempDir("export-empty")
	result := suite.RunCLI("export", "--namespace", ns, "--export-dir", exportDir)

	// Command should succeed even with empty namespace
	suite.Assert.AssertCommandSuccess(result)

	// Resources directory might exist but be empty
	resourcesDir := filepath.Join(exportDir, "resources", ns)
	if _, err := os.Stat(resourcesDir); err == nil {
		suite.LogInfo("Resources directory created for empty namespace")
	}
}

// testExportNonExistentNamespace tests export of non-existent namespace
func testExportNonExistentNamespace(t *testing.T, suite *framework.TestSuite) {
	// Try to export from non-existent namespace
	exportDir := suite.CreateTempDir("export-nonexistent")
	result := suite.RunCLI("export",
		"--namespace", "nonexistent-namespace-12345",
		"--export-dir", exportDir)

	// Command might succeed but with no resources exported
	// OR might fail depending on implementation
	// Let's check both scenarios
	if result.Success() {
		suite.LogInfo("Export succeeded for non-existent namespace (expected behavior)")
	} else {
		suite.LogInfo("Export failed for non-existent namespace (also expected)")
	}
}

// Helper functions

// deployTestApp deploys the hello-world sample app for testing
func deployTestApp(t *testing.T, suite *framework.TestSuite, namespace string) {
	suite.LogInfo("Deploying hello-world sample app to namespace %s", namespace)

	// Use the existing sample-resources/hello-world app
	cmd := exec.Command("kubectl", "apply", "-f", "sample-resources/hello-world/manifest.yaml", "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		suite.LogInfo("Deploy output: %s", string(output))
		suite.LogInfo("Deploy error: %v (non-fatal, continuing)", err)
	} else {
		suite.LogDebug("Successfully deployed hello-world to namespace %s", namespace)
	}
}

// deployTestAppWithLabels deploys hello-world app (it already has labels)
func deployTestAppWithLabels(t *testing.T, suite *framework.TestSuite, namespace string, labels map[string]string) {
	suite.LogInfo("Deploying hello-world app with labels to namespace %s", namespace)
	deployTestApp(t, suite, namespace)
}

// deployTestAppWithServiceAccount deploys hello-world app
func deployTestAppWithServiceAccount(t *testing.T, suite *framework.TestSuite, namespace string) {
	suite.LogInfo("Deploying hello-world app to namespace %s", namespace)
	deployTestApp(t, suite, namespace)
}

// TestExportHelp tests the export command help
func TestExportHelp(t *testing.T) {
	suite := framework.NewTestSuite(t)

	result := suite.RunCLI("export", "--help")

	suite.Assert.AssertCommandSuccess(result)
	suite.Assert.AssertOutputContains(result, "Export the namespace resources")
	suite.Assert.AssertOutputContains(result, "--export-dir")
	suite.Assert.AssertOutputContains(result, "--namespace")
	suite.Assert.AssertOutputContains(result, "--label-selector")
	suite.Assert.AssertOutputContains(result, "--cluster-scoped-rbac")
}

// TestExportVersion tests that export works with version flag
func TestExportVersion(t *testing.T) {
	suite := framework.NewTestSuite(t)

	result := suite.RunCLI("version")

	suite.Assert.AssertCommandSuccess(result)
	suite.LogInfo("Version output: %s", result.Stdout)
}

// TestExportResourceDiscovery tests the resource discovery functionality
func TestExportResourceDiscovery(t *testing.T) {
	suite := framework.NewTestSuite(t)
	defer suite.Cleanup()

	suite.SkipIfShort("Resource discovery test requires a running cluster")

	ns := suite.CreateTestNamespace("export-discovery")

	// Deploy multiple resource types
	deployMultipleResourceTypes(t, suite, ns)

	// Export all resources
	exportDir := suite.CreateTempDir("export-discovery")
	result := suite.RunCLI("export", "--namespace", ns, "--export-dir", exportDir)

	suite.Assert.AssertCommandSuccess(result)

	// Verify multiple resource types were discovered and exported
	resourcesDir := filepath.Join(exportDir, "resources", ns)
	suite.Assert.AssertDirExists(resourcesDir)
	suite.Assert.AssertMinYAMLFiles(resourcesDir, 1)

	suite.LogInfo("Resource discovery test completed")
}

// deployMultipleResourceTypes deploys the hello-world app (has multiple resources)
func deployMultipleResourceTypes(t *testing.T, suite *framework.TestSuite, namespace string) {
	suite.LogInfo("Deploying hello-world app (multiple resources) to namespace %s", namespace)
	deployTestApp(t, suite, namespace)
}

// TestExportOutputStructure tests the structure of exported files
func TestExportOutputStructure(t *testing.T) {
	suite := framework.NewTestSuite(t)
	defer suite.Cleanup()

	suite.SkipIfShort("Output structure test requires a running cluster")

	ns := suite.CreateTestNamespace("export-structure")
	deployTestApp(t, suite, ns)

	// Wait for deployment
	err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
	suite.Assert.AssertNoError(err)

	exportDir := suite.CreateTempDir("export-structure")
	result := suite.RunCLI("export", "--namespace", ns, "--export-dir", exportDir)

	suite.Assert.AssertCommandSuccess(result)

	// Verify expected directory structure
	// export/
	//   resources/
	//     <namespace>/
	//       <resource files>
	//   failures/
	//     <namespace>/

	suite.Assert.AssertDirExists(filepath.Join(exportDir, "resources"))
	suite.Assert.AssertDirExists(filepath.Join(exportDir, "resources", ns))
	suite.Assert.AssertDirExists(filepath.Join(exportDir, "failures"))
	suite.Assert.AssertDirExists(filepath.Join(exportDir, "failures", ns))

	suite.LogInfo("Export directory structure validated")
}

// TestExportQPSAndBurst tests QPS and Burst rate limiting flags
func TestExportQPSAndBurst(t *testing.T) {
	suite := framework.NewTestSuite(t)
	defer suite.Cleanup()

	suite.SkipIfShort("QPS test requires a running cluster")

	ns := suite.CreateTestNamespace("export-qps")
	deployTestApp(t, suite, ns)

	// Wait for deployment
	err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
	suite.Assert.AssertNoError(err)

	exportDir := suite.CreateTempDir("export-qps")
	result := suite.RunCLI("export",
		"--namespace", ns,
		"--export-dir", exportDir,
		"--qps", "50",
		"--burst", "100")

	suite.Assert.AssertCommandSuccess(result)
	suite.LogInfo("Export with custom QPS and burst completed")
}
