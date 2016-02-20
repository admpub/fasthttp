package fasthttp

import (
	"runtime"
)

func init() {
	if runtime.GOOS == "windows" {
		rootFS.PathRewrite = NewPathPrefixStripper(1)
	}
}
