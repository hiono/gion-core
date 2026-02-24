package repospec

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

var (
	cueCtx *cue.Context
)

func init() {
	cueCtx = cuecontext.New()
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func ValidateSpec(spec Spec) error {
	// Basic validation first
	if spec.EndPoint.Host == "" {
		return &ValidationError{Field: "endPoint.host", Message: "host is required"}
	}
	if spec.Registry.Group == "" {
		return &ValidationError{Field: "registry.group", Message: "group is required"}
	}
	if spec.Repository.Repo == "" {
		return &ValidationError{Field: "repository.repo", Message: "repo is required"}
	}

	// Additional validation: no .git in middle
	for i, part := range spec.Registry.SubGroups {
		if strings.HasSuffix(part, ".git") && i < len(spec.Registry.SubGroups)-1 {
			return &ValidationError{Field: "registry.subGroups", Message: ".git suffix not allowed in subgroups"}
		}
	}

	// Validate port format if present (1-65535)
	if spec.EndPoint.Port != "" {
		port := spec.EndPoint.Port
		// Port should be numeric
		for _, c := range port {
			if c < '0' || c > '9' {
				return &ValidationError{Field: "endPoint.port", Message: "port must be numeric"}
			}
		}
		// Validate port range
		if port == "0" || len(port) > 5 || (len(port) == 5 && port > "65535") {
			return &ValidationError{Field: "endPoint.port", Message: "port must be between 1 and 65535"}
		}
	}

	// Validate basePath format if present
	if spec.EndPoint.BasePath != "" && !strings.HasPrefix(spec.EndPoint.BasePath, "/") {
		return &ValidationError{Field: "endPoint.basePath", Message: "basePath must start with /"}
	}

	// Provider-specific constraints
	if err := validateProviderConstraints(spec); err != nil {
		return err
	}

	return nil
}

func validateProviderConstraints(spec Spec) error {
	switch spec.Registry.Provider {
	case ProviderGitHub, ProviderBitbucket:
		// GitHub and Bitbucket do not support subGroups (only owner/repo)
		if len(spec.Registry.SubGroups) > 0 {
			return &ValidationError{
				Field:   "registry.subGroups",
				Message: fmt.Sprintf("subGroups not supported for %s (use owner/repo format)", spec.Registry.Provider),
			}
		}
	case ProviderGitLab:
		// GitLab supports nested groups (max 20 levels total = 19 subgroups + 1 repo)
		if len(spec.Registry.SubGroups) > 19 {
			return &ValidationError{
				Field:   "registry.subGroups",
				Message: "subGroups exceeds maximum depth of 19 for GitLab",
			}
		}
	}
	return nil
}

func ValidateRepoKey(repoKey string) error {
	if strings.ContainsAny(repoKey, " \t\r\n") {
		return &ValidationError{Field: "repoKey", Message: "must not contain whitespace"}
	}

	trimmed := strings.TrimSuffix(strings.TrimSpace(repoKey), ".git")
	parts := strings.Split(trimmed, "/")

	if len(parts) < 3 {
		return &ValidationError{
			Field:   "repoKey",
			Message: "must be host/group/repo or host/group/subgroup/repo",
		}
	}

	for i, part := range parts {
		if strings.TrimSpace(part) == "" {
			return &ValidationError{
				Field:   "repoKey",
				Message: "must be host/group/repo or host/group/subgroup/repo",
			}
		}
		// Validate no .git suffix in middle parts
		if strings.HasSuffix(part, ".git") && i < len(parts)-1 {
			return &ValidationError{
				Field:   "repoKey",
				Message: ".git suffix only allowed at end",
			}
		}
	}

	return nil
}

func ValidateURL(url string) error {
	// Basic URL validation
	trimmed := strings.TrimSpace(url)

	if trimmed == "" {
		return &ValidationError{Field: "url", Message: "cannot be empty"}
	}

	// Check for valid URL patterns
	isSSH := strings.HasPrefix(trimmed, "git@") || strings.HasPrefix(trimmed, "ssh://")
	isHTTPS := strings.HasPrefix(trimmed, "https://")
	isHTTP := strings.HasPrefix(trimmed, "http://")
	isFile := strings.HasPrefix(trimmed, "file://")

	if !isSSH && !isHTTPS && !isHTTP && !isFile {
		return &ValidationError{
			Field:   "url",
			Message: "must be SSH (git@host:path), HTTPS, HTTP, or file:// URL",
		}
	}

	// Parse and validate the URL using Normalize
	_, err := Normalize(trimmed)
	if err != nil {
		return &ValidationError{Field: "url", Message: err.Error()}
	}

	return nil
}

func parseCUEError(err error) error {
	errStr := err.Error()

	// Parse CUE error messages and make them user-friendly
	if strings.Contains(errStr, "invalid value") {
		// Extract field name from error
		return &ValidationError{Message: errStr}
	}

	return &ValidationError{Message: errStr}
}

func NormalizeAndValidate(input string) (Spec, error) {
	spec, err := Normalize(input)
	if err != nil {
		return Spec{}, err
	}

	if err := ValidateSpec(spec); err != nil {
		return Spec{}, err
	}

	return spec, nil
}
