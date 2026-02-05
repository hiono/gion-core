package gitref

import "testing"

func TestParseOriginHeadRef(t *testing.T) {
	branch, ok := ParseOriginHeadRef("refs/remotes/origin/main")
	if !ok || branch != "main" {
		t.Fatalf("ParseOriginHeadRef() = (%q,%v)", branch, ok)
	}
	if _, ok := ParseOriginHeadRef("refs/heads/main"); ok {
		t.Fatalf("expected false for non-origin-head ref")
	}
}

func TestFormatOriginTarget(t *testing.T) {
	target, ok := FormatOriginTarget("main")
	if !ok || target != "origin/main" {
		t.Fatalf("FormatOriginTarget() = (%q,%v)", target, ok)
	}
	if _, ok := FormatOriginTarget(" "); ok {
		t.Fatalf("expected false for empty branch")
	}
}
