package applyplan

import "github.com/tasuku43/gion-core/planner"

type Execution struct {
	WorkspaceRemovals       []planner.WorkspaceChange
	WorkspaceUpdateRemovals []planner.WorkspaceChange
	WorkspaceUpdateRenames  []planner.WorkspaceChange
	WorkspaceAdds           []planner.WorkspaceChange
	WorkspaceUpdateAdds     []planner.WorkspaceChange
}

func BuildExecution(changes []planner.WorkspaceChange) Execution {
	exec := Execution{}
	for _, change := range changes {
		switch change.Kind {
		case planner.WorkspaceRemove:
			exec.WorkspaceRemovals = append(exec.WorkspaceRemovals, change)
		case planner.WorkspaceAdd:
			exec.WorkspaceAdds = append(exec.WorkspaceAdds, change)
		case planner.WorkspaceUpdate:
			removals := filterRepoRemovals(change.Repos)
			if len(removals) > 0 {
				exec.WorkspaceUpdateRemovals = append(exec.WorkspaceUpdateRemovals, planner.WorkspaceChange{
					Kind:        change.Kind,
					WorkspaceID: change.WorkspaceID,
					Repos:       removals,
				})
			}
			renames := filterRepoRenames(change.Repos)
			if len(renames) > 0 {
				exec.WorkspaceUpdateRenames = append(exec.WorkspaceUpdateRenames, planner.WorkspaceChange{
					Kind:        change.Kind,
					WorkspaceID: change.WorkspaceID,
					Repos:       renames,
				})
			}
			adds := filterRepoAdds(change.Repos)
			if len(adds) > 0 {
				exec.WorkspaceUpdateAdds = append(exec.WorkspaceUpdateAdds, planner.WorkspaceChange{
					Kind:        change.Kind,
					WorkspaceID: change.WorkspaceID,
					Repos:       adds,
				})
			}
		}
	}
	return exec
}

func filterRepoRemovals(changes []planner.RepoChange) []planner.RepoChange {
	var filtered []planner.RepoChange
	for _, change := range changes {
		switch change.Kind {
		case planner.RepoRemove:
			filtered = append(filtered, change)
		case planner.RepoUpdate:
			if IsInPlaceBranchRename(change) {
				continue
			}
			filtered = append(filtered, change)
		}
	}
	return filtered
}

func filterRepoRenames(changes []planner.RepoChange) []planner.RepoChange {
	var filtered []planner.RepoChange
	for _, change := range changes {
		if IsInPlaceBranchRename(change) {
			filtered = append(filtered, change)
		}
	}
	return filtered
}

func filterRepoAdds(changes []planner.RepoChange) []planner.RepoChange {
	var filtered []planner.RepoChange
	for _, change := range changes {
		switch change.Kind {
		case planner.RepoAdd:
			filtered = append(filtered, change)
		case planner.RepoUpdate:
			if IsInPlaceBranchRename(change) {
				continue
			}
			filtered = append(filtered, change)
		}
	}
	return filtered
}
