package repostore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestList(t *testing.T) {
	root := t.TempDir()
	reposRoot := filepath.Join(root, "bare")
	if err := os.MkdirAll(filepath.Join(reposRoot, "github.com", "org", "repo.git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(reposRoot, "github.com", "org", "not-git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	entries, warnings, err := List(reposRoot)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none", warnings)
	}
	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
	if entries[0].RepoKey != "github.com/org/repo.git" {
		t.Fatalf("RepoKey = %q", entries[0].RepoKey)
	}
}

func TestList_NoRoot(t *testing.T) {
	entries, warnings, err := List(filepath.Join(t.TempDir(), "missing"))
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none", warnings)
	}
	if len(entries) != 0 {
		t.Fatalf("entries = %v, want empty", entries)
	}
}
