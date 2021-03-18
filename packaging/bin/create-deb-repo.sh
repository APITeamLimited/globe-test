#!/bin/bash
set -eEuo pipefail

# External dependencies:
# - https://salsa.debian.org/apt-team/apt (apt-ftparchive, packaged in apt-utils)
# - https://github.com/s3tools/s3cmd
#   s3cmd expects AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set in the
#   environment.
# - https://gnupg.org/
#   For signing the script expects the private signing key to already be
#   imported and PGPKEYID and PGP_SIGN_KEY_PASSPHRASE to be set in the
#   environment.

_s3bucket="$***REMOVED***S3_BUCKET-dl-k6-io***REMOVED***"
_usage="Usage: $0 <pkgdir> <repodir> [s3bucket=$***REMOVED***_s3bucket***REMOVED***]"
PKGDIR="$***REMOVED***1?$***REMOVED***_usage***REMOVED******REMOVED***"  # The directory where .deb files are located
REPODIR="$***REMOVED***2?$***REMOVED***_usage***REMOVED******REMOVED***" # The package repository working directory
S3PATH="$***REMOVED***3-$***REMOVED***_s3bucket***REMOVED******REMOVED***/deb"

# We don't publish i386 packages, but the repo structure is needed for
# compatibility on some systems. See https://unix.stackexchange.com/a/272916 .
architectures="amd64 i386"
# TODO: Replace with CDN URL
#repobaseurl="https://dl.k6.io/deb"
repobaseurl="http://test-dl-k6-io.s3-website.eu-north-1.amazonaws.com/deb"

# TODO: Remove old package versions?
# Something like: https://github.com/kopia/kopia/blob/master/tools/apt-publish.sh#L23-L25

mkdir -p "$REPODIR" && cd "$_"

for arch in $architectures; do
  bindir="dists/stable/main/binary-$arch"
  mkdir -p "$bindir"
  # Download existing packages via the CDN to avoid S3 egress costs.
  # An optimization might be to just append to the Packages file and upload it
  # and the new package only, but updating the index.html would get messy and
  # would be inconsistent with the RPM script which does require all packages to
  # be present because of how createrepo works.
  # TODO: Also check their hashes? Or just sync them with s3cmd which does MD5 checks...
  files=$(s3cmd ls "s3://$***REMOVED***S3PATH***REMOVED***/$***REMOVED***bindir***REMOVED***/" | ***REMOVED*** grep -oP "(?<=/$***REMOVED***S3PATH***REMOVED***/).*\.(deb|asc)" || true; ***REMOVED***)
  # curl supports parallel downloads with the -Z option since v7.68.0, but
  # unfortunately Debian carries an older version, hence xargs.
  echo "$files" | xargs -r -I***REMOVED******REMOVED*** -n1 -P"$(nproc)" curl -fsSLR -o "***REMOVED******REMOVED***" "$repobaseurl/***REMOVED******REMOVED***"

  # Copy the new packages in
  find "$PKGDIR" -name "*$arch*.deb" -type f -print0 | xargs -r0 cp -t "$bindir"
  # Generate signatures for files that don't have it
  # TODO: Switch to debsign instead? This is currently done as Bintray did it,
  # but the signature is not validated by apt/dpkg.
  # https://blog.packagecloud.io/eng/2014/10/28/howto-gpg-sign-verify-deb-packages-apt-repositories/
  find "$bindir" -type f -name '*.deb' -print0 | while read -r -d $'\0' f; do
    if ! [ -r "$***REMOVED***f***REMOVED***.asc" ]; then
      gpg2 --default-key="$PGPKEYID" --passphrase="$PGP_SIGN_KEY_PASSPHRASE" \
        --pinentry-mode=loopback --yes --detach-sign --armor -o "$***REMOVED***f***REMOVED***.asc" "$f"
    fi
  done
  apt-ftparchive packages "$bindir" | tee "$bindir/Packages"
  gzip -fk "$bindir/Packages"
  bzip2 -fk "$bindir/Packages"
done

echo "Creating release file..."
apt-ftparchive release \
  -o APT::FTPArchive::Release::Origin="k6" \
  -o APT::FTPArchive::Release::Label="k6" \
  -o APT::FTPArchive::Release::Suite="stable" \
  -o APT::FTPArchive::Release::Codename="stable" \
  -o APT::FTPArchive::Release::Architectures="$architectures" \
  -o APT::FTPArchive::Release::Components="main" \
  -o APT::FTPArchive::Release::Date="$(date -Ru)" \
  "dists/stable" > "dists/stable/Release"

# Sign release file
gpg2 --default-key="$PGPKEYID" --passphrase="$PGP_SIGN_KEY_PASSPHRASE" \
  --pinentry-mode=loopback --yes --detach-sign --armor \
  -o "dists/stable/Release.gpg" "dists/stable/Release"
gpg2 --default-key="$PGPKEYID" --passphrase="$PGP_SIGN_KEY_PASSPHRASE" \
  --pinentry-mode=loopback --yes --clear-sign \
  -o "dists/stable/InRelease" "dists/stable/Release"
