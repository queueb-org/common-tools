package tools

import (
	"testing"
)

const (
	testDir  = "testdata/example-dir"
	testFile = "testdata/example.txt"
)

// TestOSSuite runs tests for Exists, IsDir, IsFile
func TestOSSuite(t *testing.T) {
	nonExistent := "non existent"
	for _, entry := range []struct {
		name     string
		in       string
		exister  func(string) bool
		expected bool
	}{
		// Exists
		{"exists/ok-dir", testDir, Exists, true},
		{"exists/ok-file", testFile, Exists, true},
		{"exists/not", nonExistent, Exists, false},
		// IsDir
		{"is-dir/ok", testDir, IsDir, true},
		{"is-dir/not-a-dir", testFile, IsDir, false},
		{"is-dir/not-exist", nonExistent, IsDir, false},
		// IsFile
		{"is-file/ok", testFile, IsFile, true},
		{"is-file/not-a-file", testDir, IsFile, false},
		{"is-file/not-exist", nonExistent, IsFile, false},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := entry.exister(entry.in); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}
