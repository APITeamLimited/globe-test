package sarama

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
)

type none struct***REMOVED******REMOVED***

// make []int32 sortable so we can sort partition numbers
type int32Slice []int32

func (slice int32Slice) Len() int ***REMOVED***
	return len(slice)
***REMOVED***

func (slice int32Slice) Less(i, j int) bool ***REMOVED***
	return slice[i] < slice[j]
***REMOVED***

func (slice int32Slice) Swap(i, j int) ***REMOVED***
	slice[i], slice[j] = slice[j], slice[i]
***REMOVED***

func dupInt32Slice(input []int32) []int32 ***REMOVED***
	ret := make([]int32, 0, len(input))
	for _, val := range input ***REMOVED***
		ret = append(ret, val)
	***REMOVED***
	return ret
***REMOVED***

func withRecover(fn func()) ***REMOVED***
	defer func() ***REMOVED***
		handler := PanicHandler
		if handler != nil ***REMOVED***
			if err := recover(); err != nil ***REMOVED***
				handler(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	fn()
***REMOVED***

func safeAsyncClose(b *Broker) ***REMOVED***
	tmp := b // local var prevents clobbering in goroutine
	go withRecover(func() ***REMOVED***
		if connected, _ := tmp.Connected(); connected ***REMOVED***
			if err := tmp.Close(); err != nil ***REMOVED***
				Logger.Println("Error closing broker", tmp.ID(), ":", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Encoder is a simple interface for any type that can be encoded as an array of bytes
// in order to be sent as the key or value of a Kafka message. Length() is provided as an
// optimization, and must return the same as len() on the result of Encode().
type Encoder interface ***REMOVED***
	Encode() ([]byte, error)
	Length() int
***REMOVED***

// make strings and byte slices encodable for convenience so they can be used as keys
// and/or values in kafka messages

// StringEncoder implements the Encoder interface for Go strings so that they can be used
// as the Key or Value in a ProducerMessage.
type StringEncoder string

func (s StringEncoder) Encode() ([]byte, error) ***REMOVED***
	return []byte(s), nil
***REMOVED***

func (s StringEncoder) Length() int ***REMOVED***
	return len(s)
***REMOVED***

// ByteEncoder implements the Encoder interface for Go byte slices so that they can be used
// as the Key or Value in a ProducerMessage.
type ByteEncoder []byte

func (b ByteEncoder) Encode() ([]byte, error) ***REMOVED***
	return b, nil
***REMOVED***

func (b ByteEncoder) Length() int ***REMOVED***
	return len(b)
***REMOVED***

// bufConn wraps a net.Conn with a buffer for reads to reduce the number of
// reads that trigger syscalls.
type bufConn struct ***REMOVED***
	net.Conn
	buf *bufio.Reader
***REMOVED***

func newBufConn(conn net.Conn) *bufConn ***REMOVED***
	return &bufConn***REMOVED***
		Conn: conn,
		buf:  bufio.NewReader(conn),
	***REMOVED***
***REMOVED***

func (bc *bufConn) Read(b []byte) (n int, err error) ***REMOVED***
	return bc.buf.Read(b)
***REMOVED***

// KafkaVersion instances represent versions of the upstream Kafka broker.
type KafkaVersion struct ***REMOVED***
	// it's a struct rather than just typing the array directly to make it opaque and stop people
	// generating their own arbitrary versions
	version [4]uint
***REMOVED***

func newKafkaVersion(major, minor, veryMinor, patch uint) KafkaVersion ***REMOVED***
	return KafkaVersion***REMOVED***
		version: [4]uint***REMOVED***major, minor, veryMinor, patch***REMOVED***,
	***REMOVED***
***REMOVED***

// IsAtLeast return true if and only if the version it is called on is
// greater than or equal to the version passed in:
//    V1.IsAtLeast(V2) // false
//    V2.IsAtLeast(V1) // true
func (v KafkaVersion) IsAtLeast(other KafkaVersion) bool ***REMOVED***
	for i := range v.version ***REMOVED***
		if v.version[i] > other.version[i] ***REMOVED***
			return true
		***REMOVED*** else if v.version[i] < other.version[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Effective constants defining the supported kafka versions.
var (
	V0_8_2_0   = newKafkaVersion(0, 8, 2, 0)
	V0_8_2_1   = newKafkaVersion(0, 8, 2, 1)
	V0_8_2_2   = newKafkaVersion(0, 8, 2, 2)
	V0_9_0_0   = newKafkaVersion(0, 9, 0, 0)
	V0_9_0_1   = newKafkaVersion(0, 9, 0, 1)
	V0_10_0_0  = newKafkaVersion(0, 10, 0, 0)
	V0_10_0_1  = newKafkaVersion(0, 10, 0, 1)
	V0_10_1_0  = newKafkaVersion(0, 10, 1, 0)
	V0_10_2_0  = newKafkaVersion(0, 10, 2, 0)
	V0_11_0_0  = newKafkaVersion(0, 11, 0, 0)
	V1_0_0_0   = newKafkaVersion(1, 0, 0, 0)
	minVersion = V0_8_2_0
)

func ParseKafkaVersion(s string) (KafkaVersion, error) ***REMOVED***
	var major, minor, veryMinor, patch uint
	var err error
	if s[0] == '0' ***REMOVED***
		err = scanKafkaVersion(s, `^0\.\d+\.\d+\.\d+$`, "0.%d.%d.%d", [3]*uint***REMOVED***&minor, &veryMinor, &patch***REMOVED***)
	***REMOVED*** else ***REMOVED***
		err = scanKafkaVersion(s, `^\d+\.\d+\.\d+$`, "%d.%d.%d", [3]*uint***REMOVED***&major, &minor, &veryMinor***REMOVED***)
	***REMOVED***
	if err != nil ***REMOVED***
		return minVersion, err
	***REMOVED***
	return newKafkaVersion(major, minor, veryMinor, patch), nil
***REMOVED***

func scanKafkaVersion(s string, pattern string, format string, v [3]*uint) error ***REMOVED***
	if !regexp.MustCompile(pattern).MatchString(s) ***REMOVED***
		return fmt.Errorf("invalid version `%s`", s)
	***REMOVED***
	_, err := fmt.Sscanf(s, format, v[0], v[1], v[2])
	return err
***REMOVED***

func (v KafkaVersion) String() string ***REMOVED***
	if v.version[0] == 0 ***REMOVED***
		return fmt.Sprintf("0.%d.%d.%d", v.version[1], v.version[2], v.version[3])
	***REMOVED*** else ***REMOVED***
		return fmt.Sprintf("%d.%d.%d", v.version[0], v.version[1], v.version[2])
	***REMOVED***
***REMOVED***
