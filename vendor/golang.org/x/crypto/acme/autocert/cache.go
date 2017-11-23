// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ErrCacheMiss is returned when a certificate is not found in cache.
var ErrCacheMiss = errors.New("acme/autocert: certificate cache miss")

// Cache is used by Manager to store and retrieve previously obtained certificates
// as opaque data.
//
// The key argument of the methods refers to a domain name but need not be an FQDN.
// Cache implementations should not rely on the key naming pattern.
type Cache interface ***REMOVED***
	// Get returns a certificate data for the specified key.
	// If there's no such key, Get returns ErrCacheMiss.
	Get(ctx context.Context, key string) ([]byte, error)

	// Put stores the data in the cache under the specified key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, Get, results in the original data.
	Put(ctx context.Context, key string, data []byte) error

	// Delete removes a certificate data from the cache under the specified key.
	// If there's no such key in the cache, Delete returns nil.
	Delete(ctx context.Context, key string) error
***REMOVED***

// DirCache implements Cache using a directory on the local filesystem.
// If the directory does not exist, it will be created with 0700 permissions.
type DirCache string

// Get reads a certificate data from the specified file name.
func (d DirCache) Get(ctx context.Context, name string) ([]byte, error) ***REMOVED***
	name = filepath.Join(string(d), name)
	var (
		data []byte
		err  error
		done = make(chan struct***REMOVED******REMOVED***)
	)
	go func() ***REMOVED***
		data, err = ioutil.ReadFile(name)
		close(done)
	***REMOVED***()
	select ***REMOVED***
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	***REMOVED***
	if os.IsNotExist(err) ***REMOVED***
		return nil, ErrCacheMiss
	***REMOVED***
	return data, err
***REMOVED***

// Put writes the certificate data to the specified file name.
// The file will be created with 0600 permissions.
func (d DirCache) Put(ctx context.Context, name string, data []byte) error ***REMOVED***
	if err := os.MkdirAll(string(d), 0700); err != nil ***REMOVED***
		return err
	***REMOVED***

	done := make(chan struct***REMOVED******REMOVED***)
	var err error
	go func() ***REMOVED***
		defer close(done)
		var tmp string
		if tmp, err = d.writeTempFile(name, data); err != nil ***REMOVED***
			return
		***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			// Don't overwrite the file if the context was canceled.
		default:
			newName := filepath.Join(string(d), name)
			err = os.Rename(tmp, newName)
		***REMOVED***
	***REMOVED***()
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	***REMOVED***
	return err
***REMOVED***

// Delete removes the specified file name.
func (d DirCache) Delete(ctx context.Context, name string) error ***REMOVED***
	name = filepath.Join(string(d), name)
	var (
		err  error
		done = make(chan struct***REMOVED******REMOVED***)
	)
	go func() ***REMOVED***
		err = os.Remove(name)
		close(done)
	***REMOVED***()
	select ***REMOVED***
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	***REMOVED***
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// writeTempFile writes b to a temporary file, closes the file and returns its path.
func (d DirCache) writeTempFile(prefix string, b []byte) (string, error) ***REMOVED***
	// TempFile uses 0600 permissions
	f, err := ioutil.TempFile(string(d), prefix)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if _, err := f.Write(b); err != nil ***REMOVED***
		f.Close()
		return "", err
	***REMOVED***
	return f.Name(), f.Close()
***REMOVED***
