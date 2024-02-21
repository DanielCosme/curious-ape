#!/usr/bin/env bash
set -eo pipefail

DIR=$(dirname "$(readlink -f "$0")")
ROOT_DIR=$(pwd)

source "$ROOT_DIR/env.sh"
set -x
# source "$DIR/functions.sh"

cargo test
cargo clippy -- -D warnings
cargo fmt -- --check

DB_HOST=$SERVER

# cargo sqlx prepare --workspace

docker build --tag danielcosme/curious-ape --file Dockerfile $ROOT_DIR

docker context use $REMOTE_DOCKER_CONTEXT

docker-compose up --detach --wait --build

source "$ROOT_DIR/.env"
DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
export DATABASE_URL

sqlx database create
sqlx migrate run

docker context use default
