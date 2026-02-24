package repospec

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantKey string
		wantErr bool
	}{
		{
			name:    "ssh",
			input:   "git@github.com:org/repo.git",
			wantKey: "github.com/org/repo",
		},
		{
			name:    "ssh with port",
			input:   "git@gitlab.example.com:2222:org/repo.git",
			wantKey: "gitlab.example.com:2222/org/repo",
		},
		{
			name:    "ssh nested group",
			input:   "git@gitlab.example.com:org/subgroup/team/repo.git",
			wantKey: "gitlab.example.com/org/subgroup-team-repo",
		},
		{
			name:    "ssh with port nested group",
			input:   "git@gitlab.example.com:2222:org/subgroup/team/repo.git",
			wantKey: "gitlab.example.com:2222/org/subgroup-team-repo",
		},
		{
			name:    "https nested group",
			input:   "https://gitlab.example.com/org/subgroup/team/repo.git",
			wantKey: "gitlab.example.com/org/subgroup-team-repo",
		},
		{
			name:    "ssh port min boundary",
			input:   "git@host.com:1:org/repo.git",
			wantKey: "host.com:1/org/repo",
		},
		{
			name:    "ssh port max boundary",
			input:   "git@host.com:65535:org/repo.git",
			wantKey: "host.com:65535/org/repo",
		},
		{
			name:    "https",
			input:   "https://github.com/org/repo",
			wantKey: "github.com/org/repo",
		},
		{
			name:    "https nested group",
			input:   "https://gitlab.example.com/org/subgroup/team/repo.git",
			wantKey: "gitlab.example.com/org/subgroup-team-repo",
		},
		{
			name:    "file",
			input:   "file:///tmp/mirrors/example.com/org/repo.git",
			wantKey: "example.com/org/repo",
		},
		{
			name:    "shorthand",
			input:   "github.com/org/repo.git",
			wantErr: true,
		},
		{
			name:    "invalid",
			input:   "org/repo",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := Normalize(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if spec.RepoKey != tc.wantKey {
				t.Fatalf("repo key mismatch: got %q want %q", spec.RepoKey, tc.wantKey)
			}
		})
	}
}

func TestNormalizeWithBasePath(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		basePath string
		wantKey  string
		wantErr  bool
	}{
		{
			name:     "https with base_path",
			input:    "https://host.com/git/org/repo.git",
			basePath: "/git",
			wantKey:  "host.com/org/repo",
		},
		{
			name:     "https nested group with base_path",
			input:    "https://host.com/git/org/team/repo.git",
			basePath: "/git",
			wantKey:  "host.com/org/team-repo",
		},
		{
			name:     "ssh ignores base_path",
			input:    "git@host.com:git/org/repo.git",
			basePath: "/git",
			wantKey:  "host.com/git/org-repo",
		},
		{
			name:     "ssh with port ignores base_path",
			input:    "git@host.com:2222:git/org/repo.git",
			basePath: "/git",
			wantKey:  "host.com:2222/git/org-repo",
		},
		{
			name:     "https no base_path specified",
			input:    "https://host.com/git/org/repo.git",
			basePath: "",
			wantKey:  "host.com/git/org-repo",
		},
		{
			name:     "https base_path not matching",
			input:    "https://host.com/org/repo.git",
			basePath: "/git",
			wantKey:  "host.com/org/repo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := NormalizeWithBasePath(tc.input, tc.basePath)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if spec.RepoKey != tc.wantKey {
				t.Fatalf("repo key mismatch: got %q want %q", spec.RepoKey, tc.wantKey)
			}
		})
	}
}
