#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running dependencies script"

set -e
# List of tools used for testing, validation, and report generation
tools=(
    github.com/jstemmer/go-junit-report
    github.com/axw/gocov/gocov
    github.com/AlekSi/gocov-xml
    github.com/Masterminds/glide
    github.com/golangci/golangci-lint/cmd/golangci-lint
)

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
    go get -u -v ${tool}
done

echoTitle "Installing Glide dependencies"
glide install

set +e
