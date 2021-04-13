#!/bin/bash
set -eEuo pipefail

# External dependencies:
# - https://aws.amazon.com/cli/
#   awscli expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set in the
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
  aws s3 sync --no-progress --delete "$***REMOVED***REPODIR***REMOVED***/" "s3://$***REMOVED***S3PATH***REMOVED***/"

  # Set a short cache expiration for index files and the latest MSI package.
  aws s3 cp --no-progress --recursive --exclude='*' \
    --include='*.html' \
    --cache-control='max-age=60,must-revalidate' \
    --content-type='text/html' \
    --metadata-directive=REPLACE \
    "s3://$***REMOVED***S3PATH***REMOVED***" "s3://$***REMOVED***S3PATH***REMOVED***"
  aws s3 cp --no-progress \
    --cache-control='max-age=60,must-revalidate' \
    --content-type='application/x-msi' \
    --metadata-directive=REPLACE \
    "s3://$***REMOVED***S3PATH***REMOVED***/k6-latest-amd64.msi" "s3://$***REMOVED***S3PATH***REMOVED***/k6-latest-amd64.msi"
***REMOVED***

mkdir -p "$REPODIR"

# Download existing packages
# For MSI packages this is only done to be able to generate the index.html correctly.
# Should we fake it and create empty files that have the same timestamp and size as the original ones?
aws s3 sync --no-progress --exclude='*' --include='*.msi' "s3://$***REMOVED***S3PATH***REMOVED***/" "$REPODIR/"

# Copy the new packages in
find "$PKGDIR" -name "*.msi" -type f -print0 | xargs -r0 cp -avt "$REPODIR"

delete_old_pkgs "$REPODIR"

# Update the latest package. This could be done with S3 redirects, but
# CloudFront caches redirects aggressively and I wasn't able to invalidate it.
latest="$(find "$REPODIR" -name '*.msi' -printf '%P\n' | sort | tail -1)"
cp -p "$***REMOVED***REPODIR***REMOVED***/$***REMOVED***latest***REMOVED***" "$***REMOVED***REPODIR***REMOVED***/k6-latest-amd64.msi"

log "Generating index.html ..."
(cd "$REPODIR" && generate_index.py -r)

sync_to_s3
