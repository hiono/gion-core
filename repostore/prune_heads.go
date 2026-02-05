package repostore

import (
	"sort"
	"strings"
)

func SelectPrunableHeadRefs(headRefs []string, keepBranch string, worktreeBranches []string) []string {
	keepBranch = strings.TrimSpace(keepBranch)
	worktreeSet := map[string]struct{}{}
	for _, branch := range worktreeBranches {
		branch = strings.TrimSpace(branch)
		if branch == "" {
			continue
		}
		worktreeSet[branch] = struct{}{}
	}
	prunable := make([]string, 0)
	for _, ref := range headRefs {
		if !strings.HasPrefix(ref, "refs/heads/") {
			continue
		}
		name := strings.TrimPrefix(ref, "refs/heads/")
		if name == "" {
			continue
		}
		if keepBranch != "" && name == keepBranch {
			continue
		}
		if _, ok := worktreeSet[name]; ok {
			continue
		}
		prunable = append(prunable, ref)
	}
	sort.Strings(prunable)
	return prunable
}
