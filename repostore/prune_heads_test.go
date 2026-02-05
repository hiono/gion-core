package repostore

import (
	"reflect"
	"testing"
)

func TestSelectPrunableHeadRefs(t *testing.T) {
	refs := []string{
		"refs/heads/main",
		"refs/heads/feat-1",
		"refs/heads/feat-2",
	}
	got := SelectPrunableHeadRefs(refs, "main", []string{"feat-2"})
	want := []string{"refs/heads/feat-1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SelectPrunableHeadRefs() = %#v, want %#v", got, want)
	}
}
