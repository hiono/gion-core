package workspacerisk

import "testing"

func TestAggregateForState(t *testing.T) {
	tests := []struct {
		name  string
		repos []RepoState
		want  WorkspaceRisk
	}{
		{name: "clean", repos: []RepoState{RepoStateClean}, want: WorkspaceRiskClean},
		{name: "unpushed", repos: []RepoState{RepoStateUnpushed}, want: WorkspaceRiskUnpushed},
		{name: "diverged over unpushed", repos: []RepoState{RepoStateUnpushed, RepoStateDiverged}, want: WorkspaceRiskDiverged},
		{name: "dirty over diverged", repos: []RepoState{RepoStateDiverged, RepoStateDirty}, want: WorkspaceRiskDirty},
		{name: "dirty over unknown", repos: []RepoState{RepoStateUnknown, RepoStateDirty}, want: WorkspaceRiskDirty},
		{name: "unknown when no dirty", repos: []RepoState{RepoStateUnknown, RepoStateUnpushed}, want: WorkspaceRiskUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AggregateForState(tt.repos); got != tt.want {
				t.Fatalf("AggregateForState() = %q, want %q", got, tt.want)
			}
		})
	}
}
