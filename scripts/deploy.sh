#!/usr/bin/env fish

if test "$RELEASE" = true
  echo "--- Pushing code to remote ---"
  git push; or exit 1
  echo "--- Success ---"
  echo ""

  echo "--- Starting Release ---"
  echo "\
    cd curious-ape
    git pull
    ./scripts/release.sh; or exit 1 \
    " | ssh daniel@prime ; or exit 1
    echo "--- Success ---"
  echo ""
end


echo "--- Synchornizing deployment files ---"
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
    docker compose up -d &&
    docker system prune -f; or exit 1
    " | ssh daniel@danicos.me ; or exit 1
echo "--- Success ---"
echo ""

echo "--- Done! ---"
