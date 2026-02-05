package planner

import (
	"fmt"
	"sort"
	"strings"
)

type Repo struct {
	Alias   string
	RepoKey string
	Branch  string
}

type Workspace struct {
	ID    string
	Repos []Repo
}

type Inventory struct {
	Workspaces map[string]Workspace
}

type RepoChangeKind string

const (
	RepoAdd    RepoChangeKind = "add"
	RepoRemove RepoChangeKind = "remove"
	RepoUpdate RepoChangeKind = "update"
)

type RepoChange struct {
	Kind       RepoChangeKind
	Alias      string
	FromRepo   string
	ToRepo     string
	FromBranch string
	ToBranch   string
}

type WorkspaceChangeKind string

const (
	WorkspaceAdd    WorkspaceChangeKind = "add"
	WorkspaceRemove WorkspaceChangeKind = "remove"
	WorkspaceUpdate WorkspaceChangeKind = "update"
)

type WorkspaceChange struct {
	Kind        WorkspaceChangeKind
	WorkspaceID string
	Repos       []RepoChange
}

func Diff(desired, actual Inventory) []WorkspaceChange {
	var changes []WorkspaceChange
	desiredIDs := sortedKeys(desired.Workspaces)
	actualIDs := sortedKeys(actual.Workspaces)

	actualSet := make(map[string]Workspace, len(actual.Workspaces))
	for id, ws := range actual.Workspaces {
		actualSet[id] = ws
	}

	for _, id := range desiredIDs {
		desiredWS := desired.Workspaces[id]
		actualWS, exists := actualSet[id]
		if !exists {
			changes = append(changes, WorkspaceChange{
				Kind:        WorkspaceAdd,
				WorkspaceID: id,
				Repos:       plannedRepoAdds(desiredWS),
			})
			continue
		}
		repoChanges := diffRepos(actualWS.Repos, desiredWS.Repos)
		if len(repoChanges) > 0 {
			changes = append(changes, WorkspaceChange{
				Kind:        WorkspaceUpdate,
				WorkspaceID: id,
				Repos:       repoChanges,
			})
		}
	}

	desiredSet := make(map[string]Workspace, len(desired.Workspaces))
	for id, ws := range desired.Workspaces {
		desiredSet[id] = ws
	}
	for _, id := range actualIDs {
		if _, exists := desiredSet[id]; exists {
			continue
		}
		changes = append(changes, WorkspaceChange{
			Kind:        WorkspaceRemove,
			WorkspaceID: id,
		})
	}

	sort.Slice(changes, func(i, j int) bool {
		if changes[i].WorkspaceID == changes[j].WorkspaceID {
			return changes[i].Kind < changes[j].Kind
		}
		return changes[i].WorkspaceID < changes[j].WorkspaceID
	})

	return changes
}

func diffRepos(actualRepos, desiredRepos []Repo) []RepoChange {
	actualByAlias := map[string]Repo{}
	for _, repo := range actualRepos {
		actualByAlias[strings.TrimSpace(repo.Alias)] = repo
	}
	desiredByAlias := map[string]Repo{}
	for _, repo := range desiredRepos {
		desiredByAlias[strings.TrimSpace(repo.Alias)] = repo
	}

	var changes []RepoChange
	for alias, desired := range desiredByAlias {
		actual, exists := actualByAlias[alias]
		if !exists {
			changes = append(changes, RepoChange{
				Kind:     RepoAdd,
				Alias:    alias,
				ToRepo:   desired.RepoKey,
				ToBranch: desired.Branch,
			})
			continue
		}
		if actual.RepoKey != desired.RepoKey || actual.Branch != desired.Branch {
			changes = append(changes, RepoChange{
				Kind:       RepoUpdate,
				Alias:      alias,
				FromRepo:   actual.RepoKey,
				ToRepo:     desired.RepoKey,
				FromBranch: actual.Branch,
				ToBranch:   desired.Branch,
			})
		}
	}

	for alias, actual := range actualByAlias {
		if _, exists := desiredByAlias[alias]; exists {
			continue
		}
		changes = append(changes, RepoChange{
			Kind:       RepoRemove,
			Alias:      alias,
			FromRepo:   actual.RepoKey,
			FromBranch: actual.Branch,
		})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Alias < changes[j].Alias
	})
	return changes
}

func plannedRepoAdds(ws Workspace) []RepoChange {
	var changes []RepoChange
	for _, repo := range ws.Repos {
		changes = append(changes, RepoChange{
			Kind:     RepoAdd,
			Alias:    strings.TrimSpace(repo.Alias),
			ToRepo:   repo.RepoKey,
			ToBranch: repo.Branch,
		})
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Alias < changes[j].Alias
	})
	return changes
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (k WorkspaceChangeKind) String() string {
	switch k {
	case WorkspaceAdd:
		return "add"
	case WorkspaceRemove:
		return "remove"
	case WorkspaceUpdate:
		return "update"
	default:
		return fmt.Sprintf("unknown(%s)", string(k))
	}
}
