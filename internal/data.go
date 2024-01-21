package internal

import (
	_ "embed"
	"encoding/json"
	"log/slog"
)

// rawPackages is the raw JSON containing all known packages
//
//go:embed resources/packages.json
var rawPackages []byte

// packages is a map of all known packages
var packages map[string]string

func init() {
	// unmarshal the embedded JSON data into a usable map
	if err := json.Unmarshal(rawPackages, &packages); err != nil {
		slog.Error("failed to decode known packages", "error", err)
		panic(err)
	}
}
