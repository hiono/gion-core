package gitparse

import "strings"

func ParseHeadRefs(showRefOutput string) []string {
	refs := make([]string, 0)
	lines := strings.Split(strings.TrimSpace(showRefOutput), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		ref := parts[1]
		if !strings.HasPrefix(ref, "refs/heads/") {
			continue
		}
		refs = append(refs, ref)
	}
	return refs
}
