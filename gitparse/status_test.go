package gitparse

import "testing"

func TestParseStatusPorcelainV2(t *testing.T) {
	out := "# branch.oid 94a67ef\n# branch.head main\n# branch.upstream origin/main\n# branch.ab +2 -1\n1 .M N... 100644 100644 100644 abcdef0 abcdef0 file.txt\n? new.txt\nu UU N... 100644 100644 100644 abcdef0 abcdef0 abcdef0 conflict.txt\n"
	got := ParseStatusPorcelainV2(out, "fallback")

	if got.Branch != "main" {
		t.Fatalf("Branch = %q, want main", got.Branch)
	}
	if got.Upstream != "origin/main" {
		t.Fatalf("Upstream = %q, want origin/main", got.Upstream)
	}
	if got.Head != "94a67ef" {
		t.Fatalf("Head = %q, want 94a67ef", got.Head)
	}
	if got.Detached {
		t.Fatalf("Detached = true, want false")
	}
	if got.HeadMissing {
		t.Fatalf("HeadMissing = true, want false")
	}
	if !got.Dirty {
		t.Fatalf("Dirty = false, want true")
	}
	if got.UntrackedCount != 1 {
		t.Fatalf("UntrackedCount = %d, want 1", got.UntrackedCount)
	}
	if got.StagedCount != 0 {
		t.Fatalf("StagedCount = %d, want 0", got.StagedCount)
	}
	if got.UnstagedCount != 1 {
		t.Fatalf("UnstagedCount = %d, want 1", got.UnstagedCount)
	}
	if got.UnmergedCount != 1 {
		t.Fatalf("UnmergedCount = %d, want 1", got.UnmergedCount)
	}
	if got.AheadCount != 2 {
		t.Fatalf("AheadCount = %d, want 2", got.AheadCount)
	}
	if got.BehindCount != 1 {
		t.Fatalf("BehindCount = %d, want 1", got.BehindCount)
	}
}

func TestParseChangedFilesPorcelainV2(t *testing.T) {
	out := "# branch.oid 94a67ef\n# branch.head main\n# branch.upstream origin/main\n# branch.ab +2 -1\n1 .M N... 100644 100644 100644 abcdef0 abcdef0 file.txt\n? new.txt\nu UU N... 100644 100644 100644 abcdef0 abcdef0 abcdef0 conflict.txt\n2 R. N... 100644 100644 100644 abcdef0 abcdef0 R100 old.txt new.txt\n"
	got := ParseChangedFilesPorcelainV2(out)
	want := []string{
		" M file.txt",
		"?? new.txt",
		"UU conflict.txt",
		"R  old.txt -> new.txt",
	}
	if len(got) != len(want) {
		t.Fatalf("len(files) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("files[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
