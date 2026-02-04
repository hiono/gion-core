package gcplan

import "fmt"

func MergeCheckRefs(branch string, target string) (headRef string, targetRef string) {
	return fmt.Sprintf("refs/heads/%s", branch), fmt.Sprintf("refs/remotes/%s", target)
}
