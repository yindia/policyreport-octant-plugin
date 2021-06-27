package settings

import (
	"fmt"
)

const (
	name        = "policy-report"
	description = "Kubernetes-native policies"
	rootNavIcon = "boat"
)

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

func GetName() string {
	return name
}

func GetDescription(version VersionInfo) string {
	return fmt.Sprintf("%s (%s, %s, %s)", description, version.Version, version.Commit, version.Date)
}
