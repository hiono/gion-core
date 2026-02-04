package gitparse

import (
	"slices"
	"testing"
)

func TestParseWorktreeBranchNames(t *testing.T) {
	out := "worktree /tmp/w1\nHEAD abc\nbranch refs/heads/main\n\nworktree /tmp/w2\nbranch refs/heads/feat-1\n"
	got := ParseWorktreeBranchNames(out)
	slices.Sort(got)
	want := []string{"feat-1", "main"}
	if !slices.Equal(got, want) {
		t.Fatalf("ParseWorktreeBranchNames() = %#v, want %#v", got, want)
	}
}
