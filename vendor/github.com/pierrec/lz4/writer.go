package lz4

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"runtime"
)

// Writer implements the LZ4 frame encoder.
type Writer struct ***REMOVED***
	Header
	dst      io.Writer
	checksum hash.Hash32 // frame checksum
	data     []byte      // data to be compressed, only used when dealing with block dependency as we need 64Kb to work with
	window   []byte      // last 64KB of decompressed data (block dependency) + blockMaxSize buffer

	zbCompressBuf []byte // buffer for compressing lz4 blocks
	writeSizeBuf  []byte // four-byte slice for writing checksums and sizes in writeblock
***REMOVED***

// NewWriter returns a new LZ4 frame encoder.
// No access to the underlying io.Writer is performed.
// The supplied Header is checked at the first Write.
// It is ok to change it before the first Write but then not until a Reset() is performed.
func NewWriter(dst io.Writer) *Writer ***REMOVED***
	return &Writer***REMOVED***
		dst:      dst,
		checksum: hashPool.Get(),
		Header: Header***REMOVED***
			BlockMaxSize: 4 << 20,
		***REMOVED***,
		writeSizeBuf: make([]byte, 4),
	***REMOVED***
***REMOVED***

// writeHeader builds and writes the header (magic+header) to the underlying io.Writer.
func (z *Writer) writeHeader() error ***REMOVED***
	// Default to 4Mb if BlockMaxSize is not set
	if z.Header.BlockMaxSize == 0 ***REMOVED***
		z.Header.BlockMaxSize = 4 << 20
	***REMOVED***
	// the only option that need to be validated
	bSize, ok := bsMapValue[z.Header.BlockMaxSize]
	if !ok ***REMOVED***
		return fmt.Errorf("lz4: invalid block max size: %d", z.Header.BlockMaxSize)
	***REMOVED***

	// magic number(4) + header(flags(2)+[Size(8)+DictID(4)]+checksum(1)) does not exceed 19 bytes
	// Size and DictID are optional
	var buf [19]byte

	// set the fixed size data: magic number, block max size and flags
	binary.LittleEndian.PutUint32(buf[0:], frameMagic)
	flg := byte(Version << 6)
	if !z.Header.BlockDependency ***REMOVED***
		flg |= 1 << 5
	***REMOVED***
	if z.Header.BlockChecksum ***REMOVED***
		flg |= 1 << 4
	***REMOVED***
	if z.Header.Size > 0 ***REMOVED***
		flg |= 1 << 3
	***REMOVED***
	if !z.Header.NoChecksum ***REMOVED***
		flg |= 1 << 2
	***REMOVED***
	//  if z.Header.Dict ***REMOVED***
	//      flg |= 1
	//  ***REMOVED***
	buf[4] = flg
	buf[5] = bSize << 4

	// current buffer size: magic(4) + flags(1) + block max size (1)
	n := 6
	// optional items
	if z.Header.Size > 0 ***REMOVED***
		binary.LittleEndian.PutUint64(buf[n:], z.Header.Size)
		n += 8
	***REMOVED***
	//  if z.Header.Dict ***REMOVED***
	//      binary.LittleEndian.PutUint32(buf[n:], z.Header.DictID)
	//      n += 4
	//  ***REMOVED***

	// header checksum includes the flags, block max size and optional Size and DictID
	z.checksum.Write(buf[4:n])
	buf[n] = byte(z.checksum.Sum32() >> 8 & 0xFF)
	z.checksum.Reset()

	// header ready, write it out
	if _, err := z.dst.Write(buf[0 : n+1]); err != nil ***REMOVED***
		return err
	***REMOVED***
	z.Header.done = true

	// initialize buffers dependent on header info
	z.zbCompressBuf = make([]byte, winSize+z.BlockMaxSize)

	return nil
***REMOVED***

