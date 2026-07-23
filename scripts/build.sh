#!/bin/sh

set -u

echo "Building production (static) binary into $PROD_OUTPUT"
go build \
	-buildvcs=true \
  -ldflags="-s -extldflags=-static" \
  -tags=prod \
  -o=$PROD_OUTPUT ./cmd/web
