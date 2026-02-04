package gitparse

import (
	"fmt"
	"strconv"
	"strings"
)

type StatusPorcelainV2 struct {
	Branch         string
	Upstream       string
	Head           string
	Detached       bool
	HeadMissing    bool
	Dirty          bool
	UntrackedCount int
	StagedCount    int
	UnstagedCount  int
	UnmergedCount  int
	AheadCount     int
	BehindCount    int
}

func ParseStatusPorcelainV2(output, fallbackBranch string) StatusPorcelainV2 {
	status := StatusPorcelainV2{Branch: fallbackBranch}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			switch fields[1] {
			case "branch.oid":
				if fields[2] == "(initial)" {
					status.HeadMissing = true
				} else {
					status.Head = shortSHA(fields[2])
				}
			case "branch.head":
				switch fields[2] {
				case "(detached)":
					status.Detached = true
				case "(unknown)":
					status.HeadMissing = true
				default:
					status.Branch = fields[2]
				}
			case "branch.upstream":
				if fields[2] != "(unknown)" {
					status.Upstream = fields[2]
				}
			case "branch.ab":
				for _, field := range fields[2:] {
					if strings.HasPrefix(field, "+") {
						status.AheadCount = parseCount(field[1:])
					}
					if strings.HasPrefix(field, "-") {
						status.BehindCount = parseCount(field[1:])
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "? ") {
			status.UntrackedCount++
			status.Dirty = true
			continue
		}

		if strings.HasPrefix(line, "u ") {
			status.UnmergedCount++
			status.Dirty = true
			continue
		}
		if strings.HasPrefix(line, "1 ") || strings.HasPrefix(line, "2 ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				xy := fields[1]
				if len(xy) >= 2 {
					if xy[0] != '.' {
						status.StagedCount++
					}
					if xy[1] != '.' {
						status.UnstagedCount++
					}
					if xy[0] != '.' || xy[1] != '.' {
						status.Dirty = true
					}
				}
			}
			continue
		}
		status.Dirty = true
	}

	return status
}

func ParseChangedFilesPorcelainV2(output string) []string {
	var files []string
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			continue
		}
		switch {
		case strings.HasPrefix(line, "? "):
			path := strings.TrimSpace(strings.TrimPrefix(line, "? "))
			if path != "" {
				files = append(files, fmt.Sprintf("?? %s", path))
			}
		case strings.HasPrefix(line, "u "):
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			xy := shortXY(fields[1])
			path := fields[len(fields)-1]
			if strings.TrimSpace(path) != "" {
				files = append(files, fmt.Sprintf("%s %s", xy, path))
			}
		case strings.HasPrefix(line, "1 "):
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			xy := shortXY(fields[1])
			path := fields[len(fields)-1]
			if strings.TrimSpace(path) != "" {
				files = append(files, fmt.Sprintf("%s %s", xy, path))
			}
		case strings.HasPrefix(line, "2 "):
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			xy := shortXY(fields[1])
			oldPath := fields[len(fields)-2]
			newPath := fields[len(fields)-1]
			if strings.TrimSpace(oldPath) == "" || strings.TrimSpace(newPath) == "" {
				continue
			}
			files = append(files, fmt.Sprintf("%s %s -> %s", xy, oldPath, newPath))
		}
	}
	return files
}

func shortXY(xy string) string {
	value := strings.TrimSpace(xy)
	if value == "" {
		return "??"
	}
	value = strings.ReplaceAll(value, ".", " ")
	return value
}

func shortSHA(oid string) string {
	if len(oid) <= 7 {
		return oid
	}
	return oid[:7]
}

func parseCount(value string) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	if n < 0 {
		return 0
	}
	return n
}
