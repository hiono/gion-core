package repospec

import "strings"

func Slugify(parts []string) string {
	var result []string
	for _, p := range parts {
		result = append(result, strings.ReplaceAll(p, "-", "--"))
	}
	return strings.Join(result, "-")
}

func Unslugify(slug string) []string {
	parts := strings.Split(slug, "-")
	var result []string
	for _, p := range parts {
		result = append(result, strings.ReplaceAll(p, "--", "-"))
	}
	return result
}
