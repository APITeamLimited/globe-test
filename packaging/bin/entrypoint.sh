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

# Setup RPM signing
cat > "$HOME/.rpmmacros" <<EOF
%_gpgbin        $(which gpg2)
%_gpg_path      $HOME/.gnupg
%_gpg_name      k6
%_gpg_pass      -
%__gpg_sign_cmd   %***REMOVED***__gpg***REMOVED*** gpg2 --default-key="$PGPKEYID" --no-verbose --no-armor --pinentry-mode=loopback --yes --passphrase="$PGP_SIGN_KEY_PASSPHRASE" --no-secmem-warning --detach-sign -o %***REMOVED***__signature_filename***REMOVED*** %***REMOVED***__plaintext_filename***REMOVED***
EOF

for repo in deb rpm msi; do
  log "Creating $***REMOVED***repo***REMOVED*** repository ..."
  "create-$***REMOVED***repo***REMOVED***-repo.sh" "$PWD/dist" "$***REMOVED***pkgdir***REMOVED***/$***REMOVED***repo***REMOVED***"
done

# Create and sync the RPM repository package if it doesn't exist already.
# NOTE: `s3cmd info` requires GetPolicy AWS permissions, and `s3cmd ls`
# exits with 0 for missing objects, so use awscli for this check.
aws s3 ls "s3://$***REMOVED***s3bucket***REMOVED***/rpm/repo.rpm" >/dev/null || ***REMOVED***
  mkdir -p "$HOME/rpmbuild/SOURCES"
  cp -av "$***REMOVED***pkgdir***REMOVED***/key.gpg" "$HOME/rpmbuild/SOURCES/RPM-GPG-KEY-k6-io"
  rpmbuild -ba "$HOME/rpmbuild/SPECS/k6-rpm-repo.spec"
  s3cmd put "$(find "$HOME/rpmbuild/RPMS/" -type f -name '*.rpm')" "s3://$***REMOVED***s3bucket***REMOVED***/rpm/repo.rpm"
***REMOVED***

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
