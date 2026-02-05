package gitref

import "strings"

type RemoteBranch struct {
	Remote string
	Branch string
}

func ParseRemoteBranch(input string) (RemoteBranch, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return RemoteBranch{}, false
	}
	if strings.HasPrefix(trimmed, "refs/remotes/") {
		trimmed = strings.TrimPrefix(trimmed, "refs/remotes/")
	}
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return RemoteBranch{}, false
	}
	remote := strings.TrimSpace(parts[0])
	branch := strings.TrimSpace(parts[1])
	if remote == "" || branch == "" {
		return RemoteBranch{}, false
	}
	return RemoteBranch{Remote: remote, Branch: branch}, true
}
