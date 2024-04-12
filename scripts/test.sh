#!/usr/bin/env fish

if type -q gotest
  gotest ./...
else
  go test ./...
end
