package main

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/target"
	"github.com/magefile/mage/mg"
)

// The `live:` tasks below are used together for development builds and will live-reload the server
type Live mg.Namespace

func (Live) All() {
	mg.Deps(
		Live.Server,
	)
}

func (Live) Server() error {
	t := target.New("go")
	t.Args([]string{
		"tool",
		"air",
		"-build.cmd", fmt.Sprintf("go build -tags=dev -o %s ./cmd/web", devOutput),
		"-build.entrypoint", devOutput,
		"-build.include_dir", "cmd,pkg",
		"-build.include_ext", "go",
	}...)

	return r.RunV("live server", t)
}
