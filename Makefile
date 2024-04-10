# include .envrc
DB_PATH=${HOME}/.ape/server/ape.db

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

# ==================================================================================== #
# BUILD
# ==================================================================================== #

#current_time = $(shell date --iso-8601=seconds)
#git_description = $(shell git describe --always --dirty --tags --long)
#linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'
# -extldflags "-static"

## build: Builds container images for the project.
.PHONY: build
build:
	@echo 'Building cmd/web...'
	# @echo ${current_time}
	./scripts/build.sh
	# CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/ape ./cmd/web

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## test: test all code
.PHONY: test
test:
	./scripts/test.sh

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	# staticcheck ./...
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
