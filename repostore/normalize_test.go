package repostore

import (
	"context"
	"path/filepath"
	"testing"
)

type fakeNormalizerGit struct {
	localDefaultBranch      string
	remoteDefaultBranch     string
	worktreeBranches        []string
	headRefs                []string
	remoteCalls             int
	fetchCalls              int
	setRemoteHeadCalls      int
	deletedRefs             []string
	configuredRemoteFetches int
}

func (f *fakeNormalizerGit) ConfigureRemoteFetch(ctx context.Context, storePath string) error {
	f.configuredRemoteFetches++
	return nil
}
func (f *fakeNormalizerGit) LocalDefaultBranch(ctx context.Context, storePath string) (string, error) {
	return f.localDefaultBranch, nil
}
func (f *fakeNormalizerGit) DefaultBranchFromRemote(ctx context.Context, storePath string) (string, error) {
	f.remoteCalls++
	return f.remoteDefaultBranch, nil
}
func (f *fakeNormalizerGit) SetRemoteHead(ctx context.Context, storePath, branch string) error {
	f.setRemoteHeadCalls++
	return nil
}
func (f *fakeNormalizerGit) FetchPrune(ctx context.Context, storePath string, log bool) error {
	f.fetchCalls++
	return nil
}
func (f *fakeNormalizerGit) WorktreeBranches(ctx context.Context, storePath string) ([]string, error) {
	return f.worktreeBranches, nil
}
func (f *fakeNormalizerGit) HeadRefs(ctx context.Context, storePath string) ([]string, error) {
	return f.headRefs, nil
}
func (f *fakeNormalizerGit) DeleteRef(ctx context.Context, storePath, ref string) error {
	f.deletedRefs = append(f.deletedRefs, ref)
	return nil
}
func (f *fakeNormalizerGit) TouchFetchHead(storePath string) error { return nil }

func TestEnsureDefaultBranch_RemoteWhenNoCache(t *testing.T) {
	storePath := t.TempDir()
	git := &fakeNormalizerGit{remoteDefaultBranch: "main"}
	branch, err := EnsureDefaultBranch(context.Background(), git, storePath, false, "30", true)
	if err != nil {
		t.Fatalf("EnsureDefaultBranch() error = %v", err)
	}
	if branch != "main" {
		t.Fatalf("branch = %q, want main", branch)
	}
	if git.remoteCalls != 1 {
		t.Fatalf("remoteCalls = %d, want 1", git.remoteCalls)
	}
}

func TestNormalizeStore_PrunesLocalHeads(t *testing.T) {
	storePath := t.TempDir()
	if err := TouchFetchHead(storePath); err != nil {
		t.Fatalf("TouchFetchHead() error = %v", err)
	}
	git := &fakeNormalizerGit{
		localDefaultBranch: "main",
		worktreeBranches:   []string{"feat-2"},
		headRefs:           []string{"refs/heads/main", "refs/heads/feat-1", "refs/heads/feat-2"},
	}
	branch, err := NormalizeStore(context.Background(), git, filepath.Clean(storePath), false, "30", true)
	if err != nil {
		t.Fatalf("NormalizeStore() error = %v", err)
	}
	if branch != "main" {
		t.Fatalf("branch = %q, want main", branch)
	}
	if len(git.deletedRefs) != 1 || git.deletedRefs[0] != "refs/heads/feat-1" {
		t.Fatalf("deletedRefs = %#v", git.deletedRefs)
	}
}
