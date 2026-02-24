package repostore

import (
	"path/filepath"

	"github.com/hiono/gion-core/repospec"
)

func StorePath(bareRoot string, spec repospec.Spec) string {
	// Level1: host/group merged (e.g., github.com/owner) - 1 level
	level1 := spec.EndPoint.Host + "/" + spec.Registry.Group

	// Handle empty repo case
	if spec.Repository.Repo == "" {
		return filepath.Join(bareRoot, level1) + ".git"
	}

	// Build full path for slug: subgroups + repo
	fullPath := append(spec.Registry.SubGroups, spec.Repository.Repo)
	// Level2: slugified path (e.g., group-subgroup-repo) - 1 level
	// Total: 2 levels from root
	level2 := repospec.Slugify(fullPath)
	return filepath.Join(bareRoot, level1, level2+".git")
}
