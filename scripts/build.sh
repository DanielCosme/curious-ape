#! /usr/bin/env fish

docker build \
  --tag danielcosme/curious-ape \
  --target ape \
  .

docker build \
  --tag danielcosme/migrate-ape \
  --target migrate \
  .
