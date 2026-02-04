package repostore

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	RepoKey   string
	StorePath string
	RemoteURL string
}

type EnsureStoreRequest struct {
	RepoKey       string
	RemoteURL     string
	StorePath     string
	RepoSpec      string
	MustExist     bool
	Fetch         bool
	FetchGraceEnv string
	Log           bool
}

type EnsureStoreResult struct {
	Store   Store
	Created bool
}

type StoreAccess interface {
	DirExists(path string) (bool, error)
	MkdirAll(path string, perm os.FileMode) error
	CloneBare(ctx context.Context, remoteURL, storePath string) error
	NormalizeStore(ctx context.Context, storePath string, fetch bool, fetchGraceEnv string, log bool) error
}

type ErrStoreNotFound struct {
	RepoSpec string
}

func (e ErrStoreNotFound) Error() string {
	spec := strings.TrimSpace(e.RepoSpec)
	if spec == "" {
		return "repo store not found"
	}
	return fmt.Sprintf("repo store not found, run: gion repo get %s", spec)
}

func EnsureStore(ctx context.Context, access StoreAccess, req EnsureStoreRequest) (EnsureStoreResult, error) {
	storePath := strings.TrimSpace(req.StorePath)
	if storePath == "" {
		return EnsureStoreResult{}, fmt.Errorf("store path is required")
	}

	exists, err := access.DirExists(storePath)
	if err != nil {
		return EnsureStoreResult{}, err
	}
	if req.MustExist && !exists {
		return EnsureStoreResult{}, ErrStoreNotFound{RepoSpec: req.RepoSpec}
	}

	created := false
	if !exists {
		if err := access.MkdirAll(filepath.Dir(storePath), 0o750); err != nil {
			return EnsureStoreResult{}, fmt.Errorf("create repo store dir: %w", err)
		}
		if err := access.CloneBare(ctx, req.RemoteURL, storePath); err != nil {
			return EnsureStoreResult{}, err
		}
		created = true
	}

	if err := access.NormalizeStore(ctx, storePath, req.Fetch, req.FetchGraceEnv, req.Log); err != nil {
		return EnsureStoreResult{}, err
	}

	return EnsureStoreResult{
		Store: Store{
			RepoKey:   strings.TrimSpace(req.RepoKey),
			StorePath: storePath,
			RemoteURL: strings.TrimSpace(req.RemoteURL),
		},
		Created: created,
	}, nil
}
