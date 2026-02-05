package applyplan

import (
	"testing"

	"github.com/tasuku43/gion-core/planner"
)

func TestBuildExecution(t *testing.T) {
	changes := []planner.WorkspaceChange{
		{
			Kind:        planner.WorkspaceAdd,
			WorkspaceID: "WS-ADD",
		},
		{
			Kind:        planner.WorkspaceRemove,
			WorkspaceID: "WS-REMOVE",
		},
		{
			Kind:        planner.WorkspaceUpdate,
			WorkspaceID: "WS-UPDATE",
			Repos: []planner.RepoChange{
				{
					Kind:       planner.RepoRemove,
					Alias:      "repo-remove",
					FromRepo:   "example.com/org/remove",
					FromBranch: "main",
				},
				{
					Kind:       planner.RepoUpdate,
					Alias:      "repo-rename",
					FromRepo:   "example.com/org/rename",
					ToRepo:     "example.com/org/rename",
					FromBranch: "WS-1",
					ToBranch:   "WS-2",
				},
				{
					Kind:       planner.RepoUpdate,
					Alias:      "repo-move",
					FromRepo:   "example.com/org/a",
					ToRepo:     "example.com/org/b",
					FromBranch: "WS-1",
					ToBranch:   "WS-1",
				},
				{
					Kind:     planner.RepoAdd,
					Alias:    "repo-add",
					ToRepo:   "example.com/org/add",
					ToBranch: "WS-1",
				},
			},
		},
	}

	exec := BuildExecution(changes)

	if len(exec.WorkspaceRemovals) != 1 || exec.WorkspaceRemovals[0].WorkspaceID != "WS-REMOVE" {
		t.Fatalf("WorkspaceRemovals = %+v, want WS-REMOVE", exec.WorkspaceRemovals)
	}
	if len(exec.WorkspaceAdds) != 1 || exec.WorkspaceAdds[0].WorkspaceID != "WS-ADD" {
		t.Fatalf("WorkspaceAdds = %+v, want WS-ADD", exec.WorkspaceAdds)
	}
	if len(exec.WorkspaceUpdateRemovals) != 1 {
		t.Fatalf("len(WorkspaceUpdateRemovals) = %d, want 1", len(exec.WorkspaceUpdateRemovals))
	}
	if got := len(exec.WorkspaceUpdateRemovals[0].Repos); got != 2 {
		t.Fatalf("len(WorkspaceUpdateRemovals[0].Repos) = %d, want 2", got)
	}
	if got := len(exec.WorkspaceUpdateRenames); got != 1 {
		t.Fatalf("len(WorkspaceUpdateRenames) = %d, want 1", got)
	}
	if got := len(exec.WorkspaceUpdateRenames[0].Repos); got != 1 {
		t.Fatalf("len(WorkspaceUpdateRenames[0].Repos) = %d, want 1", got)
	}
	if got := len(exec.WorkspaceUpdateAdds); got != 1 {
		t.Fatalf("len(WorkspaceUpdateAdds) = %d, want 1", got)
	}
	if got := len(exec.WorkspaceUpdateAdds[0].Repos); got != 2 {
		t.Fatalf("len(WorkspaceUpdateAdds[0].Repos) = %d, want 2", got)
	}
}
