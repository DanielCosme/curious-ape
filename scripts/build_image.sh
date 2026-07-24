#!/usr/bin/env sh

set -e
set -o pipefail

VERSION=$(mage version)
echo Building version: $VERSION

docker build \
  --tag curious-ape \
  --target ape \
  .
