package repostore

import (
	"path/filepath"
	"testing"

	"github.com/hiono/gion-core/repospec"
)

func TestStorePath(t *testing.T) {
	bareRoot := filepath.Join("/tmp", "bare")
	spec := repospec.Spec{
		EndPoint: repospec.EndPoint{
			Host: "github.com",
		},
		Registry: repospec.Registry{
			Group: "org",
		},
		Repository: repospec.Repository{
			Repo: "repo",
		},
	}
	got := StorePath(bareRoot, spec)
	// Level1: github.com/org, Level2: repo
	want := filepath.Join(bareRoot, "github.com/org", "repo.git")
	if got != want {
		t.Fatalf("StorePath() = %q, want %q", got, want)
	}
}
