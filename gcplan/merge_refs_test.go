package gcplan

import "testing"

func TestMergeCheckRefs(t *testing.T) {
	head, target := MergeCheckRefs("WS-1", "origin/main")
	if head != "refs/heads/WS-1" {
		t.Fatalf("head = %q", head)
	}
	if target != "refs/remotes/origin/main" {
		t.Fatalf("target = %q", target)
	}
}
