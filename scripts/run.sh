#! /usr/bin/env fish

source ./scripts/env.sh

./scripts/migrate.sh up; or exit 1
./scripts/gen-sql.sh up; or exit 1

# We do it like so because of: https://github.com/golang/go/issues/51279
go build -ldflags="-X main.version=$(cat VERSION.TXT)-dev" -o=(pwd)/bin/ape ./cmd/web; or exit 1
./bin/ape

# go run cmd/web/main.go
