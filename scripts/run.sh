#! /usr/bin/env fish

source ./scripts/env.sh

./scripts/migrate.sh up
go run ./cmd/web/main.go
