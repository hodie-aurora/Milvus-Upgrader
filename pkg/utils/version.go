package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Version 表示版本号
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion 解析版本字符串（如 "2.1.3" 或 "v2.1.3"）
func ParseVersion(v string) (Version, error) {
	v = strings.TrimPrefix(v, "v") // 移除 "v" 前缀
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

// IsMinorUpgrade 判断是否为小版本升级
func IsMinorUpgrade(current, target Version) bool {
	return current.Major == target.Major
}
