#!/bin/bash
set -eEuo pipefail

# Generate repositories for k6 packages and sync them from/to AWS S3.
# These scripts use both s3cmd and aws-cli, since the latter doesn't
# preserve timestamps and doesn't compute/check hashes, relying instead
# on timestamps or file size to determine whether a sync should be done.
# aws-cli is still used for invalidating the CloudFront cache, which
# doesn't work with s3cmd...
# See:
# - https://github.com/aws/aws-cli/issues/3069#issuecomment-818732309
# - https://github.com/s3tools/s3cmd/issues/536
# - https://github.com/s3tools/s3cmd/issues/790

log() ***REMOVED***
  echo "$(date -Iseconds) $*"
***REMOVED***

signkeypath="$PWD/sign-key.gpg"
s3bucket="$***REMOVED***S3_BUCKET-dl.k6.io***REMOVED***"
pkgdir="$PWD/Packages"

if ! [ -r "$signkeypath" ]; then
  log "ERROR: Signing key not found at '$signkeypath'"
  exit 1
fi

gpg2 --import --batch --passphrase="$PGP_SIGN_KEY_PASSPHRASE" "$signkeypath"
export PGPKEYID="$(gpg2 --list-secret-keys --with-colons | grep '^sec' | cut -d: -f5)"
# Export and sync the GPG pub key if it doesn't exist in S3 already.
mkdir -p "$pkgdir"
s3cmd get "s3://$***REMOVED***s3bucket***REMOVED***/key.gpg" "$***REMOVED***pkgdir***REMOVED***/key.gpg" || ***REMOVED***
  gpg2 --export --armor --output "$***REMOVED***pkgdir***REMOVED***/key.gpg" "$PGPKEYID"
  s3cmd put "$***REMOVED***pkgdir***REMOVED***/key.gpg" "s3://$***REMOVED***s3bucket***REMOVED***/key.gpg"
***REMOVED***

for repo in deb rpm msi; do
  log "Creating $***REMOVED***repo***REMOVED*** repository ..."
  "create-$***REMOVED***repo***REMOVED***-repo.sh" "$PWD/dist" "$***REMOVED***pkgdir***REMOVED***/$***REMOVED***repo***REMOVED***"
done

# Generate and sync the main index.html
(cd "$pkgdir" && generate_index.py)
s3cmd put --add-header='Cache-Control: max-age=60,must-revalidate' \
  "$***REMOVED***pkgdir***REMOVED***/index.html" "s3://$***REMOVED***s3bucket***REMOVED***/index.html"

# Invalidate CloudFront cache for index files, repo metadata and the latest MSI package.
IFS=' ' read -ra indexes <<< \
  "$(find "$***REMOVED***pkgdir***REMOVED***" -name 'index.html' -type f | sed "s:^$***REMOVED***pkgdir***REMOVED***::" | sort | paste -sd' ')"
aws cloudfront create-invalidation --distribution-id "$AWS_CF_DISTRIBUTION" \
  --paths "$***REMOVED***indexes[@]***REMOVED***" "/msi/k6-latest-amd64.msi" \
  "/deb/dists/stable/"***REMOVED***Release,Release.gpg,InRelease***REMOVED*** \
  "/deb/dists/stable/main/binary-amd64"/Packages***REMOVED***,.gz,.bz2***REMOVED*** \
  "/rpm/x86_64/repodata/*"

exec "$@"
