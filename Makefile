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

## run: Runs the application
.PHONY: run
run:
	./scripts/run.sh

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
