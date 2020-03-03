#!/bin/bash

SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
WORKDIR="$( dirname "${SCRIPT_DIR}" )"

cd "${WORKDIR}"

GOOS=${GOOS:=linux}
GOARCH=${GOARCH:=amd64}
BUILD_DIR=${BUILD_DIR:=.build}
VERSION=${VERSION:=0.x}
ARCHIVE=NO

for i in "$@"; do
case $i in
    --archive)
        ARCHIVE=YES
        shift
    ;;
    *)
        # unknown option
    ;;
esac
done

DIR="${BUILD_DIR}/${GOOS}/${GOARCH}"
mkdir -p "${DIR}"

FILE="cron"
if [ "windows" = "${GOOS}" ]; then
    FILE="cron.exe"
fi

SHA=$( git rev-parse HEAD 2>/dev/null | head -c7 )
if [ -z "${SHA}" ]; then
    SHA="dev"
fi

PACKAGE="cmd/cron/main.go"

if [ "linux" = "${GOOS}" ]; then
    CGO_ENABLED=0 go build -a -installsuffix cgo -o "${DIR}/${FILE}" -ldflags "-X main.version=${VERSION}-${SHA}" ${PACKAGE}
else
    go build -o "${DIR}/${FILE}" -ldflags "-X main.version=${VERSION}-${SHA}" ${PACKAGE}
fi

md5sum --tag "${DIR}/${FILE}" > "${DIR}/md5"

if [ "YES" = "${ARCHIVE}" ]; then
    cd "${DIR}"
    cp "${WORKDIR}/LICENSE" ./
    cp "${WORKDIR}/README.md" ./
    TAR_FILE="${WORKDIR}/${BUILD_DIR}/cron_${GOOS}_${GOARCH}.tar.gz"
    tar -czf "${TAR_FILE}" *
    md5sum --tag "${TAR_FILE}" > "${TAR_FILE}.md5"
    rm -rf "${WORKDIR}/$( dirname "${DIR}" )"
fi