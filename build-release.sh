#!/bin/bash

set -eEuo pipefail

eval "$(go env)"

set -x
export OUT_DIR="$***REMOVED***1-dist***REMOVED***"
# To override the latest git tag as the version, pass something else as the second arg.
export VERSION=$***REMOVED***2:-$(git describe --tags --always --dirty)***REMOVED***

# To overwrite the version details, pass something as the third arg. Empty string disables it.
export VERSION_DETAILS=$***REMOVED***3-"$(date -u +"%FT%T%z")/$(git describe --always --long --dirty)"***REMOVED***
set +x

build() ***REMOVED***
    local ALIAS="$1" SUFFIX="$***REMOVED***2***REMOVED***"  # Any other arguments are passed to the go build command as env vars
    local NAME="k6-$***REMOVED***VERSION***REMOVED***-$***REMOVED***ALIAS***REMOVED***"

    local BUILD_ARGS=(-o "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***/k6$***REMOVED***SUFFIX***REMOVED***" -trimpath)
    if [ -n "$VERSION_DETAILS" ]; then
        BUILD_ARGS+=(-ldflags "-X go.k6.io/k6/lib/consts.VersionDetails=$VERSION_DETAILS")
    fi

    local PACKAGE_FORMATS
    IFS="," read -ra PACKAGE_FORMATS <<< "$***REMOVED***3***REMOVED***"

    local ENV_VARS
    IFS="," read -ra ENV_VARS <<< "$***REMOVED***4***REMOVED***"

    echo "- Building platform: $***REMOVED***ALIAS***REMOVED*** (" "$***REMOVED***ENV_VARS[@]***REMOVED***" "go build" "$***REMOVED***BUILD_ARGS[@]***REMOVED***" ")"

    mkdir -p "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***"

    # Subshell to not mess with the current env vars or CWD
    (
        export "$***REMOVED***ENV_VARS[@]***REMOVED***"
        # Build a binary
        go build "$***REMOVED***BUILD_ARGS[@]***REMOVED***"

        for format in "$***REMOVED***PACKAGE_FORMATS[@]***REMOVED***"; do
            package "$format"
        done
    )
***REMOVED***

package() ***REMOVED***
    local FMT="$1"
    echo "- Creating $***REMOVED***NAME***REMOVED***.$***REMOVED***FMT***REMOVED*** package..."
    case $FMT in
    deb|rpm)
        # nfpm can't substitute env vars in file paths, so we have to cd...
        cd "$***REMOVED***OUT_DIR***REMOVED***/$***REMOVED***NAME***REMOVED***"
        set -x
        nfpm package --config ../../packaging/nfpm.yaml --packager "$***REMOVED***FMT***REMOVED***" \
            --target "../k6-$***REMOVED***VERSION***REMOVED***-linux-$***REMOVED***GOARCH***REMOVED***.$***REMOVED***FMT***REMOVED***"
        set +x
        cd -
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

build linux-amd64   ""      tgz,rpm,deb     GOOS=linux,GOARCH=amd64,CGO_ENABLED=0
build linux-arm64   ""      tgz             GOOS=linux,GOARCH=arm64,CGO_ENABLED=0 # TODO: package rpm and dep too
build macos-amd64   ""      zip             GOOS=darwin,GOARCH=amd64
build macos-arm64   ""      zip             GOOS=darwin,GOARCH=arm64
build windows-amd64 .exe    zip             GOOS=windows,GOARCH=amd64