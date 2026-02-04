package workspacerisk

func AggregateForState(repos []RepoState) WorkspaceRisk {
	hasDirty := false
	hasUnknown := false
	hasDiverged := false
	hasUnpushed := false
	for _, repo := range repos {
		switch repo {
		case RepoStateDirty:
			hasDirty = true
		case RepoStateUnknown:
			hasUnknown = true
		case RepoStateDiverged:
			hasDiverged = true
		case RepoStateUnpushed:
			hasUnpushed = true
		}
	}
	switch {
	case hasDirty:
		return WorkspaceRiskDirty
	case hasUnknown:
		return WorkspaceRiskUnknown
	case hasDiverged:
		return WorkspaceRiskDiverged
	case hasUnpushed:
		return WorkspaceRiskUnpushed
	default:
		return WorkspaceRiskClean
	}
}
