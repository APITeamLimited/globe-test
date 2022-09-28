package netext

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/libWorker/types"
	"github.com/APITeamLimited/globe-test/worker/metrics"
)

// Dialer wraps net.Dialer and provides k6 specific functionality -
// tracing, blacklists and DNS cache and aliases.
type Dialer struct ***REMOVED***
	net.Dialer

	Resolver         Resolver
	Blacklist        []*libWorker.IPNet
	BlockedHostnames *types.HostnameTrie
	Hosts            map[string]*libWorker.HostAddress

	BytesRead    int64
	BytesWritten int64
***REMOVED***

// NewDialer constructs a new Dialer with the given DNS resolver.
func NewDialer(dialer net.Dialer, resolver Resolver) *Dialer ***REMOVED***
	return &Dialer***REMOVED***
		Dialer:   dialer,
		Resolver: resolver,
	***REMOVED***
***REMOVED***

// BlackListedIPError is an error that is returned when a given IP is blacklisted
type BlackListedIPError struct ***REMOVED***
	ip  net.IP
	net *libWorker.IPNet
***REMOVED***

func (b BlackListedIPError) Error() string ***REMOVED***
	return fmt.Sprintf("IP (%s) is in a blacklisted range (%s)", b.ip, b.net)
***REMOVED***

// BlockedHostError is returned when a given hostname is blocked
type BlockedHostError struct ***REMOVED***
	hostname string
	match    string
***REMOVED***

func (b BlockedHostError) Error() string ***REMOVED***
	return fmt.Sprintf("hostname (%s) is in a blocked pattern (%s)", b.hostname, b.match)
***REMOVED***

