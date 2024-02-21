#!/usr/bin/env bash
set -eo pipefail

if ! [ -x "$(command -v psql)" ]; then
    echo >&2 "Error: psql is not installed."
    exit 1
fi

ROOT_DIR=$(pwd)
source "$ROOT_DIR/env.sh"

PGPASSWORD="${DB_PASSWORD}"
export PGPASSWORD

psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d ${DB_NAME}
