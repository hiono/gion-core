package repostore

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FetchGraceDuration(envValue string) time.Duration {
	val := strings.TrimSpace(envValue)
	if val == "" {
		return 30 * time.Second
	}
	secs, err := strconv.Atoi(val)
	if err != nil || secs < 0 {
		return 30 * time.Second
	}
	return time.Duration(secs) * time.Second
}

func RecentlyFetched(storePath string, grace time.Duration) bool {
	if grace <= 0 {
		return false
	}
	info, err := os.Stat(filepath.Join(storePath, "FETCH_HEAD"))
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) <= grace
}

func TouchFetchHead(storePath string) error {
	if strings.TrimSpace(storePath) == "" {
		return fmt.Errorf("store path is required")
	}
	path := filepath.Join(storePath, "FETCH_HEAD")
	now := time.Now()
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
		if err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return os.Chtimes(path, now, now)
}
