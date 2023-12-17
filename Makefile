# include .envrc
DB_PATH=${HOME}/.ape/server/ape.db

# Go migrate for slite3 support
# go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

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
	@echo ${DB}
	go run ./cmd/httpd -env dev

## db/psql: connect to the database using psql
# 	.PHONY: db/psql
# 	db/psql:
# 		psql ${APE_PG_DB_DSN}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations/sqlite -database sqlite3://${DB_PATH} up

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations/sqlite ${name}

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'
# -extldflags "-static"

## build/api: build the cmd/api application
.PHONY: build/api/linux
build/api/linux:
	@echo 'Building cmd/api...'
	@echo ${current_time}
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/aped ./cmd/httpd

## build/cli: build the cmd/cli application
.PHONY: build/cli/linux
build/cli/linux:
	@echo 'Building cmd/cli...'
	@echo ${current_time}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/ape ./cmd/cli

## build/web: build the cmd/web application
.PHONY: build/web/linux
build/web/linux:
	@echo 'Building cmd/web...'
	@echo ${current_time}
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/ape ./cmd/web

## install/cli: install the cmd/cli application
.PHONY: install/cli/linux
install/cli/linux: build/cli/linux
	@echo 'Installing cmd/cli...'
	rm ${GOBIN}/ape
	mv ./bin/ape ${GOBIN}/ape

.PHONY: build/api/mac
build/api/mac:
	@echo 'Building cmd/api...'
	@echo $current_time
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/aped ./cmd/httpd

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

production_host_ip = "danicos.me"
ssh_key_path = "~/.ssh/do_rsa"
ssh_command = "ssh -i ${ssh_key_path}"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh -i ${ssh_key_path} daniel@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api: build/api/linux
	rsync -rP -e ${ssh_command} --delete ./bin/ape ./migrations/sqlite daniel@${production_host_ip}:/home/daniel
	ssh -i ${ssh_key_path} -t daniel@${production_host_ip} 'migrate -path ~/sqlite -database sqlite3://~/.ape/server/ape.db up \
		&& sudo systemctl restart ape \
	'

.PHONY: production/deploy/api2
production/deploy/api2: 
	ssh -t daniel@${production_host_ip} '\
		cd /home/daniel/repo/curious-ape \
		&& git pull \
		&& make build/api/linux \
		&& mv ./bin/aped /home/daniel/ape \
		&& sudo systemctl restart ape \
	'

## to be run on the remote machine
.PHONY: production/reload/api
production/reload/api: build/api/linux
	mv ./bin/ape ~/
	sudo systemctl restart ape

## production/configure/ape.service: configure the production systemd ape.service file
.PHONY: production/configure/ape.service
production/configure/ape.service:
	rsync -P -e ${ssh_command} ./production/config/ape.service daniel@${production_host_ip}:~/
	ssh -i ${ssh_key_path} -t daniel@${production_host_ip} '\
		sudo mv ~/ape.service /etc/systemd/system/ \
		&& sudo systemctl enable ape \
		&& sudo systemctl restart ape \
	'

## production/configure/caddyfile: configure the production Caddyfile
.PHONY: production/configure/caddyfile
production/configure/caddyfile:
	rsync -P -e ${ssh_command} ./remote/production/Caddyfile daniel@${production_host_ip}:~
	ssh -i ${ssh_key_path} -t daniel@${production_host_ip} '\
		sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'

# ==================================================================================== #
# CLIENT
# ==================================================================================== #

.PHONY: cli/install
cli/install:
	go install ./cmd/cli/apectl.go
