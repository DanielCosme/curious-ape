#!/usr/bin/env fish

set root_dir pwd
set cur_dir (realpath (dirname (status --current-filename)))
set current_branch (git branch --show-current)

if not test $current_branch = "main"
  echo "git is NOT on the main branch"
  echo "current branch: " $current_branch
  exit 1
end

set -gx new_version (semver get release); or exit 1

# Run tests.
$cur_dir/test.sh; or exit 1

# SEMVER enums
#   - alpha
#   - beta
#   - rc (release candidate)
#   - release
#   - minor
#   - mayor

git tag $new_version
git push origin $new_version || exit 1
git push || or exit 1

$cur_dir/build.sh $new_version; or exit 1

echo "$DOCKER_HUB_PASSWORD" | docker login -u $DOCKER_HUB_USER --password-stdin; or exit 1

docker image tag curious-ape "danielcosme/curious-ape:latest"; or exit 1
docker image tag curious-ape "danielcosme/curious-ape:$new_version"; or exit 1
docker image tag migrate-ape "danielcosme/migrate-ape:latest"; or exit 1
docker image tag migrate-ape "danielcosme/migrate-ape:$new_version"; or exit 1
docker push --all-tags danielcosme/curious-ape
docker push --all-tags danielcosme/migrate-ape

echo "New version:" $new_version
