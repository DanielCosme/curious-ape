if ! [ -x "$(command -v psql)" ]; then
    echo >&2 "Error: psql is not installed."
    exit 1
fi

DB_USER="${POSTGRES_USER:=postgres}"
PGPASSWORD="${POSTGRES_PASSWORD:=password}"
DB_NAME="${POSTGRES_DB:=ape}"
DB_PORT="${POSTGRES_PORT:=5432}"
DB_HOST="${POSTGRES_HOST:=localhost}"

export PGPASSWORD

psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d ${DB_NAME}

