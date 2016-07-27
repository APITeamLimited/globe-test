/usr/local/src/cfssl-***REMOVED******REMOVED*** pillar.cfssl.version ***REMOVED******REMOVED***:
  file.managed:
    - source: https://pkg.cfssl.org/R***REMOVED******REMOVED*** pillar.cfssl.version ***REMOVED******REMOVED***/cfssl_linux-amd64
    - source_hash: ***REMOVED******REMOVED*** pillar.cfssl.hash.cfssl ***REMOVED******REMOVED***
    - mode: 755

/usr/local/bin/cfssl:
  file.symlink:
    - target: /usr/local/src/cfssl-***REMOVED******REMOVED*** pillar.cfssl.version ***REMOVED******REMOVED***

/usr/local/src/cfssljson-***REMOVED******REMOVED*** pillar.cfssl.version ***REMOVED******REMOVED***:
  file.managed:
    - source: https://pkg.cfssl.org/R***REMOVED******REMOVED*** pillar.cfssl.version ***REMOVED******REMOVED***/cfssljson_linux-amd64
    - source_hash: ***REMOVED******REMOVED*** pillar.cfssl.hash.cfssl ***REMOVED******REMOVED***
    - mode: 755

/usr/local/bin/cfssljson:
  file.symlink:
    - target: /usr/local/src/cfssljson-***REMOVED******REMOVED*** pillar.cfssl.version ***REMOVED******REMOVED***

/etc/ssl/ca:
  file.recurse:
    - source: salt://cfssl/ca

/usr/local/share/ca-certificates/ca.crt:
  file.symlink:
    - target: /etc/ssl/ca/ca.pem
  cmd.watch:
    - name: update-ca-certificates
    - watch:
      - file: /usr/local/share/ca-certificates/ca.crt
