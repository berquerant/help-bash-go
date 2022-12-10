#!/bin/bash

set -ex

readonly root="$(git rev-parse --show-toplevel)"
cd "$root"
find . -name "*_generated.go" -type f -delete
go generate ./...
