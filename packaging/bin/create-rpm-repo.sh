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

_s3bucket="$***REMOVED***S3_BUCKET-dl-k6-io***REMOVED***"
_usage="Usage: $0 <pkgdir> <repodir> [s3bucket=$***REMOVED***_s3bucket***REMOVED***]"
PKGDIR="$***REMOVED***1?$***REMOVED***_usage***REMOVED******REMOVED***"  # The directory where .rpm files are located
REPODIR="$***REMOVED***2?$***REMOVED***_usage***REMOVED******REMOVED***" # The package repository working directory
S3PATH="$***REMOVED***3-$***REMOVED***_s3bucket***REMOVED******REMOVED***/rpm"

architectures="x86_64"

# TODO: Remove old package versions?
# Something like: https://github.com/kopia/kopia/blob/master/tools/apt-publish.sh#L23-L25

mkdir -p "$REPODIR" && cd "$_"

for arch in $architectures; do
  mkdir -p "$arch" && cd "$_"

  # Download existing packages
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
done
