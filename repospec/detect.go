package repospec

import "strings"

func DetectProvider(host string) ProviderType {
	// Strip port if present
	colonIdx := strings.LastIndex(host, ":")
	if colonIdx > 0 {
		potentialPort := host[colonIdx+1:]
		if isPort(potentialPort) {
			host = host[:colonIdx]
		}
	}

	lowerHost := strings.ToLower(host)

	switch {
	case strings.HasSuffix(lowerHost, "github.com"):
		return ProviderGitHub
	case strings.HasSuffix(lowerHost, "gitlab.com"):
		return ProviderGitLab
	case strings.HasSuffix(lowerHost, "bitbucket.org"):
		return ProviderBitbucket
	default:
		return ProviderCustom
	}
}
