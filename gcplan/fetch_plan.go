package gcplan

import (
	"sort"

	"github.com/tasuku43/gion-core/gitref"
)

type RepoEntry struct {
	RepoKey string
	BaseRef string
}

type FetchTarget struct {
	Remote string
	Branch string
}

type FetchPlan struct {
	NeedsDefault bool
	Targets      []FetchTarget
}

func BuildFetchPlan(entries []RepoEntry) FetchPlan {
	targetSet := map[FetchTarget]struct{}{}
	needsDefault := false
	for _, entry := range entries {
		target, ok := gitref.ParseRemoteBranch(entry.BaseRef)
		if !ok {
			needsDefault = true
			continue
		}
		targetSet[FetchTarget{Remote: target.Remote, Branch: target.Branch}] = struct{}{}
	}
	targets := make([]FetchTarget, 0, len(targetSet))
	for target := range targetSet {
		targets = append(targets, target)
	}
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].Remote == targets[j].Remote {
			return targets[i].Branch < targets[j].Branch
		}
		return targets[i].Remote < targets[j].Remote
	})
	return FetchPlan{NeedsDefault: needsDefault, Targets: targets}
}
