// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// A Decoder loads an archive of CLDR data.
type Decoder struct ***REMOVED***
	dirFilter     []string
	sectionFilter []string
	loader        Loader
	cldr          *CLDR
	curLocale     string
***REMOVED***

// SetSectionFilter takes a list top-level LDML element names to which
// evaluation of LDML should be limited.  It automatically calls SetDirFilter.
func (d *Decoder) SetSectionFilter(filter ...string) ***REMOVED***
	d.sectionFilter = filter
	// TODO: automatically set dir filter
***REMOVED***

// SetDirFilter limits the loading of LDML XML files of the specied directories.
// Note that sections may be split across directories differently for different CLDR versions.
// For more robust code, use SetSectionFilter.
func (d *Decoder) SetDirFilter(dir ...string) ***REMOVED***
	d.dirFilter = dir
***REMOVED***

// A Loader provides access to the files of a CLDR archive.
type Loader interface ***REMOVED***
	Len() int
	Path(i int) string
	Reader(i int) (io.ReadCloser, error)
***REMOVED***

var fileRe = regexp.MustCompile(`.*[/\\](.*)[/\\](.*)\.xml`)

// Decode loads and decodes the files represented by l.
func (d *Decoder) Decode(l Loader) (cldr *CLDR, err error) ***REMOVED***
	d.cldr = makeCLDR()
	for i := 0; i < l.Len(); i++ ***REMOVED***
		fname := l.Path(i)
		if m := fileRe.FindStringSubmatch(fname); m != nil ***REMOVED***
			if len(d.dirFilter) > 0 && !in(d.dirFilter, m[1]) ***REMOVED***
				continue
			***REMOVED***
			var r io.Reader
			if r, err = l.Reader(i); err == nil ***REMOVED***
				err = d.decode(m[1], m[2], r)
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	d.cldr.finalize(d.sectionFilter)
	return d.cldr, nil
***REMOVED***

func (d *Decoder) decode(dir, id string, r io.Reader) error ***REMOVED***
	var v interface***REMOVED******REMOVED***
	var l *LDML
	cldr := d.cldr
	switch ***REMOVED***
	case dir == "supplemental":
		v = cldr.supp
	case dir == "transforms":
		return nil
	case dir == "bcp47":
		v = cldr.bcp47
	case dir == "validity":
		return nil
	default:
		ok := false
		if v, ok = cldr.locale[id]; !ok ***REMOVED***
			l = &LDML***REMOVED******REMOVED***
			v, cldr.locale[id] = l, l
		***REMOVED***
	***REMOVED***
	x := xml.NewDecoder(r)
	if err := x.Decode(v); err != nil ***REMOVED***
		log.Printf("%s/%s: %v", dir, id, err)
		return err
	***REMOVED***
	if l != nil ***REMOVED***
		if l.Identity == nil ***REMOVED***
			return fmt.Errorf("%s/%s: missing identity element", dir, id)
		***REMOVED***
		// TODO: verify when CLDR bug http://unicode.org/cldr/trac/ticket/8970
		// is resolved.
		// path := strings.Split(id, "_")
		// if lang := l.Identity.Language.Type; lang != path[0] ***REMOVED***
		// 	return fmt.Errorf("%s/%s: language was %s; want %s", dir, id, lang, path[0])
		// ***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type pathLoader []string

func makePathLoader(path string) (pl pathLoader, err error) ***REMOVED***
	err = filepath.Walk(path, func(path string, _ os.FileInfo, err error) error ***REMOVED***
		pl = append(pl, path)
		return err
	***REMOVED***)
	return pl, err
***REMOVED***

func (pl pathLoader) Len() int ***REMOVED***
	return len(pl)
***REMOVED***

func (pl pathLoader) Path(i int) string ***REMOVED***
	return pl[i]
***REMOVED***

func (pl pathLoader) Reader(i int) (io.ReadCloser, error) ***REMOVED***
	return os.Open(pl[i])
***REMOVED***

// DecodePath loads CLDR data from the given path.
func (d *Decoder) DecodePath(path string) (cldr *CLDR, err error) ***REMOVED***
	loader, err := makePathLoader(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return d.Decode(loader)
***REMOVED***

type zipLoader struct ***REMOVED***
	r *zip.Reader
***REMOVED***

func (zl zipLoader) Len() int ***REMOVED***
	return len(zl.r.File)
***REMOVED***

func (zl zipLoader) Path(i int) string ***REMOVED***
	return zl.r.File[i].Name
***REMOVED***

func (zl zipLoader) Reader(i int) (io.ReadCloser, error) ***REMOVED***
	return zl.r.File[i].Open()
***REMOVED***

// DecodeZip loads CLDR data from the zip archive for which r is the source.
func (d *Decoder) DecodeZip(r io.Reader) (cldr *CLDR, err error) ***REMOVED***
	buffer, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	archive, err := zip.NewReader(bytes.NewReader(buffer), int64(len(buffer)))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return d.Decode(zipLoader***REMOVED***archive***REMOVED***)
***REMOVED***
