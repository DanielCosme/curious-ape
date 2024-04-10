#! /usr/bin/env fish


# docker run \
#   -v ./migrations/sqlite:/migrations \
#   -v ./ape.db:/ape.db \
#   danielcosme/migrate-ape \
#   -path=/migrations/ \
#   -database "sqlite3://./ape.db" \
#   $argv --all

migrate -path "./migrations/sqlite" -database "sqlite3://./ape.db" $argv
