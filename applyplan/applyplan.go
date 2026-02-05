package applyplan

import (
	"strings"

	"github.com/tasuku43/gion-core/planner"
)

func CountWorkspaceChanges(changes []planner.WorkspaceChange) (adds, updates, removes int) {
	for _, change := range changes {
		switch change.Kind {
		case planner.WorkspaceAdd:
			adds++
		case planner.WorkspaceUpdate:
			updates++
		case planner.WorkspaceRemove:
			removes++
		}
	}
	return adds, updates, removes
}

func HasDestructiveChanges(changes []planner.WorkspaceChange) bool {
	for _, change := range changes {
		switch change.Kind {
		case planner.WorkspaceRemove:
			return true
		case planner.WorkspaceUpdate:
			if HasDestructiveRepoChanges(change.Repos) {
				return true
			}
		}
	}
	return false
}

func HasDestructiveRepoChanges(changes []planner.RepoChange) bool {
	for _, change := range changes {
		switch change.Kind {
		case planner.RepoRemove:
			return true
		case planner.RepoUpdate:
			if IsInPlaceBranchRename(change) {
				continue
			}
			return true
		}
	}
	return false
}

func IsInPlaceBranchRename(change planner.RepoChange) bool {
	if change.Kind != planner.RepoUpdate {
		return false
	}
	fromRepo := strings.TrimSpace(change.FromRepo)
	toRepo := strings.TrimSpace(change.ToRepo)
	fromBranch := strings.TrimSpace(change.FromBranch)
	toBranch := strings.TrimSpace(change.ToBranch)
	if fromRepo == "" || toRepo == "" || fromBranch == "" || toBranch == "" {
		return false
	}
	if fromRepo != toRepo {
		return false
	}
	if fromBranch == toBranch {
		return false
	}
	return true
}
