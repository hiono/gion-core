package gitparse

import "testing"

func TestParseRemoteHeadSymref(t *testing.T) {
	out := "ref: refs/heads/main\tHEAD\n1234567890abcdef\tHEAD\n"
	branch, hash := ParseRemoteHeadSymref(out)
	if branch != "main" {
		t.Fatalf("branch = %q, want main", branch)
	}
	if hash != "1234567890abcdef" {
		t.Fatalf("hash = %q", hash)
	}
}

func TestParseRemoteHeadSymref_NoBranch(t *testing.T) {
	out := "1234567890abcdef\tHEAD\n"
	branch, hash := ParseRemoteHeadSymref(out)
	if branch != "" {
		t.Fatalf("branch = %q, want empty", branch)
	}
	if hash != "1234567890abcdef" {
		t.Fatalf("hash = %q", hash)
	}
}
