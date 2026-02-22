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

func TestSpecFromKeyWithScheme(t *testing.T) {
	tests := []struct {
		name    string
		repoKey string
		isSSH   bool
		want    string
	}{
		{
			name:    "ssh without port",
			repoKey: "github.com/org/repo",
			isSSH:   true,
			want:    "git@github.com:org/repo.git",
		},
		{
			name:    "https without port",
			repoKey: "github.com/org/repo",
			isSSH:   false,
			want:    "https://github.com/org/repo.git",
		},
		{
			name:    "ssh with port",
			repoKey: "gitlab.example.com:2222/org/repo",
			isSSH:   true,
			want:    "ssh://git@gitlab.example.com:2222/org/repo.git",
		},
		{
			name:    "https with port",
			repoKey: "gitlab.example.com:2222/org/repo",
			isSSH:   false,
			want:    "https://gitlab.example.com:2222/org/repo.git",
		},
		{
			name:    "ssh with port nested group",
			repoKey: "gitlab.example.com:2222/org/subgroup/repo",
			isSSH:   true,
			want:    "ssh://git@gitlab.example.com:2222/org/subgroup/repo.git",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SpecFromKeyWithScheme(tt.repoKey, tt.isSSH)
			if got != tt.want {
				t.Errorf("SpecFromKeyWithScheme() = %q, want %q", got, tt.want)
			}
		})
	}
}
