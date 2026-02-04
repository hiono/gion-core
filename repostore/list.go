package repostore

import (
	"os"
	"path/filepath"
)

type Entry struct {
	RepoKey   string
	StorePath string
}

func List(reposRoot string) ([]Entry, []error, error) {
	info, err := os.Stat(reposRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, nil
	}

	var entries []Entry
	var warnings []error

	err = filepath.WalkDir(reposRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			warnings = append(warnings, walkErr)
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if path == reposRoot {
			return nil
		}
		if filepath.Ext(path) != ".git" {
			return nil
		}

		rel, err := filepath.Rel(reposRoot, path)
		if err != nil {
			warnings = append(warnings, err)
			return nil
		}
		repoKey := filepath.ToSlash(rel)
		entries = append(entries, Entry{RepoKey: repoKey, StorePath: path})
		return filepath.SkipDir
	})
	if err != nil {
		return nil, warnings, err
	}

	return entries, warnings, nil
}
