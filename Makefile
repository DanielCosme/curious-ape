include .envrc

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

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${APE_PG_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${APE_PG_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${APE_PG_DB_DSN} up
## db/migrations/up: apply all up database migrations

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description} -extldflags "-static"'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	@echo $current_time
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/ape ./cmd/api

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	#@echo 'Running tests...'
	#go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = "104.248.10.77"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh daniel@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -rP --delete ./bin/ape ./migrations daniel@${production_host_ip}:~
	ssh -t daniel@${production_host_ip} 'migrate -path ~/migrations -database $$APE_DB_DSN up'

## production/configure/ape.service: configure the production systemd ape.service file
.PHONY: production/configure/api.service
production/configure/ape.service:
	rsync -P ./remote/production/ape.service daniel@${production_host_ip}:~
	ssh -t daniel@${production_host_ip} '\
		sudo mv ~/ape.service /etc/systemd/system/ \
		&& sudo systemctl enable ape \
		&& sudo systemctl restart ape \
	'

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P ./remote/production/Caddyfile daniel@${production_host_ip}:~
	ssh -t daniel@${production_host_ip} '\
		sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'
