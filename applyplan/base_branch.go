package applyplan

import "strings"

func UpdateBaseBranchCandidate(candidate string, mixed bool, baseBranch string) (string, bool) {
	if mixed {
		return candidate, mixed
	}
	baseBranch = strings.TrimSpace(baseBranch)
	if baseBranch == "" {
		return candidate, mixed
	}
	if candidate == "" {
		return baseBranch, mixed
	}
	if candidate != baseBranch {
		return candidate, true
	}
	return candidate, mixed
}
