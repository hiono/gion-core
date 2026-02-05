package repostore

import (
	"path/filepath"
	"testing"

	"github.com/tasuku43/gion-core/repospec"
)

func TestStorePath(t *testing.T) {
	bareRoot := filepath.Join("/tmp", "bare")
	spec := repospec.Spec{Host: "github.com", Owner: "org", Repo: "repo"}
	got := StorePath(bareRoot, spec)
	want := filepath.Join(bareRoot, "github.com", "org", "repo.git")
	if got != want {
		t.Fatalf("StorePath() = %q, want %q", got, want)
	}
}
