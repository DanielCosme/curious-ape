package views

import (
	"embed"
	"github.com/benbjohnson/hashfs"
)

//go:embed "static/*"
var staticFS embed.FS

var StaticFS = hashfs.NewFS(staticFS)

func staticPath(path string) string {
	return "/" + StaticFS.HashName("static/"+path)
}
