package manifestlsplan

import (
	"sort"

	"github.com/tasuku43/gion-core/planner"
)

type DriftStatus string

const (
	DriftApplied DriftStatus = "applied"
	DriftMissing DriftStatus = "missing"
	DriftDrift   DriftStatus = "drift"
	DriftExtra   DriftStatus = "extra"
)

type DesiredWorkspace struct {
	ID          string
	Description string
}

type ManifestEntry struct {
	WorkspaceID  string
	Drift        DriftStatus
	Description  string
	HasWorkspace bool
}

type ExtraEntry struct {
	WorkspaceID string
	Drift       DriftStatus
}

type Counts struct {
	Applied int
	Drift   int
	Missing int
	Extra   int
}

type Result struct {
	ManifestEntries []ManifestEntry
	ExtraEntries    []ExtraEntry
	Counts          Counts
}

func Build(desired []DesiredWorkspace, changes []planner.WorkspaceChange, filesystemWorkspaceIDs []string) Result {
	statusByWorkspaceID := make(map[string]DriftStatus, len(desired))
	for _, change := range changes {
		switch change.Kind {
		case planner.WorkspaceAdd:
			statusByWorkspaceID[change.WorkspaceID] = DriftMissing
		case planner.WorkspaceUpdate:
			statusByWorkspaceID[change.WorkspaceID] = DriftDrift
		}
	}

	fsSet := map[string]struct{}{}
	for _, id := range filesystemWorkspaceIDs {
		if id == "" {
			continue
		}
		fsSet[id] = struct{}{}
	}

	sort.Slice(desired, func(i, j int) bool {
		return desired[i].ID < desired[j].ID
	})

	var counts Counts
	manifestEntries := make([]ManifestEntry, 0, len(desired))
	for _, ws := range desired {
		drift := statusByWorkspaceID[ws.ID]
		if drift == "" {
			drift = DriftApplied
		}
		_, hasWorkspace := fsSet[ws.ID]
		manifestEntries = append(manifestEntries, ManifestEntry{
			WorkspaceID:  ws.ID,
			Drift:        drift,
			Description:  ws.Description,
			HasWorkspace: hasWorkspace,
		})
		switch drift {
		case DriftApplied:
			counts.Applied++
		case DriftMissing:
			counts.Missing++
		case DriftDrift:
			counts.Drift++
		}
	}

	desiredSet := map[string]struct{}{}
	for _, ws := range desired {
		desiredSet[ws.ID] = struct{}{}
	}

	extraIDs := make([]string, 0)
	for id := range fsSet {
		if _, ok := desiredSet[id]; ok {
			continue
		}
		extraIDs = append(extraIDs, id)
	}
	sort.Strings(extraIDs)

	extras := make([]ExtraEntry, 0, len(extraIDs))
	for _, id := range extraIDs {
		extras = append(extras, ExtraEntry{WorkspaceID: id, Drift: DriftExtra})
		counts.Extra++
	}

	return Result{ManifestEntries: manifestEntries, ExtraEntries: extras, Counts: counts}
}
