# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: Run the application
.PHONY: run
run:
	./scripts/run.sh

## db: Quickly access the database
.PHONY: db
db:
	sqlite3 -box ape.db


## watch: Watches for changes in the Web UI
.PHONY: watch
watch:
	templ generate --watch

## gen: Run all generators of the project
.PHONY: gen
gen:
	./scripts/migrate.sh up
	./scripts/gen-sql.sh
	./scripts/gen-templ.sh

## gen-sql: Generate type safe SQL helpers
.PHONY: gen-sql
gen-sql:
	./scripts/migrate.sh up
	./scripts/gen-sql.sh

## migrate-up: Run SQL migrations up
.PHONY: migrate-up
migrate-up:
	./scripts/migrate.sh up

## migrate-down: Run SQL migrations down
.PHONY: migrate-down
migrate-down:
	./scripts/migrate.sh down

## build: Builds container images for the project
.PHONY: build
build:
	./scripts/build.sh

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## ci: run all CI checks
ci: test audit

## test: test all code
.PHONY: test
test:
	./scripts/test.sh

## audit: tidy dependencies, format, vet and run static checks on all code
.PHONY: audit
audit:
	./scripts/audit.sh

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
