#!/usr/bin/env fish

migrate create -seq -ext=.sql -dir=./migrations/sqlite $argv[1]
