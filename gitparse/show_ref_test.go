package gitparse

import (
	"reflect"
	"testing"
)

func TestParseHeadRefs(t *testing.T) {
	out := "aaaa refs/heads/main\n\ncccc refs/tags/v1\nbbbb refs/heads/feat-1\ninvalid\n"
	got := ParseHeadRefs(out)
	want := []string{"refs/heads/main", "refs/heads/feat-1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseHeadRefs() = %#v, want %#v", got, want)
	}
}
