package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

// Assertions provides assertion helpers for E2E tests
type Assertions struct {
	t *testing.T
}

// NewAssertions creates a new Assertions helper
func NewAssertions(t *testing.T) *Assertions {
	return &Assertions{t: t}
}

// AssertCommandSuccess asserts that a command succeeded
func (a *Assertions) AssertCommandSuccess(result *CLIResult) {
	a.t.Helper()
	if !result.Success() {
		a.t.Errorf("Command failed with exit code %d\nStdout: %s\nStderr: %s\nError: %v",
			result.ExitCode, result.Stdout, result.Stderr, result.Err)
	}
}

// AssertCommandFails asserts that a command failed
func (a *Assertions) AssertCommandFails(result *CLIResult) {
	a.t.Helper()
	if result.Success() {
		a.t.Errorf("Expected command to fail, but it succeeded\nStdout: %s\nStderr: %s",
			result.Stdout, result.Stderr)
	}
}

// AssertExitCode asserts that a command has a specific exit code
func (a *Assertions) AssertExitCode(result *CLIResult, expectedCode int) {
	a.t.Helper()
	if result.ExitCode != expectedCode {
		a.t.Errorf("Expected exit code %d, got %d\nStdout: %s\nStderr: %s",
			expectedCode, result.ExitCode, result.Stdout, result.Stderr)
	}
}

// AssertOutputContains asserts that output contains a substring
func (a *Assertions) AssertOutputContains(result *CLIResult, substring string) {
	a.t.Helper()
	if !result.ContainsOutput(substring) {
		a.t.Errorf("Expected output to contain '%s'\nStdout: %s\nStderr: %s",
			substring, result.Stdout, result.Stderr)
	}
}

// AssertOutputNotContains asserts that output does not contain a substring
func (a *Assertions) AssertOutputNotContains(result *CLIResult, substring string) {
	a.t.Helper()
	if result.ContainsOutput(substring) {
		a.t.Errorf("Expected output to NOT contain '%s'\nStdout: %s\nStderr: %s",
			substring, result.Stdout, result.Stderr)
	}
}

// AssertStdoutContains asserts that stdout contains a substring
func (a *Assertions) AssertStdoutContains(result *CLIResult, substring string) {
	a.t.Helper()
	if !result.ContainsStdout(substring) {
		a.t.Errorf("Expected stdout to contain '%s'\nStdout: %s",
			substring, result.Stdout)
	}
}

// AssertStderrContains asserts that stderr contains a substring
func (a *Assertions) AssertStderrContains(result *CLIResult, substring string) {
	a.t.Helper()
	if !result.ContainsStderr(substring) {
		a.t.Errorf("Expected stderr to contain '%s'\nStderr: %s",
			substring, result.Stderr)
	}
}

// AssertOutputMatches asserts that output matches a regex pattern
func (a *Assertions) AssertOutputMatches(result *CLIResult, pattern string) {
	a.t.Helper()
	matched, err := regexp.MatchString(pattern, result.Output())
	if err != nil {
		a.t.Errorf("Invalid regex pattern '%s': %v", pattern, err)
		return
	}
	if !matched {
		a.t.Errorf("Expected output to match pattern '%s'\nOutput: %s",
			pattern, result.Output())
	}
}

// AssertFileExists asserts that a file exists
func (a *Assertions) AssertFileExists(path string) {
	a.t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		a.t.Errorf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists asserts that a file does not exist
func (a *Assertions) AssertFileNotExists(path string) {
	a.t.Helper()
	if _, err := os.Stat(path); err == nil {
		a.t.Errorf("Expected file to NOT exist: %s", path)
	}
}

// AssertDirExists asserts that a directory exists
func (a *Assertions) AssertDirExists(path string) {
	a.t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		a.t.Errorf("Expected directory to exist: %s", path)
		return
	}
	if !info.IsDir() {
		a.t.Errorf("Expected path to be a directory: %s", path)
	}
}

// AssertDirNotEmpty asserts that a directory is not empty
func (a *Assertions) AssertDirNotEmpty(path string) {
	a.t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		a.t.Errorf("Failed to read directory %s: %v", path, err)
		return
	}
	if len(entries) == 0 {
		a.t.Errorf("Expected directory to not be empty: %s", path)
	}
}

// AssertFileContains asserts that a file contains a substring
func (a *Assertions) AssertFileContains(path, substring string) {
	a.t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		a.t.Errorf("Failed to read file %s: %v", path, err)
		return
	}
	if !strings.Contains(string(content), substring) {
		a.t.Errorf("Expected file %s to contain '%s'", path, substring)
	}
}

// AssertYAMLFileValid asserts that a file is valid YAML
func (a *Assertions) AssertYAMLFileValid(path string) {
	a.t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		a.t.Errorf("Failed to read file %s: %v", path, err)
		return
	}

	var data interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		a.t.Errorf("File %s is not valid YAML: %v", path, err)
	}
}

// AssertFilesInDir asserts that a directory contains a specific number of files
func (a *Assertions) AssertFilesInDir(path string, count int) {
	a.t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		a.t.Errorf("Failed to read directory %s: %v", path, err)
		return
	}

	fileCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			fileCount++
		}
	}

	if fileCount != count {
		a.t.Errorf("Expected %d files in directory %s, got %d", count, path, fileCount)
	}
}

