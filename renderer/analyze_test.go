package renderer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyze_WithUnreadableFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".quokka.yml"), []byte("name: test\n"), 0644)
	subdir := filepath.Join(dir, "restricted")
	os.Mkdir(subdir, 0000)
	defer os.Chmod(subdir, 0755)
	Analyze(dir, "output", "", []string{})
}
