package gcplan

import (
	"context"
	"fmt"
	"strings"
)

type MergeChecker interface {
	ShowRef(ctx context.Context, storePath, ref string) (hash string, ok bool, err error)
	IsAncestor(ctx context.Context, storePath, headRef, targetRef string) (bool, error)
}

type StrictMergeRequest struct {
	StorePath string
	Branch    string
	Target    string
}

func StrictMergedIntoTarget(ctx context.Context, checker MergeChecker, req StrictMergeRequest) (bool, error) {
	branch := strings.TrimSpace(req.Branch)
	if branch == "" {
		return false, fmt.Errorf("branch is required")
	}
	target := strings.TrimSpace(req.Target)
	if target == "" {
		return false, fmt.Errorf("target is required")
	}

	headRef, targetRef := MergeCheckRefs(branch, target)

	headHash, headOK, err := checker.ShowRef(ctx, req.StorePath, headRef)
	if err != nil {
		return false, err
	}
	if !headOK {
		return false, fmt.Errorf("ref not found: %s", headRef)
	}
	targetHash, targetOK, err := checker.ShowRef(ctx, req.StorePath, targetRef)
	if err != nil {
		return false, err
	}
	if !targetOK {
		return false, fmt.Errorf("ref not found: %s", targetRef)
	}
	if headHash == targetHash {
		return false, nil
	}

	ok, err := checker.IsAncestor(ctx, req.StorePath, headRef, targetRef)
	if err != nil {
		return false, err
	}
	return ok, nil
}
