package applyplan

import (
	"reflect"
	"testing"

	"github.com/tasuku43/gion-core/planner"
)

func TestCollectPrefetchRepoKeys(t *testing.T) {
	changes := []planner.WorkspaceChange{
		{
			Kind:        planner.WorkspaceAdd,
			WorkspaceID: "WS-ADD",
		},
		{
			Kind:        planner.WorkspaceUpdate,
			WorkspaceID: "WS-UPD",
			Repos: []planner.RepoChange{
				{Kind: planner.RepoAdd, Alias: "api", ToRepo: "example.com/org/api"},
				{Kind: planner.RepoUpdate, Alias: "web", ToRepo: "example.com/org/web"},
				{Kind: planner.RepoRemove, Alias: "old", FromRepo: "example.com/org/old"},
			},
		},
	}
	desired := planner.Inventory{
		Workspaces: map[string]planner.Workspace{
			"WS-ADD": {
				ID: "WS-ADD",
				Repos: []planner.Repo{
					{Alias: "api", RepoKey: "example.com/org/api"},
					{Alias: "cli", RepoKey: "example.com/org/cli"},
				},
			},
		},
	}

	got := CollectPrefetchRepoKeys(changes, desired)
	want := []string{
		"example.com/org/api",
		"example.com/org/cli",
		"example.com/org/web",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CollectPrefetchRepoKeys() = %#v, want %#v", got, want)
	}
}
