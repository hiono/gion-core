package gitparse

import "strings"

func ParseWorktreeBranchNames(worktreeListOutput string) []string {
	branches := map[string]struct{}{}
	lines := strings.Split(strings.TrimSpace(worktreeListOutput), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "branch ") {
			continue
		}
		ref := strings.TrimSpace(strings.TrimPrefix(line, "branch "))
		if strings.HasPrefix(ref, "refs/heads/") {
			name := strings.TrimPrefix(ref, "refs/heads/")
			if name != "" {
				branches[name] = struct{}{}
			}
		}
	}
	result := make([]string, 0, len(branches))
	for name := range branches {
		result = append(result, name)
	}
	return result
}
