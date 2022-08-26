package proto

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"

	"github.com/APITeamLimited/redis/v9/internal/util"
)

// redis resp protocol data type.
const (
	RespStatus    = '+' // +<string>\r\n
	RespError     = '-' // -<string>\r\n
	RespString    = '$' // $<length>\r\n<bytes>\r\n
	RespInt       = ':' // :<number>\r\n
	RespNil       = '_' // _\r\n
	RespFloat     = ',' // ,<floating-point-number>\r\n (golang float)
	RespBool      = '#' // true: #t\r\n false: #f\r\n
	RespBlobError = '!' // !<length>\r\n<bytes>\r\n
	RespVerbatim  = '=' // =<length>\r\nFORMAT:<bytes>\r\n
	RespBigInt    = '(' // (<big number>\r\n
	RespArray     = '*' // *<len>\r\n... (same as resp2)
	RespMap       = '%' // %<len>\r\n(key)\r\n(value)\r\n... (golang map)
	RespSet       = '~' // ~<len>\r\n... (same as Array)
	RespAttr      = '|' // |<len>\r\n(key)\r\n(value)\r\n... + command reply
	RespPush      = '>' // ><len>\r\n... (same as Array)
)

// Not used temporarily.
// Redis has not used these two data types for the time being, and will implement them later.
// Streamed           = "EOF:"
// StreamedAggregated = '?'

//------------------------------------------------------------------------------

const Nil = RedisError("redis: nil") // nolint:errname

type RedisError string

func (e RedisError) Error() string ***REMOVED*** return string(e) ***REMOVED***

func (RedisError) RedisError() ***REMOVED******REMOVED***

func ParseErrorReply(line []byte) error ***REMOVED***
	return RedisError(line[1:])
***REMOVED***

//------------------------------------------------------------------------------

type Reader struct ***REMOVED***
	rd *bufio.Reader
***REMOVED***

func NewReader(rd io.Reader) *Reader ***REMOVED***
	return &Reader***REMOVED***
		rd: bufio.NewReader(rd),
	***REMOVED***
***REMOVED***

func (r *Reader) Buffered() int ***REMOVED***
	return r.rd.Buffered()
***REMOVED***

func (r *Reader) Peek(n int) ([]byte, error) ***REMOVED***
	return r.rd.Peek(n)
***REMOVED***

func (r *Reader) Reset(rd io.Reader) ***REMOVED***
	r.rd.Reset(rd)
***REMOVED***

// PeekReplyType returns the data type of the next response without advancing the Reader,
// and discard the attribute type.
func (r *Reader) PeekReplyType() (byte, error) ***REMOVED***
	b, err := r.rd.Peek(1)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if b[0] == RespAttr ***REMOVED***
		if err = r.DiscardNext(); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return r.PeekReplyType()
	***REMOVED***
	return b[0], nil
***REMOVED***

// ReadLine Return a valid reply, it will check the protocol or redis error,
// and discard the attribute type.
func (r *Reader) ReadLine() ([]byte, error) ***REMOVED***
	line, err := r.readLine()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch line[0] ***REMOVED***
	case RespError:
		return nil, ParseErrorReply(line)
	case RespNil:
		return nil, Nil
	case RespBlobError:
		var blobErr string
		blobErr, err = r.readStringReply(line)
		if err == nil ***REMOVED***
			err = RedisError(blobErr)
		***REMOVED***
		return nil, err
	case RespAttr:
		if err = r.Discard(line); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return r.ReadLine()
	***REMOVED***

	// Compatible with RESP2
	if IsNilReply(line) ***REMOVED***
		return nil, Nil
	***REMOVED***

	return line, nil
***REMOVED***

// readLine returns an error if:
//   - there is a pending read error;
//   - or line does not end with \r\n.
func (r *Reader) readLine() ([]byte, error) ***REMOVED***
	b, err := r.rd.ReadSlice('\n')
	if err != nil ***REMOVED***
		if err != bufio.ErrBufferFull ***REMOVED***
			return nil, err
		***REMOVED***

		full := make([]byte, len(b))
		copy(full, b)

		b, err = r.rd.ReadBytes('\n')
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		full = append(full, b...) //nolint:makezero
		b = full
	***REMOVED***
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' ***REMOVED***
		return nil, fmt.Errorf("redis: invalid reply: %q", b)
	***REMOVED***
	return b[:len(b)-2], nil
***REMOVED***

