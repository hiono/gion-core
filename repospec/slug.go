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
	if slug == "" {
		return []string{}
	}
	var result []string
	var current strings.Builder
	for i := 0; i < len(slug); i++ {
		if slug[i] == '-' {
			if i+1 < len(slug) && slug[i+1] == '-' {
				current.WriteByte('-')
				i++
			} else {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(slug[i])
		}
	}
	result = append(result, current.String())
	return result
}
