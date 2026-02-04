package gcplan

import (
	"context"
	"errors"
	"testing"
)

type fakeRepoFetcher struct {
	defaultBranch    string
	defaultBranchErr error
	fetchErr         error
	fetchCalls       []FetchTarget
}

func (f *fakeRepoFetcher) DefaultBranchFromRemote(_ context.Context, _ string) (string, error) {
	if f.defaultBranchErr != nil {
		return "", f.defaultBranchErr
	}
	return f.defaultBranch, nil
}

func (f *fakeRepoFetcher) FetchRemoteBranch(_ context.Context, _ string, remote, branch string) error {
	if f.fetchErr != nil {
		return f.fetchErr
	}
	f.fetchCalls = append(f.fetchCalls, FetchTarget{Remote: remote, Branch: branch})
	return nil
}

func TestFetchRepo(t *testing.T) {
	t.Run("base refs only", func(t *testing.T) {
		fetcher := &fakeRepoFetcher{}
		result, err := FetchRepo(context.Background(), fetcher, "/tmp/store", []RepoEntry{
			{RepoKey: "owner/a", BaseRef: "origin/main"},
			{RepoKey: "owner/a", BaseRef: "origin/release"},
			{RepoKey: "owner/a", BaseRef: "origin/main"},
		})
		if err != nil {
			t.Fatalf("FetchRepo() error = %v", err)
		}
		if result.DefaultTarget != "" {
			t.Fatalf("DefaultTarget = %q, want empty", result.DefaultTarget)
		}
		want := []FetchTarget{
			{Remote: "origin", Branch: "main"},
			{Remote: "origin", Branch: "release"},
		}
		assertFetchCalls(t, fetcher.fetchCalls, want)
	})

	t.Run("needs default branch", func(t *testing.T) {
		fetcher := &fakeRepoFetcher{defaultBranch: "develop"}
		result, err := FetchRepo(context.Background(), fetcher, "/tmp/store", []RepoEntry{
			{RepoKey: "owner/a", BaseRef: ""},
		})
		if err != nil {
			t.Fatalf("FetchRepo() error = %v", err)
		}
		if result.DefaultTarget != "origin/develop" {
			t.Fatalf("DefaultTarget = %q, want origin/develop", result.DefaultTarget)
		}
		want := []FetchTarget{
			{Remote: "origin", Branch: "develop"},
		}
		assertFetchCalls(t, fetcher.fetchCalls, want)
	})

	t.Run("default branch error", func(t *testing.T) {
		fetcher := &fakeRepoFetcher{defaultBranchErr: errors.New("boom")}
		_, err := FetchRepo(context.Background(), fetcher, "/tmp/store", []RepoEntry{
			{RepoKey: "owner/a", BaseRef: ""},
		})
		if err == nil || err.Error() != "boom" {
			t.Fatalf("FetchRepo() error = %v, want boom", err)
		}
	})

	t.Run("default branch missing", func(t *testing.T) {
		fetcher := &fakeRepoFetcher{defaultBranch: " "}
		_, err := FetchRepo(context.Background(), fetcher, "/tmp/store", []RepoEntry{
			{RepoKey: "owner/a", BaseRef: ""},
		})
		if err == nil || err.Error() != "default branch unavailable" {
			t.Fatalf("FetchRepo() error = %v, want default branch unavailable", err)
		}
	})

	t.Run("fetch error", func(t *testing.T) {
		fetcher := &fakeRepoFetcher{fetchErr: errors.New("fetch failed")}
		_, err := FetchRepo(context.Background(), fetcher, "/tmp/store", []RepoEntry{
			{RepoKey: "owner/a", BaseRef: "origin/main"},
		})
		if err == nil || err.Error() != "fetch failed" {
			t.Fatalf("FetchRepo() error = %v, want fetch failed", err)
		}
	})
}

func assertFetchCalls(t *testing.T, got []FetchTarget, want []FetchTarget) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("fetch calls len = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("fetch calls[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}
