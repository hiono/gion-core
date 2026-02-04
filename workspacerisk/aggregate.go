package workspacerisk

type RepoState string

type WorkspaceRisk string

const (
	RepoStateUnknown  RepoState = "unknown"
	RepoStateDirty    RepoState = "dirty"
	RepoStateDiverged RepoState = "diverged"
	RepoStateUnpushed RepoState = "unpushed"
	RepoStateClean    RepoState = "clean"
)

const (
	WorkspaceRiskUnknown  WorkspaceRisk = "unknown"
	WorkspaceRiskDirty    WorkspaceRisk = "dirty"
	WorkspaceRiskDiverged WorkspaceRisk = "diverged"
	WorkspaceRiskUnpushed WorkspaceRisk = "unpushed"
	WorkspaceRiskClean    WorkspaceRisk = "clean"
)

func Aggregate(repos []RepoState) WorkspaceRisk {
	hasDirty := false
	hasUnknown := false
	hasDiverged := false
	hasUnpushed := false
	for _, repo := range repos {
		switch repo {
		case RepoStateUnknown:
			hasUnknown = true
		case RepoStateDirty:
			hasDirty = true
		case RepoStateDiverged:
			hasDiverged = true
		case RepoStateUnpushed:
			hasUnpushed = true
		}
	}
	switch {
	case hasUnknown:
		return WorkspaceRiskUnknown
	case hasDirty:
		return WorkspaceRiskDirty
	case hasDiverged:
		return WorkspaceRiskDiverged
	case hasUnpushed:
		return WorkspaceRiskUnpushed
	default:
		return WorkspaceRiskClean
	}
}
