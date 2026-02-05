package gcplan

import "testing"

func TestBuildFetchPlan(t *testing.T) {
	plan := BuildFetchPlan([]RepoEntry{
		{RepoKey: "example.com/org/repo", BaseRef: "origin/main"},
		{RepoKey: "example.com/org/repo", BaseRef: "refs/remotes/upstream/release"},
		{RepoKey: "example.com/org/repo", BaseRef: "origin/main"},
		{RepoKey: "example.com/org/repo", BaseRef: ""},
	})
	if !plan.NeedsDefault {
		t.Fatalf("NeedsDefault = false, want true")
	}
	if got := len(plan.Targets); got != 2 {
		t.Fatalf("len(Targets) = %d, want 2", got)
	}
	if plan.Targets[0].Remote != "origin" || plan.Targets[0].Branch != "main" {
		t.Fatalf("Targets[0] = %#v", plan.Targets[0])
	}
	if plan.Targets[1].Remote != "upstream" || plan.Targets[1].Branch != "release" {
		t.Fatalf("Targets[1] = %#v", plan.Targets[1])
	}
}
