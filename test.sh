#!/bin/bash
set -e
export GIN_MODE=release

for d in $(go list ./... | grep -v vendor); do
    go test -v $d
done
