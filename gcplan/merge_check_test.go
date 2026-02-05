package gcplan

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type fakeMergeChecker struct {
	showRefErr   map[string]error
	refs         map[string]string
	ancestor     map[string]bool
	ancestorErr  error
	ancestorCall int
}

func (f *fakeMergeChecker) ShowRef(_ context.Context, _ string, ref string) (string, bool, error) {
	if err := f.showRefErr[ref]; err != nil {
		return "", false, err
	}
	hash, ok := f.refs[ref]
	return hash, ok, nil
}

func (f *fakeMergeChecker) IsAncestor(_ context.Context, _ string, headRef, targetRef string) (bool, error) {
	f.ancestorCall++
	if f.ancestorErr != nil {
		return false, f.ancestorErr
	}
	return f.ancestor[fmt.Sprintf("%s->%s", headRef, targetRef)], nil
}

func TestStrictMergedIntoTarget(t *testing.T) {
	tests := []struct {
		name          string
		req           StrictMergeRequest
		checker       *fakeMergeChecker
		want          bool
		wantErr       string
		wantAncestors int
	}{
		{
			name:    "branch required",
			req:     StrictMergeRequest{Target: "origin/main"},
			checker: &fakeMergeChecker{},
			wantErr: "branch is required",
		},
		{
			name:    "target required",
			req:     StrictMergeRequest{Branch: "feat"},
			checker: &fakeMergeChecker{},
			wantErr: "target is required",
		},
		{
			name: "head ref missing",
			req:  StrictMergeRequest{Branch: "feat", Target: "origin/main"},
			checker: &fakeMergeChecker{
				refs: map[string]string{
					"refs/remotes/origin/main": "bbb",
				},
			},
			wantErr: "ref not found: refs/heads/feat",
		},
		{
			name: "target ref missing",
			req:  StrictMergeRequest{Branch: "feat", Target: "origin/main"},
			checker: &fakeMergeChecker{
				refs: map[string]string{
					"refs/heads/feat": "aaa",
				},
			},
			wantErr: "ref not found: refs/remotes/origin/main",
		},
		{
			name: "show ref error",
			req:  StrictMergeRequest{Branch: "feat", Target: "origin/main"},
			checker: &fakeMergeChecker{
				showRefErr: map[string]error{
					"refs/heads/feat": errors.New("boom"),
				},
			},
			wantErr: "boom",
		},
		{
			name: "same hash is not merged",
			req:  StrictMergeRequest{Branch: "feat", Target: "origin/main"},
			checker: &fakeMergeChecker{
				refs: map[string]string{
					"refs/heads/feat":          "aaa",
					"refs/remotes/origin/main": "aaa",
				},
			},
			want:          false,
			wantAncestors: 0,
		},
		{
			name: "ancestor true",
			req:  StrictMergeRequest{Branch: "feat", Target: "origin/main"},
			checker: &fakeMergeChecker{
				refs: map[string]string{
					"refs/heads/feat":          "aaa",
					"refs/remotes/origin/main": "bbb",
				},
				ancestor: map[string]bool{
					"refs/heads/feat->refs/remotes/origin/main": true,
				},
			},
			want:          true,
			wantAncestors: 1,
		},
		{
			name: "ancestor false",
			req:  StrictMergeRequest{Branch: "feat", Target: "origin/main"},
			checker: &fakeMergeChecker{
				refs: map[string]string{
					"refs/heads/feat":          "aaa",
					"refs/remotes/origin/main": "bbb",
				},
				ancestor: map[string]bool{
					"refs/heads/feat->refs/remotes/origin/main": false,
				},
			},
			want:          false,
			wantAncestors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrictMergedIntoTarget(context.Background(), tt.checker, tt.req)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("StrictMergedIntoTarget() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("StrictMergedIntoTarget() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("StrictMergedIntoTarget() = %v, want %v", got, tt.want)
			}
			if tt.checker.ancestorCall != tt.wantAncestors {
				t.Fatalf("IsAncestor calls = %d, want %d", tt.checker.ancestorCall, tt.wantAncestors)
			}
		})
	}
}
