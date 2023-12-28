#!/usr/bin/env bash
set -eo pipefail

DIR=$(dirname "$(readlink -f "$0")")
ROOT_DIR=$(pwd)

source "$ROOT_DIR/tmp/env.sh"
source "$DIR/functions.sh"

echo "--- Deploying to $SERVER ---"

echo "Sending configuration file"
do_scp "$ROOT_DIR/tmp/prod.env.json" "$APP_HOME/.$APP_NAME/server/prod.env.json"

do_ssh 'bash -s' <<-STDIN
  set -eo pipefail

  rm $APP_HOME/$APP_NAME
  cd $APP_HOME/repo/curious-ape
  git pull

  echo ""
  make build/web/linux
  mv ./bin/$APP_NAME $APP_HOME/$APP_NAME

  echo ""
  echo Running migrations...
  migrate -path ./migrations/sqlite/ -database sqlite3://$APP_HOME/.$APP_NAME/server/$APP_NAME.db up
  sudo systemctl restart $APP_NAME
STDIN

echo "Done!"