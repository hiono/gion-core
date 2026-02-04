package importplan

import "sort"

type RepoSnapshot struct {
	Alias   string
	RepoKey string
	Branch  string
}

type WorkspaceSnapshot struct {
	ID          string
	Description string
	Mode        string
	PresetName  string
	SourceURL   string
	BaseBranch  string
	Repos       []RepoSnapshot
}

type Repo struct {
	Alias   string
	RepoKey string
	Branch  string
	BaseRef string
}

type Workspace struct {
	Description string
	Mode        string
	PresetName  string
	SourceURL   string
	Repos       []Repo
}

type Inventory struct {
	Workspaces map[string]Workspace
}

func BuildInventory(snapshots []WorkspaceSnapshot) Inventory {
	workspaces := make(map[string]Workspace, len(snapshots))
	for _, snapshot := range snapshots {
		repos := make([]Repo, 0, len(snapshot.Repos))
		for _, repo := range snapshot.Repos {
			repos = append(repos, Repo{
				Alias:   repo.Alias,
				RepoKey: repo.RepoKey,
				Branch:  repo.Branch,
				BaseRef: snapshot.BaseBranch,
			})
		}
		sort.Slice(repos, func(i, j int) bool {
			return repos[i].Alias < repos[j].Alias
		})
		workspaces[snapshot.ID] = Workspace{
			Description: snapshot.Description,
			Mode:        snapshot.Mode,
			PresetName:  snapshot.PresetName,
			SourceURL:   snapshot.SourceURL,
			Repos:       repos,
		}
	}
	return Inventory{Workspaces: workspaces}
}

func CollectWorkspaceIDs(names []string) []string {
	ids := make([]string, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		ids = append(ids, name)
	}
	sort.Strings(ids)
	return ids
}
