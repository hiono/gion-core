package repospec

import (
	"fmt"
	"net/url"
	"strings"
)

type Spec struct {
	Host    string
	Owner   string
	Repo    string
	RepoKey string
	IsSSH   bool
}

func Normalize(input string) (Spec, error) {
	return NormalizeWithBasePath(input, "")
}

func NormalizeWithBasePath(input string, basePath string) (Spec, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return Spec{}, fmt.Errorf("repo spec is empty")
	}

	var host string
	var path string
	var isSSH bool

	switch {
	case strings.HasPrefix(trimmed, "git@"):
		isSSH = true
		at := strings.Index(trimmed, "@")
		if at < 0 {
			return Spec{}, fmt.Errorf("invalid ssh repo spec: %q", input)
		}
		rest := trimmed[at+1:]
		colons := findAllColons(rest)
		if len(colons) == 0 {
			return Spec{}, fmt.Errorf("invalid ssh repo spec: %q", input)
		}
		if len(colons) >= 2 && isPort(rest[colons[0]+1:colons[1]]) {
			host = rest[:colons[0]]
			path = rest[colons[1]+1:]
		} else {
			host = rest[:colons[0]]
			path = rest[colons[0]+1:]
		}
	case strings.HasPrefix(trimmed, "https://"):
		u, err := url.Parse(trimmed)
		if err != nil {
			return Spec{}, fmt.Errorf("invalid https repo spec: %q", input)
		}
		host = u.Hostname()
		path = strings.TrimPrefix(u.Path, "/")
	case strings.HasPrefix(trimmed, "file://"):
		u, err := url.Parse(trimmed)
		if err != nil {
			return Spec{}, fmt.Errorf("invalid file repo spec: %q", input)
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) < 3 {
			return Spec{}, fmt.Errorf("file repo spec must end with <host>/<owner>/<repo>: %q", input)
		}
		host = parts[len(parts)-3]
		owner := parts[len(parts)-2]
		repo := parts[len(parts)-1]
		path = fmt.Sprintf("%s/%s", owner, repo)
	default:
		return Spec{}, fmt.Errorf("repo spec must be ssh, https, or file: %q", input)
	}

	if basePath != "" && !isSSH {
		path = removeBasePath(path, basePath)
	}

	owner, repo, err := splitOwnerRepo(path)
	if err != nil {
		return Spec{}, err
	}
	if host == "" {
		return Spec{}, fmt.Errorf("host is required in repo spec: %q", input)
	}

	spec := Spec{
		Host:    host,
		Owner:   owner,
		Repo:    repo,
		RepoKey: fmt.Sprintf("%s/%s/%s", host, owner, repo),
		IsSSH:   isSSH,
	}
	return spec, nil
}

func removeBasePath(path string, basePath string) string {
	basePath = strings.Trim(basePath, "/")
	if basePath == "" {
		return path
	}
	path = strings.TrimPrefix(path, basePath+"/")
	return path
}

func splitOwnerRepo(path string) (string, string, error) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "", "", fmt.Errorf("repo path is empty")
	}

	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("repo path must have at least owner/repo")
	}

	repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
	if repo == "" {
		return "", "", fmt.Errorf("repo name cannot be empty")
	}
	owner := strings.Join(parts[:len(parts)-1], "/")
	if owner == "" {
		return "", "", fmt.Errorf("owner/namespace cannot be empty")
	}

	return owner, repo, nil
}

func findAllColons(s string) []int {
	var indices []int
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			indices = append(indices, i)
		}
	}
	return indices
}

func isPort(s string) bool {
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
