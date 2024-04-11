#! /usr/bin/env fish

set -gx DOCKER_BUILDKIT 1

docker build \
  --tag danielcosme/curious-ape \
  --target ape \
  .

docker build \
  --tag danielcosme/migrate-ape \
  --target migrate \
  .
