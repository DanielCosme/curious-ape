#!/usr/bin/env fish

set -gx ape_version (mage version); or exit

if test $argv[1] && test $argv[1] = prod
    echo "Building production (static) binary into $PROD_OUTPUT"
    go build \
        -ldflags="-s -extldflags=-static -X main.version=$ape_version" \
        -o=$PROD_OUTPUT ./cmd/web
else
    echo "Building dev binary into $DEV_OUTPUT"
    go build \
        -ldflags="-X main.version=$ape_version-dev" \
        -o=$DEV_OUTPUT ./cmd/web
end
