#!/usr/bin/env sh

set -e
set -o pipefail

APE_VERSION=$(mage version)
echo Version: $APE_VERSION
REGISTRY=$(mage registry)
echo Registry $REGISTRY

podman build \
  --tag curious-ape \
  --target ape \
  --build-arg="APE_VERSION=$APE_VERSION" \
  .

podman tag curious-ape $REGISTRY/daniel/curious-ape:latest
podman tag curious-ape $REGISTRY/daniel/curious-ape:$APE_VERSION

podman login --get-login https://$REGISTRY ||
  podman login https://$REGISTRY -u $REGISTRY_USER --password-stdin <<< "$REGISTRY_PASSWORD"

podman push $REGISTRY/daniel/curious-ape:latest
podman push $REGISTRY/daniel/curious-ape:$APE_VERSION
