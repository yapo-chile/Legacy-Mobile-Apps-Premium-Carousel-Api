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
    github.com/golangci/golangci-lint/cmd/golangci-lint
)

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
   GO111MODULE=off go get -u -v ${tool}
done

echoTitle "Removing outdated vendor"
rm -rf vendor glide.* go.*

echoTitle "Initializating go modules"
GO111MODULE=on go mod init $APPMODULE

echoTitle "Installing project dependencies"
GO111MODULE=on go mod tidy

set +e
