package repostore

import "context"

type NormalizerGit interface {
	ConfigureRemoteFetch(ctx context.Context, storePath string) error
	LocalDefaultBranch(ctx context.Context, storePath string) (string, error)
	DefaultBranchFromRemote(ctx context.Context, storePath string) (string, error)
	SetRemoteHead(ctx context.Context, storePath, branch string) error
	FetchPrune(ctx context.Context, storePath string, log bool) error
	WorktreeBranches(ctx context.Context, storePath string) ([]string, error)
	HeadRefs(ctx context.Context, storePath string) ([]string, error)
	DeleteRef(ctx context.Context, storePath, ref string) error
	TouchFetchHead(storePath string) error
}

func EnsureDefaultBranch(ctx context.Context, git NormalizerGit, storePath string, fetch bool, fetchGraceEnv string, log bool) (string, error) {
	if err := git.ConfigureRemoteFetch(ctx, storePath); err != nil {
		return "", err
	}
	defaultBranch, _ := git.LocalDefaultBranch(ctx, storePath)

	remoteChecked := false
	grace := FetchGraceDuration(fetchGraceEnv)
	if !RecentlyFetched(storePath, grace) || defaultBranch == "" {
		branch, err := git.DefaultBranchFromRemote(ctx, storePath)
		if err != nil {
			return "", err
		}
		defaultBranch = branch
		remoteChecked = true
	}
	if defaultBranch != "" {
		_ = git.SetRemoteHead(ctx, storePath, defaultBranch)
	}

	if fetch {
		if err := git.FetchPrune(ctx, storePath, log); err != nil {
			return "", err
		}
	} else if remoteChecked {
		if err := git.TouchFetchHead(storePath); err != nil {
			return "", err
		}
	}

	return defaultBranch, nil
}

func NormalizeStore(ctx context.Context, git NormalizerGit, storePath string, fetch bool, fetchGraceEnv string, log bool) (string, error) {
	defaultBranch, err := EnsureDefaultBranch(ctx, git, storePath, fetch, fetchGraceEnv, log)
	if err != nil {
		return "", err
	}
	worktreeBranches, _ := git.WorktreeBranches(ctx, storePath)
	headRefs, err := git.HeadRefs(ctx, storePath)
	if err != nil {
		return "", err
	}
	prunable := SelectPrunableHeadRefs(headRefs, defaultBranch, worktreeBranches)
	for _, ref := range prunable {
		_ = git.DeleteRef(ctx, storePath, ref)
	}
	return defaultBranch, nil
}
