#!/usr/bin/env fish

set root_dir pwd
set cur_dir (realpath (dirname (status --current-filename)))
set current_branch (git branch --show-current)

if not test $current_branch = "main"
  echo "git is NOT on the main branch"
  echo "current branch: " $current_branch
  exit 1
end

# SEMVER enums
#   - alpha
#   - beta
#   - rc (release candidate)
#   - release
#   - minor
#   - mayor

set -gx new_version (semver get release); or exit 1

# Run tests.
$cur_dir/test.sh; or exit 1
$cur_dir/audit.sh; or exit 1

git diff --exit-code; or echo "Working tree cannot be dirty" and exit 1

git tag $new_version
git push || or exit 1
git push origin $new_version || exit 1

set -Ux APE_VERSION $new_version
echo $new_version
echo $APE_VERSION

echo "--- Success ---"
