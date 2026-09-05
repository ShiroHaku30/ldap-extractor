package version

import "fmt"

var (
	Version   = "dev"
	GitTag    = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func String() string {
	return fmt.Sprintf(
		"Version: %s\nTag: %s\nCommit: %s\nBuilt: %s",
		Version,
		GitTag,
		GitCommit,
		BuildTime,
	)
}
