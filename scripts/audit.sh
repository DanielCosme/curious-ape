#!/usr/bin/env fish

	echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	go fmt ./...
	echo 'Vetting code...'
	go vet ./...
	staticcheck ./...