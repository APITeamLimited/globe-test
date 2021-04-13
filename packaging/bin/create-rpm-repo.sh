#!/bin/bash
set -eEuo pipefail

# External dependencies:
# - https://github.com/rpm-software-management/createrepo
# - https://github.com/s3tools/s3cmd
#   s3cmd expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set in the
#   environment.
# - https://gnupg.org/
#   For signing the script expects the private signing key to already be
#   imported and the `rpm` command configured for signing, e.g. ~/.rpmmacros
#   should exist.
# - generate_index.py
#   For generating the index.html of each directory. It's available in the
#   packaging/bin directory of the k6 repo, and should be in $PATH.

_s3bucket="$***REMOVED***S3_BUCKET-dl.k6.io***REMOVED***"
_usage="Usage: $0 <pkgdir> <repodir> [s3bucket=$***REMOVED***_s3bucket***REMOVED***]"
PKGDIR="$***REMOVED***1?$***REMOVED***_usage***REMOVED******REMOVED***"  # The directory where .rpm files are located
REPODIR="$***REMOVED***2?$***REMOVED***_usage***REMOVED******REMOVED***" # The package repository working directory
S3PATH="$***REMOVED***3-$***REMOVED***_s3bucket***REMOVED******REMOVED***/rpm"
# Remove packages older than N number of days (730 is roughly ~2 years).
REMOVE_PKG_DAYS=730

log() ***REMOVED***
  echo "$(date -Iseconds) $*"
***REMOVED***

delete_old_pkgs() ***REMOVED***
  find "$1" -name '*.rpm' -type f -daystart -mtime "+$***REMOVED***REMOVE_PKG_DAYS***REMOVED***" -print0 | xargs -r0 rm -v
***REMOVED***

sync_to_s3() ***REMOVED***
  log "Syncing to S3 ..."
  s3cmd sync --delete-removed "$***REMOVED***REPODIR***REMOVED***/" "s3://$***REMOVED***S3PATH***REMOVED***/"

  # Set a short cache expiration for index and repo metadata files.
  s3cmd modify --recursive --exclude='*' \
    --include='index.html' --include='/repodata/*' \
    --add-header='Cache-Control: max-age=60,must-revalidate' "s3://$***REMOVED***S3PATH***REMOVED***/"
***REMOVED***

architectures="x86_64"

pushd . > /dev/null
mkdir -p "$REPODIR" && cd "$_"

for arch in $architectures; do
  mkdir -p "$arch" && cd "$_"

  # Download existing packages
  # TODO: Consider doing this over the CDN with curl to avoid S3 egress costs,
  # but that would involve parsing the index.html, checking the checksum
  # manually, etc.
  s3cmd sync --exclude='*' --include='*.rpm' "s3://$***REMOVED***S3PATH***REMOVED***/$***REMOVED***arch***REMOVED***/" ./

  # Copy the new packages in and generate signatures
  # FIXME: The architecture naming used by yum docs and in public RPM repos is
  # "x86_64", whereas our packages are named with "amd64". So we do a replacement
  # here, but we should probably consider naming them with "x86_64" instead.
  find "$PKGDIR" -name "*$***REMOVED***arch/x86_64/amd64***REMOVED****.rpm" -type f -print0 | while read -r -d $'\0' f; do
    cp -av "$f" "$PWD/"
    rpm --addsign "$***REMOVED***f##*/***REMOVED***"
  done
  createrepo .
  cd -

  delete_old_pkgs "$arch"
done

log "Generating index.html ..."
generate_index.py -r

popd > /dev/null

sync_to_s3
