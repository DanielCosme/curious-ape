#!/usr/bin/env sh

set -e
set -o pipefail

APE_VERSION=$(mage version_image)
REGISTRY=$(mage registry)
REGISTRY_USER=daniel

echo Version: $APE_VERSION
echo Registry $REGISTRY
docker tag curious-ape $REGISTRY/daniel/curious-ape:latest
docker tag curious-ape $REGISTRY/daniel/curious-ape:$APE_VERSION

docker login https://$REGISTRY -u $REGISTRY_USER --password-stdin <<< "$REGISTRY_PASSWORD"

docker push $REGISTRY/daniel/curious-ape:latest
docker push $REGISTRY/daniel/curious-ape:$APE_VERSION
