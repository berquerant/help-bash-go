#!/bin/bash

# build.sh builds bash-help executable binary.
# This will not change the code.
# $1: output absolute path

set -ex

readonly output="$1"
if test -z "$output" ; then
    echo Output required >&2
    exit 1
fi
readonly root="$(git rev-parse --show-toplevel)"
readonly version="$(git describe --tags --abbrev=0)"
readonly commit="$(git rev-parse HEAD)"
readonly goversion="$(go version)"

gen_ldflags() {
    readonly ldflags_package="main"
    readonly ldflags_array=(
        "AuthorName=berquerant"
        "Version=${version}"
        "GoVersion=${goversion}"
        "Commit=${commit}"
        "Project=help-bash-go"
        "GithubUser=berquerant"
    )
    for x in "${ldflags_array[@]}" ; do echo "-X '${ldflags_package}.${x}'" ; done | tr "\n" " "
}
readonly ldflags="-ldflags=$(gen_ldflags)"

cd "$root"
go get .
go build -v -o "$output" "$ldflags"
"$output" -v
