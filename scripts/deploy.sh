#!/usr/bin/env fish

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
