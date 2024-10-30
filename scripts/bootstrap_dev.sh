#!/usr/bin/env fish

go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/maykonlsf/semver-cli/cmd/semver@latest
go install github.com/rakyll/gotest@latest
go install github.com/stephenafamo/bob/gen/bobgen-sqlite@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/a-h/templ/cmd/templ@latest
