// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// NewListener returns a net.Listener that listens on the standard TLS
// port (443) on all interfaces and returns *tls.Conn connections with
// LetsEncrypt certificates for the provided domain or domains.
//
// It enables one-line HTTPS servers:
//
//     log.Fatal(http.Serve(autocert.NewListener("example.com"), handler))
//
// NewListener is a convenience function for a common configuration.
// More complex or custom configurations can use the autocert.Manager
// type instead.
//
// Use of this function implies acceptance of the LetsEncrypt Terms of
// Service. If domains is not empty, the provided domains are passed
// to HostWhitelist. If domains is empty, the listener will do
// LetsEncrypt challenges for any requested domain, which is not
// recommended.
//
// Certificates are cached in a "golang-autocert" directory under an
// operating system-specific cache or temp directory. This may not
// be suitable for servers spanning multiple machines.
//
// The returned listener uses a *tls.Config that enables HTTP/2, and
// should only be used with servers that support HTTP/2.
//
// The returned Listener also enables TCP keep-alives on the accepted
// connections. The returned *tls.Conn are returned before their TLS
// handshake has completed.
func NewListener(domains ...string) net.Listener ***REMOVED***
	m := &Manager***REMOVED***
		Prompt: AcceptTOS,
	***REMOVED***
	if len(domains) > 0 ***REMOVED***
		m.HostPolicy = HostWhitelist(domains...)
	***REMOVED***
	dir := cacheDir()
	if err := os.MkdirAll(dir, 0700); err != nil ***REMOVED***
		log.Printf("warning: autocert.NewListener not using a cache: %v", err)
	***REMOVED*** else ***REMOVED***
		m.Cache = DirCache(dir)
	***REMOVED***
	return m.Listener()
***REMOVED***

// Listener listens on the standard TLS port (443) on all interfaces
// and returns a net.Listener returning *tls.Conn connections.
//
// The returned listener uses a *tls.Config that enables HTTP/2, and
// should only be used with servers that support HTTP/2.
//
// The returned Listener also enables TCP keep-alives on the accepted
// connections. The returned *tls.Conn are returned before their TLS
// handshake has completed.
//
// Unlike NewListener, it is the caller's responsibility to initialize
// the Manager m's Prompt, Cache, HostPolicy, and other desired options.
func (m *Manager) Listener() net.Listener ***REMOVED***
	ln := &listener***REMOVED***
		conf: m.TLSConfig(),
	***REMOVED***
	ln.tcpListener, ln.tcpListenErr = net.Listen("tcp", ":443")
	return ln
***REMOVED***

type listener struct ***REMOVED***
	conf *tls.Config

	tcpListener  net.Listener
	tcpListenErr error
***REMOVED***

func (ln *listener) Accept() (net.Conn, error) ***REMOVED***
	if ln.tcpListenErr != nil ***REMOVED***
		return nil, ln.tcpListenErr
	***REMOVED***
	conn, err := ln.tcpListener.Accept()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tcpConn := conn.(*net.TCPConn)

	// Because Listener is a convenience function, help out with
	// this too.  This is not possible for the caller to set once
	// we return a *tcp.Conn wrapping an inaccessible net.Conn.
	// If callers don't want this, they can do things the manual
	// way and tweak as needed. But this is what net/http does
	// itself, so copy that. If net/http changes, we can change
	// here too.
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(3 * time.Minute)

	return tls.Server(tcpConn, ln.conf), nil
***REMOVED***

func (ln *listener) Addr() net.Addr ***REMOVED***
	if ln.tcpListener != nil ***REMOVED***
		return ln.tcpListener.Addr()
	***REMOVED***
	// net.Listen failed. Return something non-nil in case callers
	// call Addr before Accept:
	return &net.TCPAddr***REMOVED***IP: net.IP***REMOVED***0, 0, 0, 0***REMOVED***, Port: 443***REMOVED***
***REMOVED***

func (ln *listener) Close() error ***REMOVED***
	if ln.tcpListenErr != nil ***REMOVED***
		return ln.tcpListenErr
	***REMOVED***
	return ln.tcpListener.Close()
***REMOVED***

func homeDir() string ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	***REMOVED***
	if h := os.Getenv("HOME"); h != "" ***REMOVED***
		return h
	***REMOVED***
	return "/"
***REMOVED***

func cacheDir() string ***REMOVED***
	const base = "golang-autocert"
	switch runtime.GOOS ***REMOVED***
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Caches", base)
	case "windows":
		for _, ev := range []string***REMOVED***"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"***REMOVED*** ***REMOVED***
			if v := os.Getenv(ev); v != "" ***REMOVED***
				return filepath.Join(v, base)
			***REMOVED***
		***REMOVED***
		// Worst case:
		return filepath.Join(homeDir(), base)
	***REMOVED***
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" ***REMOVED***
		return filepath.Join(xdg, base)
	***REMOVED***
	return filepath.Join(homeDir(), ".cache", base)
***REMOVED***
