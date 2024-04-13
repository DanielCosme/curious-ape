#!/usr/bin/env fish

# docker run \
#  -v ./migrations/sqlite:/migrations \
#  -v $(pwd)/ape.db:/ape.db \
#  danielcosme/migrate-ape \
#  -path=/migrations/ \
#  -database "sqlite3://./ape.db" \
#  $argv

echo "Running migrations..."
migrate -path "./migrations/sqlite_backup" -database "sqlite3://./ape.db" $argv
