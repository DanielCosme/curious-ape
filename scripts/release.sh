#!/usr/bin/env fish

set root_dir pwd
set cur_dir (realpath (dirname (status --current-filename)))

if test -z $argv[1]
    echo "Ape version cannot be empty"
    exit 1
end

set -gx APE_VERSION (cat VERSION.txt); or exit

$cur_dir/build.sh $APE_VERSION; or exit 1

docker run curious-ape-ci; or exit 1

echo "$DOCKER_HUB_PASSWORD" | docker login -u $DOCKER_HUB_USER --password-stdin; or exit 1

docker image tag curious-ape "danielcosme/curious-ape:latest"; or exit 1
docker image tag curious-ape "danielcosme/curious-ape:$new_version"; or exit 1
docker image tag migrate-ape "danielcosme/migrate-ape:latest"; or exit 1
docker image tag migrate-ape "danielcosme/migrate-ape:$new_version"; or exit 1
docker push --all-tags danielcosme/curious-ape
docker push --all-tags danielcosme/migrate-ape

echo "New version:" $APE_VERSION