func (r *Reader) ReadReply() (interface***REMOVED******REMOVED***, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch line[0] ***REMOVED***
	case RespStatus:
		return string(line[1:]), nil
	case RespInt:
		return util.ParseInt(line[1:], 10, 64)
	case RespFloat:
		return r.readFloat(line)
	case RespBool:
		return r.readBool(line)
	case RespBigInt:
		return r.readBigInt(line)

	case RespString:
		return r.readStringReply(line)
	case RespVerbatim:
		return r.readVerb(line)

	case RespArray, RespSet, RespPush:
		return r.readSlice(line)
	case RespMap:
		return r.readMap(line)
	***REMOVED***
	return nil, fmt.Errorf("redis: can't parse %.100q", line)
***REMOVED***

func (r *Reader) readFloat(line []byte) (float64, error) ***REMOVED***
	v := string(line[1:])
	switch string(line[1:]) ***REMOVED***
	case "inf":
		return math.Inf(1), nil
	case "-inf":
		return math.Inf(-1), nil
	***REMOVED***
	return strconv.ParseFloat(v, 64)
***REMOVED***

func (r *Reader) readBool(line []byte) (bool, error) ***REMOVED***
	switch string(line[1:]) ***REMOVED***
	case "t":
		return true, nil
	case "f":
		return false, nil
	***REMOVED***
	return false, fmt.Errorf("redis: can't parse bool reply: %q", line)
***REMOVED***

func (r *Reader) readBigInt(line []byte) (*big.Int, error) ***REMOVED***
	i := new(big.Int)
	if i, ok := i.SetString(string(line[1:]), 10); ok ***REMOVED***
		return i, nil
	***REMOVED***
	return nil, fmt.Errorf("redis: can't parse bigInt reply: %q", line)
***REMOVED***

func (r *Reader) readStringReply(line []byte) (string, error) ***REMOVED***
	n, err := replyLen(line)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	b := make([]byte, n+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return util.BytesToString(b[:n]), nil
***REMOVED***

func (r *Reader) readVerb(line []byte) (string, error) ***REMOVED***
	s, err := r.readStringReply(line)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(s) < 4 || s[3] != ':' ***REMOVED***
		return "", fmt.Errorf("redis: can't parse verbatim string reply: %q", line)
	***REMOVED***
	return s[4:], nil
***REMOVED***

func (r *Reader) readSlice(line []byte) ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	n, err := replyLen(line)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	val := make([]interface***REMOVED******REMOVED***, n)
	for i := 0; i < len(val); i++ ***REMOVED***
		v, err := r.ReadReply()
		if err != nil ***REMOVED***
			if err == Nil ***REMOVED***
				val[i] = nil
				continue
			***REMOVED***
			if err, ok := err.(RedisError); ok ***REMOVED***
				val[i] = err
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		val[i] = v
	***REMOVED***
	return val, nil
***REMOVED***

func (r *Reader) readMap(line []byte) (map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, error) ***REMOVED***
	n, err := replyLen(line)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m := make(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***, n)
	for i := 0; i < n; i++ ***REMOVED***
		k, err := r.ReadReply()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		v, err := r.ReadReply()
		if err != nil ***REMOVED***
			if err == Nil ***REMOVED***
				m[k] = nil
				continue
			***REMOVED***
			if err, ok := err.(RedisError); ok ***REMOVED***
				m[k] = err
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		m[k] = v
	***REMOVED***
	return m, nil
***REMOVED***

// -------------------------------

func (r *Reader) ReadInt() (int64, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch line[0] ***REMOVED***
	case RespInt, RespStatus:
		return util.ParseInt(line[1:], 10, 64)
	case RespString:
		s, err := r.readStringReply(line)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return util.ParseInt([]byte(s), 10, 64)
	case RespBigInt:
		b, err := r.readBigInt(line)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if !b.IsInt64() ***REMOVED***
			return 0, fmt.Errorf("bigInt(%s) value out of range", b.String())
		***REMOVED***
		return b.Int64(), nil
	***REMOVED***
	return 0, fmt.Errorf("redis: can't parse int reply: %.100q", line)
***REMOVED***

func (r *Reader) ReadFloat() (float64, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch line[0] ***REMOVED***
	case RespFloat:
		return r.readFloat(line)
	case RespStatus:
		return strconv.ParseFloat(string(line[1:]), 64)
	case RespString:
		s, err := r.readStringReply(line)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		return strconv.ParseFloat(s, 64)
	***REMOVED***
	return 0, fmt.Errorf("redis: can't parse float reply: %.100q", line)
***REMOVED***

func (r *Reader) ReadString() (string, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	switch line[0] ***REMOVED***
	case RespStatus, RespInt, RespFloat:
		return string(line[1:]), nil
	case RespString:
		return r.readStringReply(line)
	case RespBool:
		b, err := r.readBool(line)
		return strconv.FormatBool(b), err
	case RespVerbatim:
		return r.readVerb(line)
	case RespBigInt:
		b, err := r.readBigInt(line)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		return b.String(), nil
	***REMOVED***
	return "", fmt.Errorf("redis: can't parse reply=%.100q reading string", line)
***REMOVED***

func (r *Reader) ReadBool() (bool, error) ***REMOVED***
	s, err := r.ReadString()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return s == "OK" || s == "1" || s == "true", nil
***REMOVED***

func (r *Reader) ReadSlice() ([]interface***REMOVED******REMOVED***, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.readSlice(line)
***REMOVED***

// ReadFixedArrayLen read fixed array length.
func (r *Reader) ReadFixedArrayLen(fixedLen int) error ***REMOVED***
	n, err := r.ReadArrayLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n != fixedLen ***REMOVED***
		return fmt.Errorf("redis: got %d elements in the array, wanted %d", n, fixedLen)
	***REMOVED***
	return nil
***REMOVED***

// ReadArrayLen Read and return the length of the array.
func (r *Reader) ReadArrayLen() (int, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch line[0] ***REMOVED***
	case RespArray, RespSet, RespPush:
		return replyLen(line)
	default:
		return 0, fmt.Errorf("redis: can't parse array/set/push reply: %.100q", line)
	***REMOVED***
***REMOVED***

// ReadFixedMapLen reads fixed map length.
func (r *Reader) ReadFixedMapLen(fixedLen int) error ***REMOVED***
	n, err := r.ReadMapLen()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if n != fixedLen ***REMOVED***
		return fmt.Errorf("redis: got %d elements in the map, wanted %d", n, fixedLen)
	***REMOVED***
	return nil
***REMOVED***

// ReadMapLen reads the length of the map type.
// If responding to the array type (RespArray/RespSet/RespPush),
// it must be a multiple of 2 and return n/2.
// Other types will return an error.
func (r *Reader) ReadMapLen() (int, error) ***REMOVED***
	line, err := r.ReadLine()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch line[0] ***REMOVED***
	case RespMap:
		return replyLen(line)
	case RespArray, RespSet, RespPush:
		// Some commands and RESP2 protocol may respond to array types.
		n, err := replyLen(line)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		if n%2 != 0 ***REMOVED***
			return 0, fmt.Errorf("redis: the length of the array must be a multiple of 2, got: %d", n)
		***REMOVED***
		return n / 2, nil
	default:
		return 0, fmt.Errorf("redis: can't parse map reply: %.100q", line)
	***REMOVED***
***REMOVED***

// DiscardNext read and discard the data represented by the next line.
func (r *Reader) DiscardNext() error ***REMOVED***
	line, err := r.readLine()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return r.Discard(line)
***REMOVED***

// Discard the data represented by line.
func (r *Reader) Discard(line []byte) (err error) ***REMOVED***
	if len(line) == 0 ***REMOVED***
		return errors.New("redis: invalid line")
	***REMOVED***
	switch line[0] ***REMOVED***
	case RespStatus, RespError, RespInt, RespNil, RespFloat, RespBool, RespBigInt:
		return nil
	***REMOVED***

	n, err := replyLen(line)
	if err != nil && err != Nil ***REMOVED***
		return err
	***REMOVED***

	switch line[0] ***REMOVED***
	case RespBlobError, RespString, RespVerbatim:
		// +\r\n
		_, err = r.rd.Discard(n + 2)
		return err
	case RespArray, RespSet, RespPush:
		for i := 0; i < n; i++ ***REMOVED***
			if err = r.DiscardNext(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	case RespMap, RespAttr:
		// Read key & value.
		for i := 0; i < n*2; i++ ***REMOVED***
			if err = r.DiscardNext(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	return fmt.Errorf("redis: can't parse %.100q", line)
***REMOVED***

func replyLen(line []byte) (n int, err error) ***REMOVED***
	n, err = util.Atoi(line[1:])
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if n < -1 ***REMOVED***
		return 0, fmt.Errorf("redis: invalid reply: %q", line)
	***REMOVED***

	switch line[0] ***REMOVED***
	case RespString, RespVerbatim, RespBlobError,
		RespArray, RespSet, RespPush, RespMap, RespAttr:
		if n == -1 ***REMOVED***
			return 0, Nil
		***REMOVED***
	***REMOVED***
	return n, nil
***REMOVED***

// IsNilReply detects redis.Nil of RESP2.
func IsNilReply(line []byte) bool ***REMOVED***
	return len(line) == 3 &&
		(line[0] == RespString || line[0] == RespArray) &&
		line[1] == '-' && line[2] == '1'
***REMOVED***
