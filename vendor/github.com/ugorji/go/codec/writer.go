// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "io"

// encWriter abstracts writing to a byte array or to an io.Writer.
type encWriter interface ***REMOVED***
	writeb([]byte)
	writestr(string)
	writeqstr(string) // write string wrapped in quotes ie "..."
	writen1(byte)
	writen2(byte, byte)
	// writen will write up to 7 bytes at a time.
	writen(b [rwNLen]byte, num uint8)
	end()
***REMOVED***

// ---------------------------------------------

// bufioEncWriter
type bufioEncWriter struct ***REMOVED***
	w io.Writer

	buf []byte

	n int

	b [16]byte // scratch buffer and padding (cache-aligned)
***REMOVED***

func (z *bufioEncWriter) reset(w io.Writer, bufsize int, blist *bytesFreelist) ***REMOVED***
	z.w = w
	z.n = 0
	if bufsize <= 0 ***REMOVED***
		bufsize = defEncByteBufSize
	***REMOVED***
	// bufsize must be >= 8, to accomodate writen methods (where n <= 8)
	if bufsize <= 8 ***REMOVED***
		bufsize = 8
	***REMOVED***
	if cap(z.buf) < bufsize ***REMOVED***
		if len(z.buf) > 0 && &z.buf[0] != &z.b[0] ***REMOVED***
			blist.put(z.buf)
		***REMOVED***
		if len(z.b) > bufsize ***REMOVED***
			z.buf = z.b[:]
		***REMOVED*** else ***REMOVED***
			z.buf = blist.get(bufsize)
		***REMOVED***
	***REMOVED***
	z.buf = z.buf[:cap(z.buf)]
***REMOVED***

//go:noinline - flush only called intermittently
func (z *bufioEncWriter) flushErr() (err error) ***REMOVED***
	n, err := z.w.Write(z.buf[:z.n])
	z.n -= n
	if z.n > 0 && err == nil ***REMOVED***
		err = io.ErrShortWrite
	***REMOVED***
	if n > 0 && z.n > 0 ***REMOVED***
		copy(z.buf, z.buf[n:z.n+n])
	***REMOVED***
	return err
***REMOVED***

func (z *bufioEncWriter) flush() ***REMOVED***
	if err := z.flushErr(); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *bufioEncWriter) writeb(s []byte) ***REMOVED***
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a ***REMOVED***
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	***REMOVED***
	z.n += copy(z.buf[z.n:], s)
***REMOVED***

func (z *bufioEncWriter) writestr(s string) ***REMOVED***
	// z.writeb(bytesView(s)) // inlined below
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a ***REMOVED***
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	***REMOVED***
	z.n += copy(z.buf[z.n:], s)
***REMOVED***

func (z *bufioEncWriter) writeqstr(s string) ***REMOVED***
	// z.writen1('"')
	// z.writestr(s)
	// z.writen1('"')

	if z.n+len(s)+2 > len(z.buf) ***REMOVED***
		z.flush()
	***REMOVED***
	z.buf[z.n] = '"'
	z.n++
LOOP:
	a := len(z.buf) - z.n
	if len(s)+1 > a ***REMOVED***
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	***REMOVED***
	z.n += copy(z.buf[z.n:], s)
	z.buf[z.n] = '"'
	z.n++
***REMOVED***

func (z *bufioEncWriter) writen1(b1 byte) ***REMOVED***
	if 1 > len(z.buf)-z.n ***REMOVED***
		z.flush()
	***REMOVED***
	z.buf[z.n] = b1
	z.n++
***REMOVED***

func (z *bufioEncWriter) writen2(b1, b2 byte) ***REMOVED***
	if 2 > len(z.buf)-z.n ***REMOVED***
		z.flush()
	***REMOVED***
	z.buf[z.n+1] = b2
	z.buf[z.n] = b1
	z.n += 2
***REMOVED***

func (z *bufioEncWriter) writen(b [rwNLen]byte, num uint8) ***REMOVED***
	if int(num) > len(z.buf)-z.n ***REMOVED***
		z.flush()
	***REMOVED***
	copy(z.buf[z.n:], b[:num])
	z.n += int(num)
