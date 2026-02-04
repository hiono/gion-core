package repostore

import (
	"path/filepath"
	"testing"
	"time"
)

func TestFetchGraceDuration(t *testing.T) {
	if got := FetchGraceDuration(""); got != 30*time.Second {
		t.Fatalf("default grace = %v", got)
	}
	if got := FetchGraceDuration("15"); got != 15*time.Second {
		t.Fatalf("parsed grace = %v", got)
	}
	if got := FetchGraceDuration("-1"); got != 30*time.Second {
		t.Fatalf("negative grace fallback = %v", got)
	}
}

func TestTouchFetchHeadAndRecentlyFetched(t *testing.T) {
	storePath := t.TempDir()
	if err := TouchFetchHead(storePath); err != nil {
		t.Fatalf("TouchFetchHead() error = %v", err)
	}
	if !RecentlyFetched(storePath, 1*time.Minute) {
		t.Fatalf("RecentlyFetched() = false, want true")
	}
	if RecentlyFetched(filepath.Join(storePath, "missing"), 1*time.Minute) {
		t.Fatalf("RecentlyFetched() on missing path = true")
	}
}
