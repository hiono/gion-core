package repospec

import "testing"

func TestDisplaySpec(t *testing.T) {
	got := DisplaySpec("https://github.com/org/repo")
	if got != "git@github.com:org/repo.git" {
		t.Fatalf("DisplaySpec() = %q", got)
	}
}

func TestDisplayName(t *testing.T) {
	got := DisplayName("https://github.com/org/repo")
	if got != "repo" {
		t.Fatalf("DisplayName() = %q", got)
	}
}

func TestSpecFromKey(t *testing.T) {
	got := SpecFromKey("github.com/org/repo")
	if got != "git@github.com:org/repo.git" {
		t.Fatalf("SpecFromKey() = %q", got)
	}
}
