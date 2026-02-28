package repospec

import (
	"testing"
)

func TestDetectProvider(t *testing.T) {
	tests := []struct {
		name string
		host string
		want ProviderType
	}{
		// GitHub
		{"github.com", "github.com", ProviderGitHub},
		{"github.com with www", "www.github.com", ProviderGitHub},
		// Enterprise URLs cannot be detected from string alone - requires API access
		{"github enterprise", "github.mycompany.com", ProviderCustom},

		// GitLab
		{"gitlab.com", "gitlab.com", ProviderGitLab},
		{"gitlab with www", "www.gitlab.com", ProviderGitLab},
		// Enterprise URLs cannot be detected from string alone
		{"self-hosted gitlab", "gitlab.mycompany.com", ProviderCustom},

		// Bitbucket
		{"bitbucket.org", "bitbucket.org", ProviderBitbucket},
		{"bitbucket with www", "www.bitbucket.org", ProviderBitbucket},
		// Enterprise URLs cannot be detected from string alone
		{"self-hosted bitbucket", "bitbucket.mycompany.com", ProviderCustom},

		// Custom
		{"custom domain", "git.example.com", ProviderCustom},
		{"custom with subdomain", "git.company.com", ProviderCustom},
		{"unknown", "example.com", ProviderCustom},

		// Case insensitive
		{"uppercase github", "GITHUB.COM", ProviderGitHub},
		{"mixed case", "GitHub.Com", ProviderGitHub},
		{"uppercase gitlab", "GITLAB.COM", ProviderGitLab},

		// Port stripping
		{"github with port", "github.com:8080", ProviderGitHub},
		{"gitlab with port", "gitlab.com:443", ProviderGitLab},
		{"custom with port", "git.example.com:2222", ProviderCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectProvider(tt.host); got != tt.want {
				t.Errorf("DetectProvider(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestDetectProviderEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		host string
		want ProviderType
	}{
		// Invalid URLs - GitHub doesn't accept path/query/fragment in repo URLs
		{"with path", "github.com/path/to/repo", ProviderCustom},
		{"with query", "github.com?ref=main", ProviderCustom},
		{"with fragment", "github.com#readme", ProviderCustom},
		{"IP address", "192.168.1.1", ProviderCustom},
		{"double port", "github.com:8080:9090", ProviderCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectProvider(tt.host); got != tt.want {
				t.Errorf("DetectProvider(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestProviderString(t *testing.T) {
	tests := []struct {
		pt   ProviderType
		want string
	}{
		{ProviderGitHub, "github"},
		{ProviderGitLab, "gitlab"},
		{ProviderBitbucket, "bitbucket"},
		{ProviderCustom, "custom"},
		{ProviderType("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := string(tt.pt)
			if got != tt.want {
				t.Errorf("ProviderType(%v) = %q, want %q", tt.pt, got, tt.want)
			}
		})
	}
}
