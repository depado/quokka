package renderer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyze_WithUnreadableFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".quokka.yml"), []byte("name: test\n"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	subdir := filepath.Join(dir, "restricted")

	if err := os.Mkdir(subdir, 0000); err != nil {
		t.Fatalf("Failed to create restricted directory: %v", err)
	}
	defer os.Chmod(subdir, 0755) //nolint:errcheck

	if err := Analyze(dir, "output", "", []string{}); err == nil {
		t.Error("Expected Analyze to return an error for unreadable directory, but got nil")
	}
}
