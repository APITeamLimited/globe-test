#!/bin/bash
set -e

if [[ -z "$***REMOVED***VERSION***REMOVED***" ]]; then
    echo "\$VERSION enviroment variable undefined"
    exit 2
fi

rm -rf pkg-build

mkdir -p dist
go-bin-deb generate --file packaging/deb.json -a amd64 --version $VERSION -o dist/k6-v$VERSION-amd64.deb
go-bin-rpm generate --file packaging/rpm.json -a amd64 --version $VERSION -o dist/k6-v$VERSION-amd64.rpm