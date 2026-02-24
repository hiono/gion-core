package repospec

import (
	"fmt"
	"strings"
)

type Config struct {
	Version    int                  `yaml:"version"`
	Presets    map[string]Preset    `yaml:"presets"`
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

type Preset struct {
	Repos []PresetRepoEntry `yaml:"repos"`
}

type PresetRepoEntry struct {
	Repo     string `yaml:"repo"`
	Provider string `yaml:"provider"`
	BasePath string `yaml:"base_path"`
}

type Workspace struct {
	Description string               `yaml:"description"`
	Mode        string               `yaml:"mode"`
	SourceURL   string               `yaml:"source_url"`
	PresetName  string               `yaml:"preset_name"`
	Repos       []WorkspaceRepoEntry `yaml:"repos"`
}

type WorkspaceRepoEntry struct {
	Alias   string `yaml:"alias"`
	RepoKey string `yaml:"repo_key"`
	Branch  string `yaml:"branch"`
	BaseRef string `yaml:"base_ref"`
}

var validModes = map[string]bool{
	"preset": true,
	"repo":   true,
	"review": true,
	"issue":  true,
	"resume": true,
	"add":    true,
}

func ValidateConfig(cfg Config) []ValidationError {
	var errs []ValidationError

	// Version validation
	if cfg.Version < 1 {
		errs = append(errs, ValidationError{Field: "version", Message: "version must be >= 1"})
	}

	// Presets validation
	if cfg.Presets != nil {
		for name, preset := range cfg.Presets {
			if err := validatePreset(name, preset); err != nil {
				errs = append(errs, err...)
			}
		}
	}

	// Workspaces validation
	if cfg.Workspaces != nil {
		for name, ws := range cfg.Workspaces {
			if err := validateWorkspace(name, ws); err != nil {
				errs = append(errs, err...)
			}
		}
	}

	return errs
}

func validatePreset(name string, preset Preset) []ValidationError {
	var errs []ValidationError

	if len(preset.Repos) == 0 {
		errs = append(errs, ValidationError{
			Field:   fmt.Sprintf("presets.%s.repos", name),
			Message: "repos cannot be empty",
		})
	}

	for i, repo := range preset.Repos {
		if repo.Repo == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("presets.%s.repos[%d].repo", name, i),
				Message: "repo is required",
			})
		}

		// Validate provider if specified
		if repo.Provider != "" {
			if !isValidProvider(repo.Provider) {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("presets.%s.repos[%d].provider", name, i),
					Message: fmt.Sprintf("invalid provider: %s (must be github, gitlab, bitbucket, or custom)", repo.Provider),
				})
			}
		}

		// Validate base_path format
		if repo.BasePath != "" && !strings.HasPrefix(repo.BasePath, "/") {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("presets.%s.repos[%d].base_path", name, i),
				Message: "base_path must start with /",
			})
		}
	}

	return errs
}

func validateWorkspace(name string, ws Workspace) []ValidationError {
	var errs []ValidationError

	// Mode validation
	if ws.Mode == "" {
		errs = append(errs, ValidationError{
			Field:   fmt.Sprintf("workspaces.%s.mode", name),
			Message: "mode is required",
		})
	} else if !validModes[ws.Mode] {
		errs = append(errs, ValidationError{
			Field:   fmt.Sprintf("workspaces.%s.mode", name),
			Message: fmt.Sprintf("invalid mode: %s (must be preset, repo, review, issue, resume, or add)", ws.Mode),
		})
	}

	// mode-specific validation
	switch ws.Mode {
	case "preset":
		if ws.PresetName == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("workspaces.%s.preset_name", name),
				Message: "preset_name is required for mode=preset",
			})
		}
	case "repo", "issue", "review", "resume", "add":
		if len(ws.Repos) == 0 {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("workspaces.%s.repos", name),
				Message: "repos cannot be empty for this mode",
			})
		}
	}

	// Repos validation
	for i, repo := range ws.Repos {
		if repo.Alias == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("workspaces.%s.repos[%d].alias", name, i),
				Message: "alias is required",
			})
		}

		if repo.RepoKey == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("workspaces.%s.repos[%d].repo_key", name, i),
				Message: "repo_key is required",
			})
		} else {
			// Validate repo_key format (host/group/repo)
			if err := validateRepoKeyFormat(repo.RepoKey); err != nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("workspaces.%s.repos[%d].repo_key", name, i),
					Message: err.Error(),
				})
			}
		}

		if repo.Branch == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("workspaces.%s.repos[%d].branch", name, i),
				Message: "branch is required",
			})
		}

		// Validate base_ref format
		if repo.BaseRef != "" && !strings.HasPrefix(repo.BaseRef, "origin/") {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("workspaces.%s.repos[%d].base_ref", name, i),
				Message: "base_ref must start with origin/",
			})
		}
	}

	return errs
}

func validateRepoKeyFormat(repoKey string) error {
	trimmed := strings.TrimSuffix(strings.TrimSpace(repoKey), ".git")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 3 {
		return fmt.Errorf("repo_key must be host/group/repo[.git] or host/group/subgroup/repo[.git]")
	}
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return fmt.Errorf("repo_key must be host/group/repo[.git] or host/group/subgroup/repo[.git]")
		}
	}
	// Validate no .git suffix in middle parts
	for i, part := range parts {
		if strings.HasSuffix(part, ".git") && i < len(parts)-1 {
			return fmt.Errorf(".git suffix only allowed at end of repo_key")
		}
	}
	return nil
}

func isValidProvider(p string) bool {
	switch ProviderType(p) {
	case ProviderGitHub, ProviderGitLab, ProviderBitbucket, ProviderCustom:
		return true
	}
	return false
}
