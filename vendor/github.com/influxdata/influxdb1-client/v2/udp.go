package client

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	// UDPPayloadSize is a reasonable default payload size for UDP packets that
	// could be travelling over the internet.
	UDPPayloadSize = 512
)

// UDPConfig is the config data needed to create a UDP Client.
type UDPConfig struct ***REMOVED***
	// Addr should be of the form "host:port"
	// or "[ipv6-host%zone]:port".
	Addr string

	// PayloadSize is the maximum size of a UDP client message, optional
	// Tune this based on your network. Defaults to UDPPayloadSize.
	PayloadSize int
***REMOVED***

// NewUDPClient returns a client interface for writing to an InfluxDB UDP
// service from the given config.
func NewUDPClient(conf UDPConfig) (Client, error) ***REMOVED***
	var udpAddr *net.UDPAddr
	udpAddr, err := net.ResolveUDPAddr("udp", conf.Addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	payloadSize := conf.PayloadSize
	if payloadSize == 0 ***REMOVED***
		payloadSize = UDPPayloadSize
	***REMOVED***

	return &udpclient***REMOVED***
		conn:        conn,
		payloadSize: payloadSize,
	***REMOVED***, nil
***REMOVED***

// Close releases the udpclient's resources.
func (uc *udpclient) Close() error ***REMOVED***
	return uc.conn.Close()
***REMOVED***

type udpclient struct ***REMOVED***
	conn        io.WriteCloser
	payloadSize int
***REMOVED***

func (uc *udpclient) Write(bp BatchPoints) error ***REMOVED***
	var b = make([]byte, 0, uc.payloadSize) // initial buffer size, it will grow as needed
	var d, _ = time.ParseDuration("1" + bp.Precision())

	var delayedError error

	var checkBuffer = func(n int) ***REMOVED***
		if len(b) > 0 && len(b)+n > uc.payloadSize ***REMOVED***
			if _, err := uc.conn.Write(b); err != nil ***REMOVED***
				delayedError = err
			***REMOVED***
			b = b[:0]
		***REMOVED***
	***REMOVED***

	for _, p := range bp.Points() ***REMOVED***
		p.pt.Round(d)
		pointSize := p.pt.StringSize() + 1 // include newline in size
		//point := p.pt.RoundedString(d) + "\n"

		checkBuffer(pointSize)

		if p.Time().IsZero() || pointSize <= uc.payloadSize ***REMOVED***
			b = p.pt.AppendString(b)
			b = append(b, '\n')
			continue
		***REMOVED***

		points := p.pt.Split(uc.payloadSize - 1) // account for newline character
		for _, sp := range points ***REMOVED***
			checkBuffer(sp.StringSize() + 1)
			b = sp.AppendString(b)
			b = append(b, '\n')
		***REMOVED***
	***REMOVED***

	if len(b) > 0 ***REMOVED***
		if _, err := uc.conn.Write(b); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return delayedError
***REMOVED***

func (uc *udpclient) Query(q Query) (*Response, error) ***REMOVED***
	return nil, fmt.Errorf("Querying via UDP is not supported")
***REMOVED***

func (uc *udpclient) QueryAsChunk(q Query) (*ChunkedResponse, error) ***REMOVED***
	return nil, fmt.Errorf("Querying via UDP is not supported")
***REMOVED***

func (uc *udpclient) Ping(timeout time.Duration) (time.Duration, string, error) ***REMOVED***
	return 0, "", nil
***REMOVED***
