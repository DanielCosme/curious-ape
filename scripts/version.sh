#!/usr/bin/env fish

echo (semver get $argv[1]) > VERSION.txt
cat VERSION.txt
