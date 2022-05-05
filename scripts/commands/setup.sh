#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running dependencies script"
# List of tools used for testing, validation, and report generation
tools=(
    github.com/jstemmer/go-junit-report
    github.com/axw/gocov/gocov
    github.com/AlekSi/gocov-xml
)

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
   GO111MODULE=off go get -u -v ${tool}
done

GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint

 # install librdkafka
if [ "$(id -nu)" != "root" ]; then
    echoTitle "Installing librdkafka"
    git clone https://github.com/edenhill/librdkafka.git
    cd librdkafka
    ./configure --prefix /usr
    make
    sudo make install
    cd ..
    rm -rf librdkafka
fi

echoTitle "Installing project dependencies"
GO111MODULE=on go mod tidy

set +e
