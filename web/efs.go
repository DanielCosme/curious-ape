package web

import "embed"

//go:embed "static" "dist"
var Files embed.FS