// Write compresses data from the supplied buffer into the underlying io.Writer.
// Write does not return until the data has been written.
//
// If the input buffer is large enough (typically in multiples of BlockMaxSize)
// the data will be compressed concurrently.
//
// Write never buffers any data unless in BlockDependency mode where it may
// do so until it has 64Kb of data, after which it never buffers any.
func (z *Writer) Write(buf []byte) (n int, err error) ***REMOVED***
	if !z.Header.done ***REMOVED***
		if err = z.writeHeader(); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if len(buf) == 0 ***REMOVED***
		return
	***REMOVED***

	if !z.NoChecksum ***REMOVED***
		z.checksum.Write(buf)
	***REMOVED***

	// with block dependency, require at least 64Kb of data to work with
	// not having 64Kb only matters initially to setup the first window
	bl := 0
	if z.BlockDependency && len(z.window) == 0 ***REMOVED***
		bl = len(z.data)
		z.data = append(z.data, buf...)
		if len(z.data) < winSize ***REMOVED***
			return len(buf), nil
		***REMOVED***
		buf = z.data
		z.data = nil
	***REMOVED***

	// Break up the input buffer into BlockMaxSize blocks, provisioning the left over block.
	// Then compress into each of them concurrently if possible (no dependency).
	var (
		zb       block
		wbuf     = buf
		zn       = len(wbuf) / z.BlockMaxSize
		zi       = 0
		leftover = len(buf) % z.BlockMaxSize
	)

loop:
	for zi < zn ***REMOVED***
		if z.BlockDependency ***REMOVED***
			if zi == 0 ***REMOVED***
				// first block does not have the window
				zb.data = append(z.window, wbuf[:z.BlockMaxSize]...)
				zb.offset = len(z.window)
				wbuf = wbuf[z.BlockMaxSize-winSize:]
			***REMOVED*** else ***REMOVED***
				// set the uncompressed data including the window from previous block
				zb.data = wbuf[:z.BlockMaxSize+winSize]
				zb.offset = winSize
				wbuf = wbuf[z.BlockMaxSize:]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			zb.data = wbuf[:z.BlockMaxSize]
			wbuf = wbuf[z.BlockMaxSize:]
		***REMOVED***

		goto write
	***REMOVED***

	// left over
	if leftover > 0 ***REMOVED***
		zb = block***REMOVED***data: wbuf***REMOVED***
		if z.BlockDependency ***REMOVED***
			if zn == 0 ***REMOVED***
				zb.data = append(z.window, zb.data...)
				zb.offset = len(z.window)
			***REMOVED*** else ***REMOVED***
				zb.offset = winSize
			***REMOVED***
		***REMOVED***

		leftover = 0
		goto write
	***REMOVED***

	if z.BlockDependency ***REMOVED***
		if len(z.window) == 0 ***REMOVED***
			z.window = make([]byte, winSize)
		***REMOVED***
		// last buffer may be shorter than the window
		if len(buf) >= winSize ***REMOVED***
			copy(z.window, buf[len(buf)-winSize:])
		***REMOVED*** else ***REMOVED***
			copy(z.window, z.window[len(buf):])
			copy(z.window[len(buf)+1:], buf)
		***REMOVED***
	***REMOVED***

	return

write:
	zb = z.compressBlock(zb)
	_, err = z.writeBlock(zb)

	written := len(zb.data)
	if bl > 0 ***REMOVED***
		if written >= bl ***REMOVED***
			written -= bl
			bl = 0
		***REMOVED*** else ***REMOVED***
			bl -= written
			written = 0
		***REMOVED***
	***REMOVED***

	n += written
	// remove the window in zb.data
	if z.BlockDependency ***REMOVED***
		if zi == 0 ***REMOVED***
			n -= len(z.window)
		***REMOVED*** else ***REMOVED***
			n -= winSize
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return
	***REMOVED***
	zi++
	goto loop
***REMOVED***

