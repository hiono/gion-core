package applyplan

import (
	"sort"

	"github.com/tasuku43/gion-core/planner"
)

func CollectPrefetchRepoKeys(changes []planner.WorkspaceChange, desired planner.Inventory) []string {
	unique := map[string]struct{}{}
	for _, change := range changes {
		switch change.Kind {
		case planner.WorkspaceAdd:
			ws, ok := desired.Workspaces[change.WorkspaceID]
			if !ok {
				continue
			}
			for _, repoEntry := range ws.Repos {
				if repoEntry.RepoKey == "" {
					continue
				}
				unique[repoEntry.RepoKey] = struct{}{}
			}
		case planner.WorkspaceUpdate:
			for _, repoChange := range change.Repos {
				switch repoChange.Kind {
				case planner.RepoAdd, planner.RepoUpdate:
					if repoChange.ToRepo == "" {
						continue
					}
					unique[repoChange.ToRepo] = struct{}{}
				}
			}
		}
	}
	keys := make([]string, 0, len(unique))
	for repoKey := range unique {
		keys = append(keys, repoKey)
	}
	sort.Strings(keys)
	return keys
}
