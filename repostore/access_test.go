package repostore

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type fakeStoreAccess struct {
	exists             bool
	dirExistsErr       error
	mkdirErr           error
	cloneErr           error
	normalizeErr       error
	mkdirCalled        int
	cloneCalled        int
	normalizeCalled    int
	lastNormalizeFetch bool
}

func (f *fakeStoreAccess) DirExists(path string) (bool, error) {
	return f.exists, f.dirExistsErr
}

func (f *fakeStoreAccess) MkdirAll(path string, perm os.FileMode) error {
	f.mkdirCalled++
	return f.mkdirErr
}

func (f *fakeStoreAccess) CloneBare(ctx context.Context, remoteURL, storePath string) error {
	f.cloneCalled++
	return f.cloneErr
}

func (f *fakeStoreAccess) NormalizeStore(ctx context.Context, storePath string, fetch bool, fetchGraceEnv string, log bool) error {
	f.normalizeCalled++
	f.lastNormalizeFetch = fetch
	return f.normalizeErr
}

func TestEnsureStore_CreateAndNormalize(t *testing.T) {
	access := &fakeStoreAccess{exists: false}
	req := EnsureStoreRequest{
		RepoKey:   "github.com/org/repo",
		RemoteURL: "git@github.com:org/repo.git",
		StorePath: filepath.Join(t.TempDir(), "repo.git"),
	}
	result, err := EnsureStore(context.Background(), access, req)
	if err != nil {
		t.Fatalf("EnsureStore() error = %v", err)
	}
	if !result.Created {
		t.Fatalf("Created = false, want true")
	}
	if access.mkdirCalled != 1 || access.cloneCalled != 1 || access.normalizeCalled != 1 {
		t.Fatalf("calls = mkdir:%d clone:%d normalize:%d", access.mkdirCalled, access.cloneCalled, access.normalizeCalled)
	}
}

func TestEnsureStore_MustExist(t *testing.T) {
	access := &fakeStoreAccess{exists: false}
	_, err := EnsureStore(context.Background(), access, EnsureStoreRequest{
		StorePath: filepath.Join(t.TempDir(), "repo.git"),
		RepoSpec:  "git@github.com:org/repo.git",
		MustExist: true,
	})
	var notFound ErrStoreNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("error = %v, want ErrStoreNotFound", err)
	}
}

func TestEnsureStore_NormalizeError(t *testing.T) {
	access := &fakeStoreAccess{exists: true, normalizeErr: errors.New("normalize failed")}
	_, err := EnsureStore(context.Background(), access, EnsureStoreRequest{
		StorePath: filepath.Join(t.TempDir(), "repo.git"),
		Fetch:     true,
	})
	if err == nil || err.Error() != "normalize failed" {
		t.Fatalf("error = %v, want normalize failed", err)
	}
	if !access.lastNormalizeFetch {
		t.Fatalf("fetch flag was not propagated")
	}
}
