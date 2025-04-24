#!/usr/bin/env sh

if command -v air >/dev/null 2>&1; then
  air
else
  go run ./cmd/web
fi
