package workspacerisk

import "strings"

type RepoStatus struct {
	Upstream    string
	AheadCount  int
	BehindCount int
	Dirty       bool
	Detached    bool
	HeadMissing bool
	Error       error
}

func ClassifyRepoStatus(status RepoStatus) RepoState {
	if status.Error != nil {
		return RepoStateUnknown
	}
	if status.Dirty {
		return RepoStateDirty
	}
	if status.Detached || status.HeadMissing {
		return RepoStateUnknown
	}
	if strings.TrimSpace(status.Upstream) == "" {
		return RepoStateUnknown
	}
	if status.AheadCount > 0 && status.BehindCount > 0 {
		return RepoStateDiverged
	}
	if status.AheadCount > 0 {
		return RepoStateUnpushed
	}
	return RepoStateClean
}
