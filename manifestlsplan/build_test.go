package manifestlsplan

import (
	"reflect"
	"testing"

	"github.com/tasuku43/gion-core/planner"
)

func TestBuild(t *testing.T) {
	result := Build(
		[]DesiredWorkspace{
			{ID: "WS-2", Description: "two"},
			{ID: "WS-1", Description: "one"},
		},
		[]planner.WorkspaceChange{
			{Kind: planner.WorkspaceAdd, WorkspaceID: "WS-2"},
			{Kind: planner.WorkspaceUpdate, WorkspaceID: "WS-1"},
			{Kind: planner.WorkspaceRemove, WorkspaceID: "WS-X"},
		},
		[]string{"WS-1", "WS-X"},
	)

	if got, want := result.Counts, (Counts{Applied: 0, Drift: 1, Missing: 1, Extra: 1}); got != want {
		t.Fatalf("Counts = %#v, want %#v", got, want)
	}

	gotEntries := []ManifestEntry{
		result.ManifestEntries[0],
		result.ManifestEntries[1],
	}
	wantEntries := []ManifestEntry{
		{WorkspaceID: "WS-1", Drift: DriftDrift, Description: "one", HasWorkspace: true},
		{WorkspaceID: "WS-2", Drift: DriftMissing, Description: "two", HasWorkspace: false},
	}
	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Fatalf("ManifestEntries = %#v, want %#v", gotEntries, wantEntries)
	}
	if len(result.ExtraEntries) != 1 || result.ExtraEntries[0].WorkspaceID != "WS-X" {
		t.Fatalf("ExtraEntries = %#v", result.ExtraEntries)
	}
}
