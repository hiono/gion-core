package repospec

import (
	"fmt"
	"strings"
)

func DisplaySpec(input string) string {
	spec, ok := normalizeForDisplay(input)
	if !ok {
		return strings.TrimSpace(input)
	}
	return fmt.Sprintf("git@%s:%s/%s.git", spec.Host, spec.Owner, spec.Repo)
}

func DisplayName(input string) string {
	spec, ok := normalizeForDisplay(input)
	if !ok || spec.Repo == "" {
		return strings.TrimSpace(input)
	}
	return spec.Repo
}

func SpecFromKey(repoKey string) string {
	return SpecFromKeyWithScheme(repoKey, true)
}

func SpecFromKeyWithScheme(repoKey string, isSSH bool) string {
	trimmed := strings.TrimSuffix(strings.TrimSpace(repoKey), ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 3 {
		return strings.TrimSpace(repoKey)
	}
	hostPart := parts[0]
	owner := strings.Join(parts[1:len(parts)-1], "/")
	repoName := parts[len(parts)-1]

	// Check if host contains port (host:port format)
	host, port := parseHostPort(hostPart)

	if isSSH {
		if port != "" {
			return fmt.Sprintf("ssh://git@%s:%s/%s/%s.git", host, port, owner, repoName)
		}
		return fmt.Sprintf("git@%s:%s/%s.git", host, owner, repoName)
	}
	if port != "" {
		return fmt.Sprintf("https://%s:%s/%s/%s.git", host, port, owner, repoName)
	}
	return fmt.Sprintf("https://%s/%s/%s.git", host, owner, repoName)
}

func parseHostPort(hostPart string) (host, port string) {
	colonIdx := strings.LastIndex(hostPart, ":")
	if colonIdx > 0 {
		potentialPort := hostPart[colonIdx+1:]
		if isAllDigits(potentialPort) {
			return hostPart[:colonIdx], potentialPort
		}
	}
	return hostPart, ""
}

func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func normalizeForDisplay(input string) (Spec, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return Spec{}, false
	}
	spec, err := Normalize(trimmed)
	if err != nil {
		return Spec{}, false
	}
	return spec, true
}