// AssertResourceExists asserts that a resource exists in the cluster
func (a *Assertions) AssertResourceExists(cluster *ClusterManager, gvr schema.GroupVersionResource, namespace, name string) {
	a.t.Helper()
	exists, err := cluster.ResourceExists(gvr, namespace, name)
	if err != nil {
		a.t.Errorf("Failed to check if resource exists: %v", err)
		return
	}
	if !exists {
		a.t.Errorf("Expected resource to exist: %s/%s in namespace %s", gvr.Resource, name, namespace)
	}
}

// AssertResourceNotExists asserts that a resource does not exist in the cluster
func (a *Assertions) AssertResourceNotExists(cluster *ClusterManager, gvr schema.GroupVersionResource, namespace, name string) {
	a.t.Helper()
	exists, err := cluster.ResourceExists(gvr, namespace, name)
	if err != nil {
		a.t.Errorf("Failed to check if resource exists: %v", err)
		return
	}
	if exists {
		a.t.Errorf("Expected resource to NOT exist: %s/%s in namespace %s", gvr.Resource, name, namespace)
	}
}

// AssertResourceHasLabel asserts that a resource has a specific label
func (a *Assertions) AssertResourceHasLabel(cluster *ClusterManager, gvr schema.GroupVersionResource, namespace, name, key, value string) {
	a.t.Helper()
	resource, err := cluster.GetResource(gvr, namespace, name)
	if err != nil {
		a.t.Errorf("Failed to get resource: %v", err)
		return
	}

	labels := resource.GetLabels()
	if labels == nil {
		a.t.Errorf("Resource has no labels")
		return
	}

	actualValue, exists := labels[key]
	if !exists {
		a.t.Errorf("Resource does not have label '%s'", key)
		return
	}

	if actualValue != value {
		a.t.Errorf("Expected label '%s' to have value '%s', got '%s'", key, value, actualValue)
	}
}

// AssertResourceCount asserts the number of resources matching a label selector
func (a *Assertions) AssertResourceCount(cluster *ClusterManager, gvr schema.GroupVersionResource, namespace, labelSelector string, expectedCount int) {
	a.t.Helper()
	count, err := cluster.CountResources(gvr, namespace, labelSelector)
	if err != nil {
		a.t.Errorf("Failed to count resources: %v", err)
		return
	}

	if count != expectedCount {
		a.t.Errorf("Expected %d resources, got %d", expectedCount, count)
	}
}

// AssertNoError asserts that an error is nil
func (a *Assertions) AssertNoError(err error) {
	a.t.Helper()
	if err != nil {
		a.t.Errorf("Expected no error, got: %v", err)
	}
}

// AssertError asserts that an error is not nil
func (a *Assertions) AssertError(err error) {
	a.t.Helper()
	if err == nil {
		a.t.Errorf("Expected an error, got nil")
	}
}

// AssertErrorContains asserts that an error message contains a substring
func (a *Assertions) AssertErrorContains(err error, substring string) {
	a.t.Helper()
	if err == nil {
		a.t.Errorf("Expected an error, got nil")
		return
	}
	if !strings.Contains(err.Error(), substring) {
		a.t.Errorf("Expected error to contain '%s', got: %v", substring, err)
	}
}

// CountYAMLFiles counts YAML files in a directory recursively
func CountYAMLFiles(dirPath string) (int, error) {
	count := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			count++
		}
		return nil
	})
	return count, err
}

// AssertMinYAMLFiles asserts minimum number of YAML files in a directory
func (a *Assertions) AssertMinYAMLFiles(dirPath string, minCount int) {
	a.t.Helper()
	count, err := CountYAMLFiles(dirPath)
	if err != nil {
		a.t.Errorf("Failed to count YAML files: %v", err)
		return
	}
	if count < minCount {
		a.t.Errorf("Expected at least %d YAML files in %s, got %d", minCount, dirPath, count)
	}
}

// LogInfo logs an informational message
func (a *Assertions) LogInfo(format string, args ...interface{}) {
	a.t.Helper()
	a.t.Logf("INFO: "+format, args...)
}

// LogDebug logs a debug message
func (a *Assertions) LogDebug(format string, args ...interface{}) {
	a.t.Helper()
	if testing.Verbose() {
		a.t.Logf("DEBUG: "+format, args...)
	}
}

// Fail fails the test with a message
func (a *Assertions) Fail(format string, args ...interface{}) {
	a.t.Helper()
	a.t.Errorf(format, args...)
}

// FailNow fails the test immediately with a message
func (a *Assertions) FailNow(format string, args ...interface{}) {
	a.t.Helper()
	a.t.Fatalf(format, args...)
}

// AssertEqual asserts that two values are equal
func (a *Assertions) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if expected != actual {
		msg := fmt.Sprintf("Expected %v, got %v", expected, actual)
		if len(msgAndArgs) > 0 {
			msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...) + ": " + msg
		}
		a.t.Error(msg)
	}
}

// AssertNotEqual asserts that two values are not equal
func (a *Assertions) AssertNotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	a.t.Helper()
	if expected == actual {
		msg := fmt.Sprintf("Expected values to be different, but both are %v", expected)
		if len(msgAndArgs) > 0 {
			msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...) + ": " + msg
		}
		a.t.Error(msg)
	}
}
