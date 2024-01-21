package main

import (
	"log/slog"
	"runtime/debug"
)

// BuildID indicates the exact version of the application being run
var BuildID string

func init() {
	// check if the build ID was set during compilation
	if BuildID != "" {
		return
	}

	// read the build information from the running binary
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		slog.Error("failed to read build info from self")
		BuildID = "unknown"
		return
	}

	revision := "dev"
	modified := false

	// loop through the build settings to pull VCS info out
	for _, setting := range buildInfo.Settings {
		// check which revision this build was based on
		if setting.Key == "vcs.revision" {
			revision = setting.Value[:7]
		}

		// check if the project has been modified since the revision
		if setting.Key == "vcs.modified" {
			modified = setting.Value == "true"
		}
	}

	// set the build ID based on this information
	if modified {
		revision += "-dirty"
	}
	BuildID = revision
}
