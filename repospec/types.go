package repospec

type ProviderType string

const (
	ProviderGitHub    ProviderType = "github"
	ProviderGitLab    ProviderType = "gitlab"
	ProviderBitbucket ProviderType = "bitbucket"
	ProviderCustom    ProviderType = "custom"
)

type EndPoint struct {
	Host     string // github.com, gitlab.com, self-hosted.example.com
	Port     string // :22, :8080 (optional)
	BasePath string // nginx reverse proxy path (only for HTTPS)
}

type Registry struct {
	Provider  ProviderType // Auto-detected or manually specified
	Group     string       // Primary group/owner (e.g., org)
	SubGroups []string     // Subgroups under group (optional, e.g., [subgroup1, subgroup2])
}

type Repository struct {
	Repo string // Repository name (e.g., myapp)
}

func (r Repository) Last() string {
	return r.Repo
}

type Spec struct {
	EndPoint
	Registry
	Repository
	RepoKey string
	IsSSH   bool
}
