package app

import "fmt"

const (
	defaultString = "undefined"
)

var (
	version            = defaultString
	commit             = defaultString
	branch             = defaultString
	buildUnixTimestamp = defaultString
)

func getBuildInfo() string {
	return fmt.Sprintf(
		"Version: %s, Commit: %s, Branch: %s, Build time: %s",
		version, commit, branch, buildUnixTimestamp,
	)
}
