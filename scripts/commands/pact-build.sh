#!/usr/bin/env bash

 # Include colors.sh
 DIR="${BASH_SOURCE%/*}"
 if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
 . "$DIR/colors.sh"

 echoTitle "Building code"
 set -e
 echo "binary: ${PACT_BINARY} main: ${PACT_MAIN_FILE}"
 go build -v -o ${PACT_BINARY} ./${PACT_MAIN_FILE}

 set +e
 echoTitle "Done"
 