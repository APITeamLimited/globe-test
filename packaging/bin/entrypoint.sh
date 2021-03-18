#!/bin/bash
set -eEuo pipefail

log() ***REMOVED***
    echo "$(date -Iseconds) $*"
***REMOVED***

signkeypath="$PWD/sign-key.gpg"

if ! [ -r "$signkeypath" ]; then
  log "ERROR: Signing key not found at '$signkeypath'"
  exit 1
fi

gpg2 --import --batch --passphrase="$PGP_SIGN_KEY_PASSPHRASE" "$signkeypath"
export PGPKEYID="$(gpg2 --list-secret-keys --with-colons | grep '^sec' | cut -d: -f5)"

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
  "create-$***REMOVED***repo***REMOVED***-repo.sh" "$PWD/dist" "$PWD/dl.k6.io/$***REMOVED***repo***REMOVED***"
done

log "Generating index.html ..."
(cd dl.k6.io && generate_index.py -r)

log "Syncing to S3 ..."
s3cmd sync ./dl.k6.io/ "s3://$***REMOVED***S3_BUCKET-dl-k6-io***REMOVED***/"

exec "$@"
