package renderer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/depado/quokka/conf"
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

	if err := Analyze(dir, "output", "", []string{}, 1, conf.InputCtx{}, false, true); err == nil {
		t.Error("Expected Analyze to return an error for unreadable directory, but got nil")
	}
}

func TestAnalyze_IgnoresList(t *testing.T) {
	dir := t.TempDir()
	out := t.TempDir()

	files := map[string]string{
		".quokka.yml":     "name: test\nignores:\n  - README.md\n  - \"docs/*\"\n",
		"README.md":       "template docs",
		"keep.txt":        "kept",
		"docs/ignored.md": "ignored",
	}
	for name, content := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := Analyze(dir, out, "", []string{}, 1, conf.InputCtx{}, false, true); err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	for _, ignored := range []string{"README.md", "docs/ignored.md"} {
		if _, err := os.Stat(filepath.Join(out, ignored)); !os.IsNotExist(err) {
			t.Errorf("Expected %s to be ignored, but it exists", ignored)
		}
	}
	if _, err := os.Stat(filepath.Join(out, "keep.txt")); err != nil {
		t.Errorf("Expected keep.txt to be rendered: %v", err)
	}
}

func TestRunCommands(t *testing.T) {
	out := t.TempDir()
	cmds := []conf.Command{
		{Cmd: "touch ran.txt"},
		{Cmd: "touch skipped.txt", If: "nope"},
		{Cmd: "false", Failure: "stop"},
	}

	if err := runCommands(cmds, out, map[string]any{"nope": false}, true, true); err != nil {
		t.Fatalf("noCommands should skip everything: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "ran.txt")); !os.IsNotExist(err) {
		t.Error("Expected no command to run with noCommands")
	}

	if err := runCommands(cmds, out, map[string]any{"nope": false}, true, false); err == nil {
		t.Error("Expected failure: stop to return an error")
	}
	if _, err := os.Stat(filepath.Join(out, "ran.txt")); err != nil {
		t.Errorf("Expected ran.txt to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "skipped.txt")); !os.IsNotExist(err) {
		t.Error("Expected conditional command to be skipped")
	}
}
