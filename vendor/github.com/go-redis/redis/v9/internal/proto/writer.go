package proto

import (
	"encoding"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9/internal/util"
)

type writer interface ***REMOVED***
	io.Writer
	io.ByteWriter
	// WriteString implement io.StringWriter.
	WriteString(s string) (n int, err error)
***REMOVED***

type Writer struct ***REMOVED***
	writer

	lenBuf []byte
	numBuf []byte
***REMOVED***

func NewWriter(wr writer) *Writer ***REMOVED***
	return &Writer***REMOVED***
		writer: wr,

		lenBuf: make([]byte, 64),
		numBuf: make([]byte, 64),
	***REMOVED***
***REMOVED***

func (w *Writer) WriteArgs(args []interface***REMOVED******REMOVED***) error ***REMOVED***
	if err := w.WriteByte(RespArray); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := w.writeLen(len(args)); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, arg := range args ***REMOVED***
		if err := w.WriteArg(arg); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (w *Writer) writeLen(n int) error ***REMOVED***
	w.lenBuf = strconv.AppendUint(w.lenBuf[:0], uint64(n), 10)
	w.lenBuf = append(w.lenBuf, '\r', '\n')
	_, err := w.Write(w.lenBuf)
	return err
***REMOVED***

func (w *Writer) WriteArg(v interface***REMOVED******REMOVED***) error ***REMOVED***
	switch v := v.(type) ***REMOVED***
	case nil:
		return w.string("")
	case string:
		return w.string(v)
	case []byte:
		return w.bytes(v)
	case int:
		return w.int(int64(v))
	case int8:
		return w.int(int64(v))
	case int16:
		return w.int(int64(v))
	case int32:
		return w.int(int64(v))
	case int64:
		return w.int(v)
	case uint:
		return w.uint(uint64(v))
	case uint8:
		return w.uint(uint64(v))
	case uint16:
		return w.uint(uint64(v))
	case uint32:
		return w.uint(uint64(v))
	case uint64:
		return w.uint(v)
	case float32:
		return w.float(float64(v))
	case float64:
		return w.float(v)
	case bool:
		if v ***REMOVED***
			return w.int(1)
		***REMOVED***
		return w.int(0)
	case time.Time:
		w.numBuf = v.AppendFormat(w.numBuf[:0], time.RFC3339Nano)
		return w.bytes(w.numBuf)
	case time.Duration:
		return w.int(v.Nanoseconds())
	case encoding.BinaryMarshaler:
		b, err := v.MarshalBinary()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return w.bytes(b)
	case net.IP:
		return w.bytes(v)
	default:
		return fmt.Errorf(
			"redis: can't marshal %T (implement encoding.BinaryMarshaler)", v)
	***REMOVED***
***REMOVED***

func (w *Writer) bytes(b []byte) error ***REMOVED***
	if err := w.WriteByte(RespString); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := w.writeLen(len(b)); err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err := w.Write(b); err != nil ***REMOVED***
		return err
	***REMOVED***

	return w.crlf()
***REMOVED***

func (w *Writer) string(s string) error ***REMOVED***
	return w.bytes(util.StringToBytes(s))
***REMOVED***

func (w *Writer) uint(n uint64) error ***REMOVED***
	w.numBuf = strconv.AppendUint(w.numBuf[:0], n, 10)
	return w.bytes(w.numBuf)
***REMOVED***

func (w *Writer) int(n int64) error ***REMOVED***
	w.numBuf = strconv.AppendInt(w.numBuf[:0], n, 10)
	return w.bytes(w.numBuf)
***REMOVED***

func (w *Writer) float(f float64) error ***REMOVED***
	w.numBuf = strconv.AppendFloat(w.numBuf[:0], f, 'f', -1, 64)
	return w.bytes(w.numBuf)
***REMOVED***

func (w *Writer) crlf() error ***REMOVED***
	if err := w.WriteByte('\r'); err != nil ***REMOVED***
		return err
	***REMOVED***
	return w.WriteByte('\n')
***REMOVED***
