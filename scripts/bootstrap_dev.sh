#!/usr/bin/env fish

go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/rakyll/gotest@latest
go install github.com/stephenafamo/bob/gen/bobgen-sqlite@v0.30.0
go install honnef.co/go/tools/cmd/staticcheck@latest
