package fibarohc2

import (
	"fmt"
	"strings"
)

const (
	// Hc2DefaultConfigFile is the default name of the configuration file
	Hc2DefaultConfigFile string = ".hc2-tools/config.json"

	// RepoName is the name of this repository
	RepoName string = "github.com/theovassiliou/hc2-tools"

	// Version contains the actuall version number. Might be replaces using the LDFLAGS.
	Version = "1.1.0-src"
)

// FormatFullVersion formats for a cmdName the version number based on version, branch and commit
func FormatFullVersion(cmdName, version, branch, commit string) string {
	var parts = []string{cmdName}

	if version != "" {
		parts = append(parts, version)
	} else {
		parts = append(parts, "unknown")
	}

	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}
