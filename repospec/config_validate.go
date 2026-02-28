package repospec

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

func ValidateConfig(cfg Config) []ValidationError {
	return nil
}
