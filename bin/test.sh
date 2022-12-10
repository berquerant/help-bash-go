#!/bin/bash

# test.sh runs unit tests and golden tests.
# $1: help-bash-go executable (absolute path)

set -ex

readonly bin="$1"
if ! test -f "$bin" ; then
    echo "${bin} is not a file" >&2
    exit 1
fi
readonly root="$(git rev-parse --show-toplevel)"

cd "$root"
echo "Run unit tests..."
go test -cover ./...

echo "Run golden tests..."
cd "${root}/help-bash" # enter submodule
./run-test.sh "$bin"
