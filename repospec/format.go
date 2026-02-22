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
	host := parts[0]
	owner := strings.Join(parts[1:len(parts)-1], "/")
	repoName := parts[len(parts)-1]
	if isSSH {
		return fmt.Sprintf("git@%s:%s/%s.git", host, owner, repoName)
	}
	return fmt.Sprintf("https://%s/%s/%s.git", host, owner, repoName)
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
