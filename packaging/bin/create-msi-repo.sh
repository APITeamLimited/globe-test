#!/bin/bash
set -eEuo pipefail

# External dependencies:
# - https://github.com/s3tools/s3cmd
#   s3cmd expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set in the
#   environment.

_s3bucket="$***REMOVED***S3_BUCKET-dl-k6-io***REMOVED***"
_usage="Usage: $0 <pkgdir> <repodir> [s3bucket=$***REMOVED***_s3bucket***REMOVED***]"
PKGDIR="$***REMOVED***1?$***REMOVED***_usage***REMOVED******REMOVED***"  # The directory where .msi files are located
REPODIR="$***REMOVED***2?$***REMOVED***_usage***REMOVED******REMOVED***" # The package repository working directory
S3PATH="$***REMOVED***3-$***REMOVED***_s3bucket***REMOVED******REMOVED***/msi"
# Number of packages to keep, older packages will be deleted.
RETAIN_PKG_COUNT=25

delete_old_pkgs() ***REMOVED***
  find "$1" -name '*.msi' -type f | sort -r | tail -n "+$((RETAIN_PKG_COUNT+1))" | xargs -r rm -v
***REMOVED***

mkdir -p "$REPODIR"

# Download existing packages
# For MSI packages this is only done to be able to generate the index.html correctly.
# Should we fake it and create empty files that have the same timestamp and size as the original ones?
s3cmd sync --exclude='*' --include='*.msi' "s3://$***REMOVED***S3PATH***REMOVED***/" "$REPODIR/"

# Copy the new packages in
find "$PKGDIR" -name "*.msi" -type f -print0 | xargs -r0 cp -t "$REPODIR"

delete_old_pkgs "$REPODIR"
