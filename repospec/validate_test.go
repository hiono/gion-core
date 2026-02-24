package repospec

import (
	"testing"
)

func TestValidateRepoKey(t *testing.T) {
	cases := []struct {
		name    string
		repoKey string
		wantErr bool
	}{
		{
			name:    "valid simple",
			repoKey: "github.com/owner/repo",
			wantErr: false,
		},
		{
			name:    "valid with subgroups",
			repoKey: "gitlab.com/group/subgroup/repo",
			wantErr: false,
		},
		{
			name:    "valid with subgroups deep",
			repoKey: "gitlab.com/org/sub/team/repo",
			wantErr: false,
		},
		{
			name:    "valid with port",
			repoKey: "gitlab.com:2222/org/repo",
			wantErr: false,
		},
		{
			name:    "valid with .git suffix",
			repoKey: "github.com/owner/repo.git",
			wantErr: false,
		},
		{
			name:    "invalid too few parts",
			repoKey: "github.com/owner",
			wantErr: true,
		},
		{
			name:    "invalid empty part",
			repoKey: "github.com//repo",
			wantErr: true,
		},
		{
			name:    "invalid whitespace",
			repoKey: "github.com/ owner/repo",
			wantErr: true,
		},
		{
			name:    "invalid .git in middle",
			repoKey: "github.com/owner/.git/repo",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateRepoKey(tc.repoKey)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateRepoKey() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid ssh",
			url:     "git@github.com:owner/repo.git",
			wantErr: false,
		},
		{
			name:    "valid https",
			url:     "https://github.com/owner/repo.git",
			wantErr: false,
		},
		{
			name:    "valid gitlab nested",
			url:     "git@gitlab.com:team/sub/repo.git",
			wantErr: false,
		},
		{
			name:    "valid file",
			url:     "file:///tmp/repos/example.com/owner/repo.git",
			wantErr: false,
		},
		{
			name:    "invalid empty",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid shorthand",
			url:     "owner/repo",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.url)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidateSpec(t *testing.T) {
	cases := []struct {
		name    string
		spec    Spec
		wantErr bool
	}{
		{
			name: "valid spec",
			spec: Spec{
				EndPoint: EndPoint{
					Host: "github.com",
				},
				Registry: Registry{
					Provider: ProviderGitHub,
					Group:    "owner",
				},
				Repository: Repository{
					Repo: "repo",
				},
				RepoKey: "github.com/owner/repo",
				IsSSH:   true,
			},
			wantErr: false,
		},
		{
			name: "valid with subgroups",
			spec: Spec{
				EndPoint: EndPoint{
					Host: "gitlab.com",
				},
				Registry: Registry{
					Provider:  ProviderGitLab,
					Group:     "team",
					SubGroups: []string{"sub"},
				},
				Repository: Repository{
					Repo: "repo",
				},
				RepoKey: "gitlab.com/team/sub-repo",
				IsSSH:   true,
			},
			wantErr: false,
		},
		{
			name: "valid with port",
			spec: Spec{
				EndPoint: EndPoint{
					Host: "gitlab.com",
					Port: "2222",
				},
				Registry: Registry{
					Provider: ProviderGitLab,
					Group:    "owner",
				},
				Repository: Repository{
					Repo: "repo",
				},
				RepoKey: "gitlab.com:2222/owner/repo",
				IsSSH:   true,
			},
			wantErr: false,
		},
		{
			name: "valid with basePath",
			spec: Spec{
				EndPoint: EndPoint{
					Host:     "git.example.com",
					BasePath: "/git",
				},
				Registry: Registry{
					Provider: ProviderCustom,
					Group:    "owner",
				},
				Repository: Repository{
					Repo: "repo",
				},
				RepoKey: "git.example.com/owner/repo",
				IsSSH:   false,
			},
			wantErr: false,
		},
		{
			name: "invalid empty host",
			spec: Spec{
				EndPoint: EndPoint{
					Host: "",
				},
				Registry: Registry{
					Provider: ProviderGitHub,
					Group:    "owner",
				},
				Repository: Repository{
					Repo: "repo",
				},
				RepoKey: "/owner/repo",
				IsSSH:   true,
			},
			wantErr: true,
		},
		{
			name: "invalid empty group",
			spec: Spec{
				EndPoint: EndPoint{
					Host: "github.com",
				},
				Registry: Registry{
					Provider: ProviderGitHub,
					Group:    "",
				},
				Repository: Repository{
					Repo: "repo",
				},
				RepoKey: "github.com//repo",
				IsSSH:   true,
			},
			wantErr: true,
		},
		{
			name: "invalid empty repo",
			spec: Spec{
				EndPoint: EndPoint{
					Host: "github.com",
				},
				Registry: Registry{
					Provider: ProviderGitHub,
					Group:    "owner",
				},
				Repository: Repository{
					Repo: "",
				},
				RepoKey: "github.com/owner/",
				IsSSH:   true,
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSpec(tc.spec)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateSpec() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestNormalizeAndValidate(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid ssh",
			input:   "git@github.com:owner/repo.git",
			wantErr: false,
		},
		{
			name:    "valid gitlab nested",
			input:   "git@gitlab.com:team/sub/repo.git",
			wantErr: false,
		},
		{
			name:    "invalid url",
			input:   "owner/repo",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NormalizeAndValidate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("NormalizeAndValidate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
