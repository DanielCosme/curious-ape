#!/usr/bin/env fish

semver up $argv[1]; or exit
echo (semver get $argv[1]) > VERSION.txt