***REMOVED***

func (z *bufioEncWriter) endErr() (err error) ***REMOVED***
	if z.n > 0 ***REMOVED***
		err = z.flushErr()
	***REMOVED***
	return
***REMOVED***

// ---------------------------------------------

// bytesEncAppender implements encWriter and can write to an byte slice.
type bytesEncAppender struct ***REMOVED***
	b   []byte
	out *[]byte
***REMOVED***

func (z *bytesEncAppender) writeb(s []byte) ***REMOVED***
	z.b = append(z.b, s...)
***REMOVED***
func (z *bytesEncAppender) writestr(s string) ***REMOVED***
	z.b = append(z.b, s...)
***REMOVED***
func (z *bytesEncAppender) writeqstr(s string) ***REMOVED***
	z.b = append(append(append(z.b, '"'), s...), '"')

	// z.b = append(z.b, '"')
	// z.b = append(z.b, s...)
	// z.b = append(z.b, '"')
***REMOVED***
func (z *bytesEncAppender) writen1(b1 byte) ***REMOVED***
	z.b = append(z.b, b1)
***REMOVED***
func (z *bytesEncAppender) writen2(b1, b2 byte) ***REMOVED***
	z.b = append(z.b, b1, b2) // cost: 81
***REMOVED***
func (z *bytesEncAppender) writen(s [rwNLen]byte, num uint8) ***REMOVED***
	// if num <= rwNLen ***REMOVED***
	if int(num) <= len(s) ***REMOVED***
		z.b = append(z.b, s[:num]...)
	***REMOVED***
***REMOVED***
func (z *bytesEncAppender) endErr() error ***REMOVED***
	*(z.out) = z.b
	return nil
***REMOVED***
func (z *bytesEncAppender) reset(in []byte, out *[]byte) ***REMOVED***
	z.b = in[:0]
	z.out = out
***REMOVED***

// --------------------------------------------------

type encWr struct ***REMOVED***
	bytes bool // encoding to []byte
	js    bool // is json encoder?
	be    bool // is binary encoder?

	c containerState

	calls uint16

	wb bytesEncAppender
	wf *bufioEncWriter
***REMOVED***

func (z *encWr) writeb(s []byte) ***REMOVED***
	if z.bytes ***REMOVED***
		z.wb.writeb(s)
	***REMOVED*** else ***REMOVED***
		z.wf.writeb(s)
	***REMOVED***
***REMOVED***
func (z *encWr) writeqstr(s string) ***REMOVED***
	if z.bytes ***REMOVED***
		z.wb.writeqstr(s)
	***REMOVED*** else ***REMOVED***
		z.wf.writeqstr(s)
	***REMOVED***
***REMOVED***
func (z *encWr) writestr(s string) ***REMOVED***
	if z.bytes ***REMOVED***
		z.wb.writestr(s)
	***REMOVED*** else ***REMOVED***
		z.wf.writestr(s)
	***REMOVED***
***REMOVED***
func (z *encWr) writen1(b1 byte) ***REMOVED***
	if z.bytes ***REMOVED***
		z.wb.writen1(b1)
	***REMOVED*** else ***REMOVED***
		z.wf.writen1(b1)
	***REMOVED***
***REMOVED***
func (z *encWr) writen2(b1, b2 byte) ***REMOVED***
	if z.bytes ***REMOVED***
		z.wb.writen2(b1, b2)
	***REMOVED*** else ***REMOVED***
		z.wf.writen2(b1, b2)
	***REMOVED***
***REMOVED***
func (z *encWr) writen(b [rwNLen]byte, num uint8) ***REMOVED***
	if z.bytes ***REMOVED***
		z.wb.writen(b, num)
	***REMOVED*** else ***REMOVED***
		z.wf.writen(b, num)
	***REMOVED***
***REMOVED***
func (z *encWr) endErr() error ***REMOVED***
	if z.bytes ***REMOVED***
		return z.wb.endErr()
	***REMOVED***
	return z.wf.endErr()
***REMOVED***

func (z *encWr) end() ***REMOVED***
	if err := z.endErr(); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

var _ encWriter = (*encWr)(nil)
