package planner

import "testing"

func TestDiff_AddUpdateRemove(t *testing.T) {
	desired := Inventory{
		Workspaces: map[string]Workspace{
			"ws-add": {
				Repos: []Repo{{Alias: "api", RepoKey: "gh:org/api", Branch: "feat-1"}},
			},
			"ws-update": {
				Repos: []Repo{
					{Alias: "api", RepoKey: "gh:org/api", Branch: "feat-2"},
					{Alias: "web", RepoKey: "gh:org/web", Branch: "feat-2"},
				},
			},
		},
	}
	actual := Inventory{
		Workspaces: map[string]Workspace{
			"ws-remove": {
				Repos: []Repo{{Alias: "api", RepoKey: "gh:org/api", Branch: "old"}},
			},
			"ws-update": {
				Repos: []Repo{
					{Alias: "api", RepoKey: "gh:org/api", Branch: "main"},
					{Alias: "old", RepoKey: "gh:org/old", Branch: "main"},
				},
			},
		},
	}

	changes := Diff(desired, actual)
	if len(changes) != 3 {
		t.Fatalf("len(changes) = %d, want 3", len(changes))
	}

	if changes[0].WorkspaceID != "ws-add" || changes[0].Kind != WorkspaceAdd {
		t.Fatalf("changes[0] = %+v, want ws-add/add", changes[0])
	}
	if len(changes[0].Repos) != 1 || changes[0].Repos[0].Kind != RepoAdd {
		t.Fatalf("changes[0].Repos = %+v, want one repo add", changes[0].Repos)
	}

	if changes[1].WorkspaceID != "ws-remove" || changes[1].Kind != WorkspaceRemove {
		t.Fatalf("changes[1] = %+v, want ws-remove/remove", changes[1])
	}

	if changes[2].WorkspaceID != "ws-update" || changes[2].Kind != WorkspaceUpdate {
		t.Fatalf("changes[2] = %+v, want ws-update/update", changes[2])
	}
	if len(changes[2].Repos) != 3 {
		t.Fatalf("len(changes[2].Repos) = %d, want 3", len(changes[2].Repos))
	}

	if changes[2].Repos[0].Alias != "api" || changes[2].Repos[0].Kind != RepoUpdate {
		t.Fatalf("changes[2].Repos[0] = %+v, want api/update", changes[2].Repos[0])
	}
	if changes[2].Repos[1].Alias != "old" || changes[2].Repos[1].Kind != RepoRemove {
		t.Fatalf("changes[2].Repos[1] = %+v, want old/remove", changes[2].Repos[1])
	}
	if changes[2].Repos[2].Alias != "web" || changes[2].Repos[2].Kind != RepoAdd {
		t.Fatalf("changes[2].Repos[2] = %+v, want web/add", changes[2].Repos[2])
	}
}

func TestWorkspaceChangeKindString(t *testing.T) {
	if got := WorkspaceAdd.String(); got != "add" {
		t.Fatalf("WorkspaceAdd.String() = %q, want add", got)
	}
	unknown := WorkspaceChangeKind("x")
	if got := unknown.String(); got != "unknown(x)" {
		t.Fatalf("unknown.String() = %q, want unknown(x)", got)
	}
}
