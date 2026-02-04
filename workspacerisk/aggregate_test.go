package workspacerisk

import "testing"

func TestAggregate(t *testing.T) {
	tests := []struct {
		name  string
		repos []RepoState
		want  WorkspaceRisk
	}{
		{name: "clean", repos: []RepoState{RepoStateClean}, want: WorkspaceRiskClean},
		{name: "unpushed", repos: []RepoState{RepoStateUnpushed}, want: WorkspaceRiskUnpushed},
		{name: "diverged over unpushed", repos: []RepoState{RepoStateUnpushed, RepoStateDiverged}, want: WorkspaceRiskDiverged},
		{name: "dirty over diverged", repos: []RepoState{RepoStateDiverged, RepoStateDirty}, want: WorkspaceRiskDirty},
		{name: "unknown wins", repos: []RepoState{RepoStateDirty, RepoStateUnknown}, want: WorkspaceRiskUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Aggregate(tt.repos); got != tt.want {
				t.Fatalf("Aggregate() = %q, want %q", got, tt.want)
			}
		})
	}
}
