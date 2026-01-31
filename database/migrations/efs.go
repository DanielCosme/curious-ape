package migrations

import "embed"

//go:embed sqlite/*
var Migrations embed.FS
