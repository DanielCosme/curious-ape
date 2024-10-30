#!/usr/bin/env fish

TEMPL_EXPERIMENT=rawgo templ generate -path ./web/view; or exit 1
mv ./web/view/*templ.go* ./internal/view/
