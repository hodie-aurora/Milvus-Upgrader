package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a version number
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string (e.g., "2.5.3" or "v2.5.3")
func ParseVersion(v string) (Version, error) {
	v = strings.TrimPrefix(v, "v") // Remove "v" prefix
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid version format: %s", v)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("invalid major version: %s", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}
	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// IsMinorUpgrade determines if it is a minor version upgrade
func IsMinorUpgrade(current, target Version) bool {
	return current.Major == target.Major
}
