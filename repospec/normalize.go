package repospec

import (
	"fmt"
	"net/url"
	"strings"
)

func Normalize(input string) (Spec, error) {
	return NormalizeWithBasePath(input, "")
}

func NormalizeWithBasePath(input string, basePath string) (Spec, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return Spec{}, fmt.Errorf("repo spec is empty")
	}

	var host, port string
	var path string
	var isSSH bool

	switch {
	case strings.HasPrefix(trimmed, "ssh://"):
		isSSH = true
		rest := strings.TrimPrefix(trimmed, "ssh://")
		at := strings.Index(rest, "@")
		if at < 0 {
			return Spec{}, fmt.Errorf("invalid ssh repo spec: %q", input)
		}
		rest = rest[at+1:]
		slashIdx := strings.Index(rest, "/")
		if slashIdx < 0 {
			return Spec{}, fmt.Errorf("invalid ssh repo spec: %q", input)
		}
		hostPart := rest[:slashIdx]
		path = rest[slashIdx+1:]
		host, port = parseHostPort(hostPart)
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
			port = rest[colons[0]+1 : colons[1]]
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

	owner, pathPartsForSlug, _, err := splitOwnerRepo(path)
	if err != nil {
		return Spec{}, err
	}
	if host == "" {
		return Spec{}, fmt.Errorf("host is required in repo spec: %q", input)
	}

	repoKey := buildRepoKey(host, port, owner, pathPartsForSlug)
	provider := DetectProvider(host)

	// Split fullPath into subgroups + repo
	// pathPartsForSlug = [group, subgroup, ..., repo]
	// subgroups = pathPartsForSlug without the last element (repo)
	var subgroups []string
	var repoName string
	if len(pathPartsForSlug) > 1 {
		subgroups = pathPartsForSlug[:len(pathPartsForSlug)-1]
		repoName = pathPartsForSlug[len(pathPartsForSlug)-1]
	} else if len(pathPartsForSlug) == 1 {
		repoName = pathPartsForSlug[0]
	} else {
		repoName = ""
	}

	spec := Spec{
		EndPoint: EndPoint{
			Host:     host,
			Port:     port,
			BasePath: basePath,
		},
		Registry: Registry{
			Provider:  provider,
			Group:     owner,
			SubGroups: subgroups,
		},
		Repository: Repository{
			Repo: repoName,
		},
		RepoKey: repoKey,
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

func splitOwnerRepo(path string) (string, []string, []string, error) {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return "", nil, nil, fmt.Errorf("repo path is empty")
	}

	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return "", nil, nil, fmt.Errorf("repo path must be host/group/repo or host/group/subgroup/repo")
	}

	repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
	if repo == "" {
		return "", nil, nil, fmt.Errorf("repo name cannot be empty")
	}

	owner := parts[0]
	if owner == "" {
		return "", nil, nil, fmt.Errorf("owner/namespace cannot be empty")
	}

	// pathPartsForSlug: parts after owner (groups + repo), with .git stripped
	pathPartsForSlug := make([]string, len(parts)-1)
	copy(pathPartsForSlug, parts[1:])
	pathPartsForSlug[len(pathPartsForSlug)-1] = repo

	// fullPath: all parts including owner, with .git stripped
	fullPath := make([]string, len(parts))
	copy(fullPath, parts)
	fullPath[len(fullPath)-1] = repo

	return owner, pathPartsForSlug, fullPath, nil
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

func parseHostPort(hostPart string) (host, port string) {
	colonIdx := strings.LastIndex(hostPart, ":")
	if colonIdx > 0 {
		potentialPort := hostPart[colonIdx+1:]
		if isPort(potentialPort) {
			return hostPart[:colonIdx], potentialPort
		}
	}
	return hostPart, ""
}

func buildRepoKey(host, port string, owner string, pathParts []string) string {
	slug := Slugify(pathParts)
	if port != "" {
		return fmt.Sprintf("%s:%s/%s/%s", host, port, owner, slug)
	}
	return fmt.Sprintf("%s/%s/%s", host, owner, slug)
}