// DialContext wraps the net.Dialer.DialContext and handles the k6 specifics
func (d *Dialer) DialContext(ctx context.Context, proto, addr string) (net.Conn, error) ***REMOVED***
	dialAddr, err := d.getDialAddr(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	conn, err := d.Dialer.DialContext(ctx, proto, dialAddr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	conn = &Conn***REMOVED***conn, &d.BytesRead, &d.BytesWritten***REMOVED***
	return conn, err
***REMOVED***

// GetTrail creates a new NetTrail instance with the Dialer
// sent and received data metrics and the supplied times and tags.
// TODO: Refactor this according to
// https://github.com/k6io/k6/pull/1203#discussion_r337938370
func (d *Dialer) GetTrail(
	startTime, endTime time.Time, fullIteration bool, emitIterations bool, tags *metrics.SampleTags,
	builtinMetrics *metrics.BuiltinMetrics,
) *NetTrail ***REMOVED***
	bytesWritten := atomic.SwapInt64(&d.BytesWritten, 0)
	bytesRead := atomic.SwapInt64(&d.BytesRead, 0)
	samples := []metrics.Sample***REMOVED***
		***REMOVED***
			Time:   endTime,
			Metric: builtinMetrics.DataSent,
			Value:  float64(bytesWritten),
			Tags:   tags,
		***REMOVED***,
		***REMOVED***
			Time:   endTime,
			Metric: builtinMetrics.DataReceived,
			Value:  float64(bytesRead),
			Tags:   tags,
		***REMOVED***,
	***REMOVED***
	if fullIteration ***REMOVED***
		samples = append(samples, metrics.Sample***REMOVED***
			Time:   endTime,
			Metric: builtinMetrics.IterationDuration,
			Value:  metrics.D(endTime.Sub(startTime)),
			Tags:   tags,
		***REMOVED***)
		if emitIterations ***REMOVED***
			samples = append(samples, metrics.Sample***REMOVED***
				Time:   endTime,
				Metric: builtinMetrics.Iterations,
				Value:  1,
				Tags:   tags,
			***REMOVED***)
		***REMOVED***
	***REMOVED***

	return &NetTrail***REMOVED***
		BytesRead:     bytesRead,
		BytesWritten:  bytesWritten,
		FullIteration: fullIteration,
		StartTime:     startTime,
		EndTime:       endTime,
		Tags:          tags,
		Samples:       samples,
	***REMOVED***
***REMOVED***

func (d *Dialer) getDialAddr(addr string) (string, error) ***REMOVED***
	remote, err := d.findRemote(addr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	for _, ipnet := range d.Blacklist ***REMOVED***
		if ipnet.Contains(remote.IP) ***REMOVED***
			return "", BlackListedIPError***REMOVED***ip: remote.IP, net: ipnet***REMOVED***
		***REMOVED***
	***REMOVED***

	return remote.String(), nil
***REMOVED***

func (d *Dialer) findRemote(addr string) (*libWorker.HostAddress, error) ***REMOVED***
	host, port, err := net.SplitHostPort(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ip := net.ParseIP(host)
	if d.BlockedHostnames != nil && ip == nil ***REMOVED***
		if match, blocked := d.BlockedHostnames.Contains(host); blocked ***REMOVED***
			return nil, BlockedHostError***REMOVED***hostname: host, match: match***REMOVED***
		***REMOVED***
	***REMOVED***

	remote, err := d.getConfiguredHost(addr, host, port)
	if err != nil || remote != nil ***REMOVED***
		return remote, err
	***REMOVED***

	if ip != nil ***REMOVED***
		return libWorker.NewHostAddress(ip, port)
	***REMOVED***

	ip, err = d.Resolver.LookupIP(host)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if ip == nil ***REMOVED***
		return nil, fmt.Errorf("lookup %s: no such host", host)
	***REMOVED***

	return libWorker.NewHostAddress(ip, port)
***REMOVED***

func (d *Dialer) getConfiguredHost(addr, host, port string) (*libWorker.HostAddress, error) ***REMOVED***
	if remote, ok := d.Hosts[addr]; ok ***REMOVED***
		return remote, nil
	***REMOVED***

	if remote, ok := d.Hosts[host]; ok ***REMOVED***
		if remote.Port != 0 || port == "" ***REMOVED***
			return remote, nil
		***REMOVED***

		newPort, err := strconv.Atoi(port)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		newRemote := *remote
		newRemote.Port = newPort

		return &newRemote, nil
	***REMOVED***

	return nil, nil
***REMOVED***

// NetTrail contains information about the exchanged data size and length of a
// series of connections from a particular netext.Dialer
type NetTrail struct ***REMOVED***
	BytesRead     int64
	BytesWritten  int64
	FullIteration bool
	StartTime     time.Time
	EndTime       time.Time
	Tags          *metrics.SampleTags
	Samples       []metrics.Sample
***REMOVED***

// Ensure that interfaces are implemented correctly
var _ metrics.ConnectedSampleContainer = &NetTrail***REMOVED******REMOVED***

// GetSamples implements the metrics.SampleContainer interface.
func (ntr *NetTrail) GetSamples() []metrics.Sample ***REMOVED***
	return ntr.Samples
***REMOVED***

// GetTags implements the metrics.ConnectedSampleContainer interface.
func (ntr *NetTrail) GetTags() *metrics.SampleTags ***REMOVED***
	return ntr.Tags
***REMOVED***

// GetTime implements the metrics.ConnectedSampleContainer interface.
func (ntr *NetTrail) GetTime() time.Time ***REMOVED***
	return ntr.EndTime
***REMOVED***

// Conn wraps net.Conn and keeps track of sent and received data size
type Conn struct ***REMOVED***
	net.Conn

	BytesRead, BytesWritten *int64
***REMOVED***

func (c *Conn) Read(b []byte) (int, error) ***REMOVED***
	n, err := c.Conn.Read(b)
	if n > 0 ***REMOVED***
		atomic.AddInt64(c.BytesRead, int64(n))
	***REMOVED***
	return n, err
***REMOVED***

func (c *Conn) Write(b []byte) (int, error) ***REMOVED***
	n, err := c.Conn.Write(b)
	if n > 0 ***REMOVED***
		atomic.AddInt64(c.BytesWritten, int64(n))
	***REMOVED***
	return n, err
***REMOVED***
