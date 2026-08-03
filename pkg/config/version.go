package config

import (
	"fmt"
	"runtime/debug"
)

const VERSION = "v2.2.1"

func Version() string {
	hash := "unknown"
	dirty := false
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				hash = s.Value[:7]
			case "vcs.modified":
				dirty = s.Value == "true"
			}
		}
	}
	if dirty {
		return fmt.Sprintf("%s-%s-dirty", VERSION, hash)
	}
	return fmt.Sprintf("%s-%s", VERSION, hash)
}
