package applyplan

import (
	"testing"

	"github.com/tasuku43/gion-core/planner"
)

func TestCountWorkspaceChanges(t *testing.T) {
	adds, updates, removes := CountWorkspaceChanges([]planner.WorkspaceChange{
		{Kind: planner.WorkspaceAdd},
		{Kind: planner.WorkspaceUpdate},
		{Kind: planner.WorkspaceRemove},
		{Kind: planner.WorkspaceUpdate},
	})
	if adds != 1 || updates != 2 || removes != 1 {
		t.Fatalf("counts = (%d,%d,%d), want (1,2,1)", adds, updates, removes)
	}
}

func TestHasDestructiveChanges(t *testing.T) {
	tests := []struct {
		name    string
		changes []planner.WorkspaceChange
		want    bool
	}{
		{
			name: "workspace remove is destructive",
			changes: []planner.WorkspaceChange{
				{Kind: planner.WorkspaceRemove},
			},
			want: true,
		},
		{
			name: "repo remove is destructive",
			changes: []planner.WorkspaceChange{
				{Kind: planner.WorkspaceUpdate, Repos: []planner.RepoChange{{Kind: planner.RepoRemove}}},
			},
			want: true,
		},
		{
			name: "in-place branch rename is not destructive",
			changes: []planner.WorkspaceChange{
				{Kind: planner.WorkspaceUpdate, Repos: []planner.RepoChange{{
					Kind:       planner.RepoUpdate,
					FromRepo:   "example.com/org/repo",
					ToRepo:     "example.com/org/repo",
					FromBranch: "WS-1",
					ToBranch:   "WS-2",
				}}},
			},
			want: false,
		},
		{
			name: "repo move is destructive",
			changes: []planner.WorkspaceChange{
				{Kind: planner.WorkspaceUpdate, Repos: []planner.RepoChange{{
					Kind:       planner.RepoUpdate,
					FromRepo:   "example.com/org/repo-a",
					ToRepo:     "example.com/org/repo-b",
					FromBranch: "WS-1",
					ToBranch:   "WS-1",
				}}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasDestructiveChanges(tt.changes); got != tt.want {
				t.Fatalf("HasDestructiveChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInPlaceBranchRename(t *testing.T) {
	tests := []struct {
		name   string
		change planner.RepoChange
		want   bool
	}{
		{
			name: "same repo and different branch",
			change: planner.RepoChange{
				Kind:       planner.RepoUpdate,
				FromRepo:   "example.com/org/repo",
				ToRepo:     "example.com/org/repo",
				FromBranch: "WS-1",
				ToBranch:   "WS-2",
			},
			want: true,
		},
		{
			name: "same branch is false",
			change: planner.RepoChange{
				Kind:       planner.RepoUpdate,
				FromRepo:   "example.com/org/repo",
				ToRepo:     "example.com/org/repo",
				FromBranch: "WS-1",
				ToBranch:   "WS-1",
			},
			want: false,
		},
		{
			name: "different repo is false",
			change: planner.RepoChange{
				Kind:       planner.RepoUpdate,
				FromRepo:   "example.com/org/repo-a",
				ToRepo:     "example.com/org/repo-b",
				FromBranch: "WS-1",
				ToBranch:   "WS-2",
			},
			want: false,
		},
		{
			name: "empty field is false",
			change: planner.RepoChange{
				Kind:       planner.RepoUpdate,
				FromRepo:   "example.com/org/repo",
				ToRepo:     "example.com/org/repo",
				FromBranch: "",
				ToBranch:   "WS-2",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInPlaceBranchRename(tt.change); got != tt.want {
				t.Fatalf("IsInPlaceBranchRename() = %v, want %v", got, tt.want)
			}
		})
	}
}