// compressBlock compresses a block.
func (z *Writer) compressBlock(zb block) block ***REMOVED***
	// compressed block size cannot exceed the input's
	var (
		n    int
		err  error
		zbuf = z.zbCompressBuf
	)
	if z.HighCompression ***REMOVED***
		n, err = CompressBlockHC(zb.data, zbuf, zb.offset)
	***REMOVED*** else ***REMOVED***
		n, err = CompressBlock(zb.data, zbuf, zb.offset)
	***REMOVED***

	// compressible and compressed size smaller than decompressed: ok!
	if err == nil && n > 0 && len(zb.zdata) < len(zb.data) ***REMOVED***
		zb.compressed = true
		zb.zdata = zbuf[:n]
	***REMOVED*** else ***REMOVED***
		zb.compressed = false
		zb.zdata = zb.data[zb.offset:]
	***REMOVED***

	if z.BlockChecksum ***REMOVED***
		xxh := hashPool.Get()
		xxh.Write(zb.zdata)
		zb.checksum = xxh.Sum32()
		hashPool.Put(xxh)
	***REMOVED***

	return zb
***REMOVED***

// writeBlock writes a frame block to the underlying io.Writer (size, data).
func (z *Writer) writeBlock(zb block) (int, error) ***REMOVED***
	bLen := uint32(len(zb.zdata))
	if !zb.compressed ***REMOVED***
		bLen |= 1 << 31
	***REMOVED***

	n := 0

	binary.LittleEndian.PutUint32(z.writeSizeBuf, bLen)
	n, err := z.dst.Write(z.writeSizeBuf)
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***

	m, err := z.dst.Write(zb.zdata)
	n += m
	if err != nil ***REMOVED***
		return n, err
	***REMOVED***

	if z.BlockChecksum ***REMOVED***
		binary.LittleEndian.PutUint32(z.writeSizeBuf, zb.checksum)
		m, err := z.dst.Write(z.writeSizeBuf)
		n += m

		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
	***REMOVED***

	return n, nil
***REMOVED***

// Flush flushes any pending compressed data to the underlying writer.
// Flush does not return until the data has been written.
// If the underlying writer returns an error, Flush returns that error.
//
// Flush is only required when in BlockDependency mode and the total of
// data written is less than 64Kb.
func (z *Writer) Flush() error ***REMOVED***
	if len(z.data) == 0 ***REMOVED***
		return nil
	***REMOVED***

	zb := z.compressBlock(block***REMOVED***data: z.data***REMOVED***)
	if _, err := z.writeBlock(zb); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Close closes the Writer, flushing any unwritten data to the underlying io.Writer, but does not close the underlying io.Writer.
func (z *Writer) Close() error ***REMOVED***
	if !z.Header.done ***REMOVED***
		if err := z.writeHeader(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// buffered data for the block dependency window
	if z.BlockDependency && len(z.data) > 0 ***REMOVED***
		zb := block***REMOVED***data: z.data***REMOVED***
		if _, err := z.writeBlock(z.compressBlock(zb)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := binary.Write(z.dst, binary.LittleEndian, uint32(0)); err != nil ***REMOVED***
		return err
	***REMOVED***
	if !z.NoChecksum ***REMOVED***
		if err := binary.Write(z.dst, binary.LittleEndian, z.checksum.Sum32()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Reset clears the state of the Writer z such that it is equivalent to its
// initial state from NewWriter, but instead writing to w.
// No access to the underlying io.Writer is performed.
func (z *Writer) Reset(w io.Writer) ***REMOVED***
	z.Header = Header***REMOVED******REMOVED***
	z.dst = w
	z.checksum.Reset()
	z.data = nil
	z.window = nil
***REMOVED***

// ReadFrom compresses the data read from the io.Reader and writes it to the underlying io.Writer.
// Returns the number of bytes read.
// It does not close the Writer.
func (z *Writer) ReadFrom(r io.Reader) (n int64, err error) ***REMOVED***
	cpus := runtime.GOMAXPROCS(0)
	buf := make([]byte, cpus*z.BlockMaxSize)
	for ***REMOVED***
		m, er := io.ReadFull(r, buf)
		n += int64(m)
		if er == nil || er == io.ErrUnexpectedEOF || er == io.EOF ***REMOVED***
			if _, err = z.Write(buf[:m]); err != nil ***REMOVED***
				return
			***REMOVED***
			if er == nil ***REMOVED***
				continue
			***REMOVED***
			return
		***REMOVED***
		return n, er
	***REMOVED***
***REMOVED***
