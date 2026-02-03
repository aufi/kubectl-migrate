package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestSuite provides a complete test environment for E2E tests
type TestSuite struct {
	T          *testing.T
	CLI        *CLIExecutor
	Cluster    *ClusterManager
	Assert     *Assertions
	TempDir    string
	FixturesDir string
	cleanupFns []func()
}

// NewTestSuite creates a new test suite with all necessary components
func NewTestSuite(t *testing.T) *TestSuite {
	cluster, err := NewClusterManager()
	if err != nil {
		t.Fatalf("Failed to create cluster manager: %v", err)
	}

	suite := &TestSuite{
		T:           t,
		CLI:         NewCLIExecutor(os.Getenv("E2E_BINARY")),
		Cluster:     cluster,
		Assert:      NewAssertions(t),
		TempDir:     t.TempDir(),
		FixturesDir: "test/fixtures",
		cleanupFns:  []func(){},
	}

	return suite
}

// NewTestSuiteWithContext creates a test suite with a specific kubeconfig context
func NewTestSuiteWithContext(t *testing.T, contextName string) *TestSuite {
	cluster, err := NewClusterManagerWithContext(contextName)
	if err != nil {
		t.Fatalf("Failed to create cluster manager with context %s: %v", contextName, err)
	}

	suite := &TestSuite{
		T:           t,
		CLI:         NewCLIExecutor(os.Getenv("E2E_BINARY")),
		Cluster:     cluster,
		Assert:      NewAssertions(t),
		TempDir:     t.TempDir(),
		FixturesDir: "test/fixtures",
		cleanupFns:  []func(){},
	}

	return suite
}

// LoadFixture loads a fixture file and returns its content
func (s *TestSuite) LoadFixture(relativePath string) string {
	path := filepath.Join(s.FixturesDir, relativePath)
	content, err := os.ReadFile(path)
	if err != nil {
		s.T.Fatalf("Failed to load fixture %s: %v", path, err)
	}
	return string(content)
}

// GetFixturePath returns the absolute path to a fixture
func (s *TestSuite) GetFixturePath(relativePath string) string {
	absPath, err := filepath.Abs(filepath.Join(s.FixturesDir, relativePath))
	if err != nil {
		s.T.Fatalf("Failed to get absolute path for fixture %s: %v", relativePath, err)
	}
	return absPath
}

// CreateTempFile creates a temporary file with content
func (s *TestSuite) CreateTempFile(name, content string) string {
	path := filepath.Join(s.TempDir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		s.T.Fatalf("Failed to create temp file %s: %v", path, err)
	}
	return path
}

// CreateTempDir creates a temporary subdirectory
func (s *TestSuite) CreateTempDir(name string) string {
	path := filepath.Join(s.TempDir, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		s.T.Fatalf("Failed to create temp directory %s: %v", path, err)
	}
	return path
}

// RegisterCleanup registers a cleanup function to be called after the test
func (s *TestSuite) RegisterCleanup(fn func()) {
	s.cleanupFns = append(s.cleanupFns, fn)
}

// Cleanup runs all registered cleanup functions
func (s *TestSuite) Cleanup() {
	for i := len(s.cleanupFns) - 1; i >= 0; i-- {
		s.cleanupFns[i]()
	}
}

// CreateTestNamespace creates a test namespace and registers cleanup
func (s *TestSuite) CreateTestNamespace(baseName string) string {
	// Generate unique namespace name
	namespaceName := fmt.Sprintf("%s-e2e-%d", baseName, os.Getpid())

	err := s.Cluster.CreateNamespace(namespaceName)
	if err != nil {
		s.T.Fatalf("Failed to create test namespace %s: %v", namespaceName, err)
	}

	// Register cleanup
	s.RegisterCleanup(func() {
		s.T.Logf("Cleaning up namespace %s", namespaceName)
		if err := s.Cluster.DeleteNamespace(namespaceName); err != nil {
			s.T.Logf("Warning: Failed to delete namespace %s: %v", namespaceName, err)
		}
	})

	return namespaceName
}

// RunCLI is a convenience method to run a CLI command
func (s *TestSuite) RunCLI(args ...string) *CLIResult {
	return s.CLI.Run(args...)
}

// LogInfo logs an informational message
func (s *TestSuite) LogInfo(format string, args ...interface{}) {
	s.T.Logf("INFO: "+format, args...)
}

// LogDebug logs a debug message (only shown in verbose mode)
func (s *TestSuite) LogDebug(format string, args ...interface{}) {
	if testing.Verbose() {
		s.T.Logf("DEBUG: "+format, args...)
	}
}

// Skip skips the test with a message
func (s *TestSuite) Skip(reason string) {
	s.T.Skip(reason)
}

// SkipIfShort skips the test if running in short mode
func (s *TestSuite) SkipIfShort(reason string) {
	if testing.Short() {
		s.T.Skipf("Skipping in short mode: %s", reason)
	}
}
