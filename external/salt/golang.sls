golang-***REMOVED******REMOVED*** pillar.golang.version ***REMOVED******REMOVED***:
  archive.extracted:
    - name: /usr/local/src/go-***REMOVED******REMOVED*** pillar.golang.version ***REMOVED******REMOVED***
    - source: https://storage.googleapis.com/golang/go***REMOVED******REMOVED*** pillar.golang.version ***REMOVED******REMOVED***.linux-amd64.tar.gz
    - source_hash: ***REMOVED******REMOVED*** pillar.golang.hash ***REMOVED******REMOVED***
    - archive_format: tar

/usr/local/bin/go:
  file.symlink:
    - target: /usr/local/src/go-***REMOVED******REMOVED*** pillar.golang.version ***REMOVED******REMOVED***/go/bin/go

/usr/local/bin/gofmt:
  file.symlink:
    - target: /usr/local/src/go-***REMOVED******REMOVED*** pillar.golang.version ***REMOVED******REMOVED***/go/bin/gofmt

/usr/local/bin/godoc:
  file.symlink:
    - target: /usr/local/src/go-***REMOVED******REMOVED*** pillar.golang.version ***REMOVED******REMOVED***/go/bin/godoc

/etc/profile.d/golang.sh:
  file.managed:
    - contents: |
        export GOPATH=$HOME/go
        test -d $GOPATH || mkdir $GOPATH
    - mode: 755
