package workspacerisk

import (
	"errors"
	"testing"
)

func TestClassifyRepoStatus(t *testing.T) {
	t.Run("error is unknown", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{Error: errors.New("boom")})
		if got != RepoStateUnknown {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateUnknown)
		}
	})

	t.Run("dirty is dirty", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{Dirty: true})
		if got != RepoStateDirty {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateDirty)
		}
	})

	t.Run("detached is unknown", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{Detached: true})
		if got != RepoStateUnknown {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateUnknown)
		}
	})

	t.Run("missing upstream is unknown", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{})
		if got != RepoStateUnknown {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateUnknown)
		}
	})

	t.Run("diverged", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{Upstream: "origin/main", AheadCount: 1, BehindCount: 1})
		if got != RepoStateDiverged {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateDiverged)
		}
	})

	t.Run("unpushed", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{Upstream: "origin/main", AheadCount: 2})
		if got != RepoStateUnpushed {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateUnpushed)
		}
	})

	t.Run("behind only is clean", func(t *testing.T) {
		got := ClassifyRepoStatus(RepoStatus{Upstream: "origin/main", BehindCount: 2})
		if got != RepoStateClean {
			t.Fatalf("ClassifyRepoStatus() = %q, want %q", got, RepoStateClean)
		}
	})
}
