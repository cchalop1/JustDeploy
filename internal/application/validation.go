package application

import (
	"regexp"
	"strings"
)

var repoPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*/[A-Za-z0-9_.-]+$`)
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,99}$`)
var labelPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)

func ValidRepo(repo string) bool {
	return repoPattern.MatchString(repo) && !strings.Contains(repo, "..")
}
func ValidID(id string) bool { return idPattern.MatchString(id) }
func ValidDomain(domain string) bool {
	if len(domain) > 253 || strings.Contains(domain, "`") {
		return false
	}
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if !labelPattern.MatchString(label) {
			return false
		}
	}
	return true
}
