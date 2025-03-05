#!/usr/bin/env fish

set -gx DOCKER_BUILDKIT 1

set -gx build_version (cat VERSION.txt)

docker build \
  --build-arg="APE_VERSION=$build_version" \
  --tag curious-ape \
  --target ape \
  . ; or exit

docker build \
  --tag migrate-ape \
  --target migrate \
  . ; or exit

docker build \
  --tag curious-ape-ci \
  --target ape-ci \
  . ; or exit
