#!/usr/bin/env fish

templ generate -path ./web/view
mv ./web/view/*templ.go* ./internal/view/
