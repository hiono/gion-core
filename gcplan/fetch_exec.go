package gcplan

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type RepoFetcher interface {
	DefaultBranchFromRemote(ctx context.Context, storePath string) (string, error)
	FetchRemoteBranch(ctx context.Context, storePath, remote, branch string) error
}

type FetchResult struct {
	DefaultTarget string
}

func FetchRepo(ctx context.Context, fetcher RepoFetcher, storePath string, entries []RepoEntry) (FetchResult, error) {
	plan := BuildFetchPlan(entries)
	targetSet := make(map[FetchTarget]struct{}, len(plan.Targets)+1)
	for _, target := range plan.Targets {
		targetSet[target] = struct{}{}
	}

	defaultTarget := ""
	if plan.NeedsDefault {
		branch, err := fetcher.DefaultBranchFromRemote(ctx, storePath)
		if err != nil {
			return FetchResult{}, err
		}
		branch = strings.TrimSpace(branch)
		if branch == "" {
			return FetchResult{}, fmt.Errorf("default branch unavailable")
		}
		defaultTarget = fmt.Sprintf("origin/%s", branch)
		targetSet[FetchTarget{Remote: "origin", Branch: branch}] = struct{}{}
	}

	targets := sortedFetchTargets(targetSet)
	for _, target := range targets {
		if err := fetcher.FetchRemoteBranch(ctx, storePath, target.Remote, target.Branch); err != nil {
			return FetchResult{}, err
		}
	}
	return FetchResult{DefaultTarget: defaultTarget}, nil
}

func sortedFetchTargets(targetSet map[FetchTarget]struct{}) []FetchTarget {
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
	return targets
}
