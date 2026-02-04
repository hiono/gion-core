package importplan

import (
	"reflect"
	"testing"
)

func TestBuildInventory(t *testing.T) {
	inv := BuildInventory([]WorkspaceSnapshot{
		{
			ID:          "WS-1",
			Description: "desc",
			Mode:        "repo",
			PresetName:  "",
			SourceURL:   "",
			BaseBranch:  "origin/main",
			Repos: []RepoSnapshot{
				{Alias: "web", RepoKey: "example.com/org/web", Branch: "WS-1"},
				{Alias: "api", RepoKey: "example.com/org/api", Branch: "WS-1"},
			},
		},
	})
	ws, ok := inv.Workspaces["WS-1"]
	if !ok {
		t.Fatalf("workspace missing")
	}
	if got := len(ws.Repos); got != 2 {
		t.Fatalf("len(repos) = %d, want 2", got)
	}
	if ws.Repos[0].Alias != "api" || ws.Repos[1].Alias != "web" {
		t.Fatalf("repos not sorted: %#v", ws.Repos)
	}
	if ws.Repos[0].BaseRef != "origin/main" {
		t.Fatalf("base_ref not set: %#v", ws.Repos[0])
	}
}

func TestCollectWorkspaceIDs(t *testing.T) {
	got := CollectWorkspaceIDs([]string{"WS-2", "", "WS-1"})
	want := []string{"WS-1", "WS-2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CollectWorkspaceIDs() = %#v, want %#v", got, want)
	}
}
