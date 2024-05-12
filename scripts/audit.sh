#!/usr/bin/env fish

	echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	echo 'Running tests...'
	go test -race -vet=off ./...
