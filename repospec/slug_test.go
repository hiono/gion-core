package repospec

import (
	"reflect"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{
			name:  "single part",
			parts: []string{"repo"},
			want:  "repo",
		},
		{
			name:  "multiple parts",
			parts: []string{"owner", "repo"},
			want:  "owner-repo",
		},
		{
			name:  "with hyphen in part",
			parts: []string{"my-repo"},
			want:  "my--repo",
		},
		{
			name:  "multiple with hyphens",
			parts: []string{"my-repo", "sub-group"},
			want:  "my--repo-sub--group",
		},
		{
			name:  "empty parts",
			parts: []string{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.parts); got != tt.want {
				t.Errorf("Slugify(%v) = %q, want %q", tt.parts, got, tt.want)
			}
		})
	}
}

func TestUnslugify(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want []string
	}{
		{
			name: "single part",
			slug: "repo",
			want: []string{"repo"},
		},
		{
			name: "multiple parts",
			slug: "owner-repo",
			want: []string{"owner", "repo"},
		},
		{
			name: "unescape hyphen",
			slug: "my--repo",
			want: []string{"my-repo"},
		},
		{
			name: "multiple escaped hyphens",
			slug: "my--repo-sub--group",
			want: []string{"my-repo", "sub-group"},
		},
		{
			name: "empty slug",
			slug: "",
			want: []string{},
		},
		{
			name: "no escaping needed",
			slug: "abc",
			want: []string{"abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unslugify(tt.slug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unslugify(%q) = %v, want %v", tt.slug, got, tt.want)
			}
		})
	}
}

func TestSlugifyUnslugifyRoundtrip(t *testing.T) {
	testCases := [][]string{
		{"repo"},
		{"owner", "repo"},
		{"my-repo"},
		{"owner", "my-repo", "sub-group"},
		{"a", "b", "c"},
	}

	for _, parts := range testCases {
		slug := Slugify(parts)
		unslugified := Unslugify(slug)
		if !reflect.DeepEqual(unslugified, parts) {
			t.Errorf("Roundtrip failed for %v: got %v", parts, unslugified)
		}
	}
}
