#!/usr/bin/env fish

if test "$RELEASE" = true
  ./scripts/push.sh; or exit 1
  echo $APE_VERSION

  echo "--- Starting Release ---"
  echo "\
    cd curious-ape
    git checkout main; or exit 1
    git pull; or exit 1
    ./scripts/release.sh $APE_VERSION; or exit 1 \
    " | ssh daniel@danicos.me ; or exit 1
    echo "--- Success ---"
  echo ""
end

echo "--- Synchronizing deployment files ---"
# Transfer directory contents, but not the directory itself
rsync \
  --verbose \
  --recursive \
  ./deployment/prod/ \
  daniel@danicos.me:~/ape-deployment/ ; or exit 1
echo "--- Success ---"
echo ""

echo "--- Refreshing containers ---"
echo "\
    cd ape-deployment &&
    docker compose pull &&
    docker compose up -d; or exit 1
    " | ssh daniel@danicos.me ; or exit 1
echo "--- Success ---"
echo ""

echo "--- Done! ---"
