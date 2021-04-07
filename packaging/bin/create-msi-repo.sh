#!/bin/bash
set -eEuo pipefail

# External dependencies:
# - https://github.com/s3tools/s3cmd
#   s3cmd expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set in the
#   environment.
# - generate_index.py
#   For generating the index.html of each directory. It's available in the
#   packaging/bin directory of the k6 repo, and should be in $PATH.

_s3bucket="$***REMOVED***S3_BUCKET-dl.k6.io***REMOVED***"
_usage="Usage: $0 <pkgdir> <repodir> [s3bucket=$***REMOVED***_s3bucket***REMOVED***]"
PKGDIR="$***REMOVED***1?$***REMOVED***_usage***REMOVED******REMOVED***"  # The directory where .msi files are located
REPODIR="$***REMOVED***2?$***REMOVED***_usage***REMOVED******REMOVED***" # The package repository working directory
S3PATH="$***REMOVED***3-$***REMOVED***_s3bucket***REMOVED******REMOVED***/msi"
# Remove packages older than N number of days (730 is roughly ~2 years).
REMOVE_PKG_DAYS=730

log() ***REMOVED***
  echo "$(date -Iseconds) $*"
***REMOVED***

delete_old_pkgs() ***REMOVED***
  find "$1" -name '*.msi' -type f -daystart -mtime "+$***REMOVED***REMOVE_PKG_DAYS***REMOVED***" -print0 | xargs -r0 rm -v
***REMOVED***

sync_to_s3() ***REMOVED***
  log "Syncing to S3 ..."
  s3cmd sync --delete-removed "$***REMOVED***REPODIR***REMOVED***/" "s3://$***REMOVED***S3PATH***REMOVED***/"

  # Disable cache for index files.
  s3cmd modify --recursive --exclude='*' --include='index.html' \
    --add-header='Cache-Control:no-cache, max-age=0' "s3://$***REMOVED***S3PATH***REMOVED***/"

  # Update latest redirect
  latest="$(find "$REPODIR" -name '*.msi' -printf '%P\n' | sort | tail -1)"
  s3cmd --force --add-header="x-amz-website-redirect-location:/msi/$***REMOVED***latest***REMOVED***" \
     --add-header='Cache-Control:no-cache, max-age=0' \
     put "$(mktemp)" "s3://$***REMOVED***S3PATH***REMOVED***/k6-latest-amd64.msi"
***REMOVED***

mkdir -p "$REPODIR"

# Download existing packages
# For MSI packages this is only done to be able to generate the index.html correctly.
# Should we fake it and create empty files that have the same timestamp and size as the original ones?
s3cmd sync --exclude='*' --include='*.msi' "s3://$***REMOVED***S3PATH***REMOVED***/" "$REPODIR/"

# Copy the new packages in
find "$PKGDIR" -name "*.msi" -type f -print0 | xargs -r0 cp -t "$REPODIR"

delete_old_pkgs "$REPODIR"

log "Generating index.html ..."
(cd "$REPODIR" && generate_index.py -r)

sync_to_s3
