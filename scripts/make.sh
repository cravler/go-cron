#!/bin/bash

SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

GOOS=linux GOARCH=amd64 "${SCRIPT_DIR}/build.sh" $@
GOOS=linux GOARCH=386 "${SCRIPT_DIR}/build.sh" $@
GOOS=linux GOARCH=arm64 "${SCRIPT_DIR}/build.sh" $@
GOOS=linux GOARCH=arm "${SCRIPT_DIR}/build.sh" $@

GOOS=darwin GOARCH=amd64 "${SCRIPT_DIR}/build.sh" $@
GOOS=darwin GOARCH=386 "${SCRIPT_DIR}/build.sh" $@

GOOS=windows GOARCH=amd64 "${SCRIPT_DIR}/build.sh" $@
GOOS=windows GOARCH=386 "${SCRIPT_DIR}/build.sh" $@