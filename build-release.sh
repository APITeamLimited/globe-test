#!/bin/bash

set -eEuo pipefail

eval "$(go env)"

OUT_DIR="$***REMOVED***1-dist***REMOVED***"
# To override the latest git tag as the version, pass something else as the second arg.
VERSION=$***REMOVED***2:-$(git describe --tags --always --dirty)***REMOVED***

# To overwrite the version details, pass something as the third arg. Empty string disables it.
VERSION_DETAILS=$***REMOVED***3-"$(date -u +"%FT%T%z")/$(git describe --always --long --dirty)"***REMOVED***

build() ***REMOVED***
    local ALIAS="$1" SUFFIX="$***REMOVED***2***REMOVED***"  # Any other arguments are passed to the go build command as env vars
    local NAME="k6-$***REMOVED***VERSION***REMOVED***-$***REMOVED***ALIAS***REMOVED***"

    local BUILD_ENV=("$***REMOVED***@:3***REMOVED***")
    local BUILD_ARGS=(-o "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***/k6$***REMOVED***SUFFIX***REMOVED***" -trimpath)

    if [ -n "$VERSION_DETAILS" ]; then
        BUILD_ARGS+=(-ldflags "-X go.k6.io/k6/lib/consts.VersionDetails=$VERSION_DETAILS")
    fi

    echo "- Building platform: $***REMOVED***ALIAS***REMOVED*** (" "$***REMOVED***BUILD_ENV[@]***REMOVED***" "go build" "$***REMOVED***BUILD_ARGS[@]***REMOVED***" ")"

    mkdir -p "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***"

    # Subshell to not mess with the current env vars or CWD
    (
        export "$***REMOVED***BUILD_ENV[@]***REMOVED***"
        # Build a binary
         go build "$***REMOVED***BUILD_ARGS[@]***REMOVED***"
    )
***REMOVED***

package() ***REMOVED***
    local ALIAS="$1" FMT="$2"
    local NAME="k6-$***REMOVED***VERSION***REMOVED***-$***REMOVED***ALIAS***REMOVED***"
    echo "- Creating $***REMOVED***NAME***REMOVED***.$***REMOVED***FMT***REMOVED*** package..."
    case $FMT in
    deb|rpm)
        # The go-bin-* tools expect the binary in /tmp/
        [ ! -r /tmp/k6 ] && cp "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***/k6" /tmp/k6
        "go-bin-$***REMOVED***FMT***REMOVED***" generate --file "packaging/$***REMOVED***FMT***REMOVED***.json" -a amd64 \
            --version "$***REMOVED***VERSION#v***REMOVED***" -o "$***REMOVED***OUT_DIR***REMOVED***/k6-$***REMOVED***VERSION***REMOVED***-amd64.$***REMOVED***FMT***REMOVED***"
        ;;
    tgz)
        tar -C "$***REMOVED***OUT_DIR***REMOVED***" -zcf "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***.tar.gz" "$NAME"
        ;;
    zip)
        (cd "$***REMOVED***OUT_DIR***REMOVED***" && zip -rq9 - "$NAME") > "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***.zip"
        ;;
    *)
        echo "Unknown format: $FMT"
        return 1
        ;;
    esac
***REMOVED***

cleanup() ***REMOVED***
    find "$OUT_DIR" -mindepth 1 -maxdepth 1 -type d -exec rm -rf ***REMOVED******REMOVED*** \;
    echo "--- Cleaned $***REMOVED***OUT_DIR***REMOVED***"
***REMOVED***
trap cleanup EXIT

echo "--- Building Release: $***REMOVED***VERSION***REMOVED***"

mkdir -p "$OUT_DIR"

build macos-amd64   ""   GOOS=darwin  GOARCH=amd64
build macos-arm64   ""   GOOS=darwin  GOARCH=arm64
build windows-amd64 .exe GOOS=windows GOARCH=amd64
build linux-amd64   ""   GOOS=linux   GOARCH=amd64  CGO_ENABLED=0
build linux-arm64   ""   GOOS=linux   GOARCH=arm64  CGO_ENABLED=0

package macos-amd64   zip
package macos-arm64   zip
package windows-amd64 zip
package linux-amd64   tgz
package linux-arm64   tgz
package linux-amd64   rpm
package linux-amd64   deb
