package storage

import (
	"fmt"
	"regexp"
)

var (
	validTagRegex = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`)
)

func validateTag(tag string) error {
	if !validTagRegex.MatchString(tag) {
		return fmt.Errorf("invalid tag")
	}
	return nil
}
