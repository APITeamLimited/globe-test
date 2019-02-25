// Copyright 2013 Ooyala, Inc.

/*
Package statsd provides a Go dogstatsd client. Dogstatsd extends the popular statsd,
adding tags and histograms and pushing upstream to Datadog.

Refer to http://docs.datadoghq.com/guides/dogstatsd/ for information about DogStatsD.

Example Usage:

    // Create the client
    c, err := statsd.New("127.0.0.1:8125")
    if err != nil ***REMOVED***
        log.Fatal(err)
    ***REMOVED***
    // Prefix every metric with the app name
    c.Namespace = "flubber."
    // Send the EC2 availability zone as a tag with every metric
    c.Tags = append(c.Tags, "us-east-1a")
    err = c.Gauge("request.duration", 1.2, nil, 1)

statsd is based on go-statsd-client.
*/
package statsd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
OptimalPayloadSize defines the optimal payload size for a UDP datagram, 1432 bytes
is optimal for regular networks with an MTU of 1500 so datagrams don't get
fragmented. It's generally recommended not to fragment UDP datagrams as losing
a single fragment will cause the entire datagram to be lost.

This can be increased if your network has a greater MTU or you don't mind UDP
datagrams getting fragmented. The practical limit is MaxUDPPayloadSize
*/
const OptimalPayloadSize = 1432

/*
MaxUDPPayloadSize defines the maximum payload size for a UDP datagram.
Its value comes from the calculation: 65535 bytes Max UDP datagram size -
8byte UDP header - 60byte max IP headers
any number greater than that will see frames being cut out.
*/
const MaxUDPPayloadSize = 65467

/*
UnixAddressPrefix holds the prefix to use to enable Unix Domain Socket
traffic instead of UDP.
*/
const UnixAddressPrefix = "unix://"

/*
Stat suffixes
*/
var (
	gaugeSuffix        = []byte("|g")
	countSuffix        = []byte("|c")
	histogramSuffix    = []byte("|h")
	distributionSuffix = []byte("|d")
	decrSuffix         = []byte("-1|c")
	incrSuffix         = []byte("1|c")
	setSuffix          = []byte("|s")
	timingSuffix       = []byte("|ms")
)

// A statsdWriter offers a standard interface regardless of the underlying
// protocol. For now UDS and UPD writers are available.
type statsdWriter interface ***REMOVED***
	Write(data []byte) (n int, err error)
	SetWriteTimeout(time.Duration) error
	Close() error
***REMOVED***

// A Client is a handle for sending messages to dogstatsd.  It is safe to
// use one Client from multiple goroutines simultaneously.
type Client struct ***REMOVED***
	// Writer handles the underlying networking protocol
	writer statsdWriter
	// Namespace to prepend to all statsd calls
	Namespace string
	// Tags are global tags to be added to every statsd call
	Tags []string
	// skipErrors turns off error passing and allows UDS to emulate UDP behaviour
	SkipErrors bool
	// BufferLength is the length of the buffer in commands.
	bufferLength int
	flushTime    time.Duration
	commands     []string
	buffer       bytes.Buffer
	stop         chan struct***REMOVED******REMOVED***
	sync.Mutex
***REMOVED***

// New returns a pointer to a new Client given an addr in the format "hostname:port" or
// "unix:///path/to/socket".
func New(addr string) (*Client, error) ***REMOVED***
	if strings.HasPrefix(addr, UnixAddressPrefix) ***REMOVED***
		w, err := newUdsWriter(addr[len(UnixAddressPrefix)-1:])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return NewWithWriter(w)
	***REMOVED***
	w, err := newUDPWriter(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewWithWriter(w)
***REMOVED***

// NewWithWriter creates a new Client with given writer. Writer is a
// io.WriteCloser + SetWriteTimeout(time.Duration) error
func NewWithWriter(w statsdWriter) (*Client, error) ***REMOVED***
	client := &Client***REMOVED***writer: w, SkipErrors: false***REMOVED***
	return client, nil
***REMOVED***

// NewBuffered returns a Client that buffers its output and sends it in chunks.
// Buflen is the length of the buffer in number of commands.
func NewBuffered(addr string, buflen int) (*Client, error) ***REMOVED***
	client, err := New(addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	client.bufferLength = buflen
	client.commands = make([]string, 0, buflen)
	client.flushTime = time.Millisecond * 100
	client.stop = make(chan struct***REMOVED******REMOVED***, 1)
	go client.watch()
	return client, nil
***REMOVED***

// format a message from its name, value, tags and rate.  Also adds global
// namespace and tags.
func (c *Client) format(name string, value interface***REMOVED******REMOVED***, suffix []byte, tags []string, rate float64) string ***REMOVED***
	var buf bytes.Buffer
	if c.Namespace != "" ***REMOVED***
		buf.WriteString(c.Namespace)
	***REMOVED***
	buf.WriteString(name)
	buf.WriteString(":")

	switch val := value.(type) ***REMOVED***
	case float64:
		buf.Write(strconv.AppendFloat([]byte***REMOVED******REMOVED***, val, 'f', 6, 64))

	case int64:
		buf.Write(strconv.AppendInt([]byte***REMOVED******REMOVED***, val, 10))

	case string:
		buf.WriteString(val)

	default:
		// do nothing
	***REMOVED***
	buf.Write(suffix)

	if rate < 1 ***REMOVED***
		buf.WriteString(`|@`)
		buf.WriteString(strconv.FormatFloat(rate, 'f', -1, 64))
	***REMOVED***

	writeTagString(&buf, c.Tags, tags)

	return buf.String()
***REMOVED***

// SetWriteTimeout allows the user to set a custom UDS write timeout. Not supported for UDP.
func (c *Client) SetWriteTimeout(d time.Duration) error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	return c.writer.SetWriteTimeout(d)
***REMOVED***

func (c *Client) watch() ***REMOVED***
	ticker := time.NewTicker(c.flushTime)

	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			c.Lock()
			if len(c.commands) > 0 ***REMOVED***
				// FIXME: eating error here
				c.flushLocked()
			***REMOVED***
			c.Unlock()
		case <-c.stop:
			ticker.Stop()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Client) append(cmd string) error ***REMOVED***
	c.Lock()
	defer c.Unlock()
	c.commands = append(c.commands, cmd)
	// if we should flush, lets do it
	if len(c.commands) == c.bufferLength ***REMOVED***
		if err := c.flushLocked(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c *Client) joinMaxSize(cmds []string, sep string, maxSize int) ([][]byte, []int) ***REMOVED***
	c.buffer.Reset() //clear buffer

	var frames [][]byte
	var ncmds []int
	sepBytes := []byte(sep)
	sepLen := len(sep)

	elem := 0
	for _, cmd := range cmds ***REMOVED***
		needed := len(cmd)

		if elem != 0 ***REMOVED***
			needed = needed + sepLen
		***REMOVED***

		if c.buffer.Len()+needed <= maxSize ***REMOVED***
			if elem != 0 ***REMOVED***
				c.buffer.Write(sepBytes)
			***REMOVED***
			c.buffer.WriteString(cmd)
			elem++
		***REMOVED*** else ***REMOVED***
			frames = append(frames, copyAndResetBuffer(&c.buffer))
			ncmds = append(ncmds, elem)
			// if cmd is bigger than maxSize it will get flushed on next loop
			c.buffer.WriteString(cmd)
			elem = 1
		***REMOVED***
	***REMOVED***

	//add whatever is left! if there's actually something
	if c.buffer.Len() > 0 ***REMOVED***
		frames = append(frames, copyAndResetBuffer(&c.buffer))
		ncmds = append(ncmds, elem)
	***REMOVED***

	return frames, ncmds
***REMOVED***

func copyAndResetBuffer(buf *bytes.Buffer) []byte ***REMOVED***
	tmpBuf := make([]byte, buf.Len())
	copy(tmpBuf, buf.Bytes())
	buf.Reset()
	return tmpBuf
***REMOVED***

// Flush forces a flush of the pending commands in the buffer
func (c *Client) Flush() error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	c.Lock()
	defer c.Unlock()
	return c.flushLocked()
***REMOVED***

// flush the commands in the buffer.  Lock must be held by caller.
func (c *Client) flushLocked() error ***REMOVED***
	frames, flushable := c.joinMaxSize(c.commands, "\n", OptimalPayloadSize)
	var err error
	cmdsFlushed := 0
	for i, data := range frames ***REMOVED***
		_, e := c.writer.Write(data)
		if e != nil ***REMOVED***
			err = e
			break
		***REMOVED***
		cmdsFlushed += flushable[i]
	***REMOVED***

	// clear the slice with a slice op, doesn't realloc
	if cmdsFlushed == len(c.commands) ***REMOVED***
		c.commands = c.commands[:0]
	***REMOVED*** else ***REMOVED***
		//this case will cause a future realloc...
		// drop problematic command though (sorry).
		c.commands = c.commands[cmdsFlushed+1:]
	***REMOVED***
	return err
***REMOVED***

func (c *Client) sendMsg(msg string) error ***REMOVED***
	// return an error if message is bigger than MaxUDPPayloadSize
	if len(msg) > MaxUDPPayloadSize ***REMOVED***
		return errors.New("message size exceeds MaxUDPPayloadSize")
	***REMOVED***

	// if this client is buffered, then we'll just append this
	if c.bufferLength > 0 ***REMOVED***
		return c.append(msg)
	***REMOVED***

	_, err := c.writer.Write([]byte(msg))

	if c.SkipErrors ***REMOVED***
		return nil
	***REMOVED***
	return err
***REMOVED***

// send handles sampling and sends the message over UDP. It also adds global namespace prefixes and tags.
func (c *Client) send(name string, value interface***REMOVED******REMOVED***, suffix []byte, tags []string, rate float64) error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	if rate < 1 && rand.Float64() > rate ***REMOVED***
		return nil
	***REMOVED***
	data := c.format(name, value, suffix, tags, rate)
	return c.sendMsg(data)
***REMOVED***

// Gauge measures the value of a metric at a particular time.
func (c *Client) Gauge(name string, value float64, tags []string, rate float64) error ***REMOVED***
	return c.send(name, value, gaugeSuffix, tags, rate)
***REMOVED***

// Count tracks how many times something happened per second.
func (c *Client) Count(name string, value int64, tags []string, rate float64) error ***REMOVED***
	return c.send(name, value, countSuffix, tags, rate)
***REMOVED***

// Histogram tracks the statistical distribution of a set of values on each host.
func (c *Client) Histogram(name string, value float64, tags []string, rate float64) error ***REMOVED***
	return c.send(name, value, histogramSuffix, tags, rate)
***REMOVED***

// Distribution tracks the statistical distribution of a set of values across your infrastructure.
func (c *Client) Distribution(name string, value float64, tags []string, rate float64) error ***REMOVED***
	return c.send(name, value, distributionSuffix, tags, rate)
***REMOVED***

// Decr is just Count of -1
func (c *Client) Decr(name string, tags []string, rate float64) error ***REMOVED***
	return c.send(name, nil, decrSuffix, tags, rate)
***REMOVED***

// Incr is just Count of 1
func (c *Client) Incr(name string, tags []string, rate float64) error ***REMOVED***
	return c.send(name, nil, incrSuffix, tags, rate)
***REMOVED***

// Set counts the number of unique elements in a group.
func (c *Client) Set(name string, value string, tags []string, rate float64) error ***REMOVED***
	return c.send(name, value, setSuffix, tags, rate)
***REMOVED***

// Timing sends timing information, it is an alias for TimeInMilliseconds
func (c *Client) Timing(name string, value time.Duration, tags []string, rate float64) error ***REMOVED***
	return c.TimeInMilliseconds(name, value.Seconds()*1000, tags, rate)
***REMOVED***

// TimeInMilliseconds sends timing information in milliseconds.
// It is flushed by statsd with percentiles, mean and other info (https://github.com/etsy/statsd/blob/master/docs/metric_types.md#timing)
func (c *Client) TimeInMilliseconds(name string, value float64, tags []string, rate float64) error ***REMOVED***
	return c.send(name, value, timingSuffix, tags, rate)
***REMOVED***

// Event sends the provided Event.
func (c *Client) Event(e *Event) error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	stat, err := e.Encode(c.Tags...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return c.sendMsg(stat)
***REMOVED***

// SimpleEvent sends an event with the provided title and text.
func (c *Client) SimpleEvent(title, text string) error ***REMOVED***
	e := NewEvent(title, text)
	return c.Event(e)
***REMOVED***

// ServiceCheck sends the provided ServiceCheck.
func (c *Client) ServiceCheck(sc *ServiceCheck) error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	stat, err := sc.Encode(c.Tags...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return c.sendMsg(stat)
***REMOVED***

// SimpleServiceCheck sends an serviceCheck with the provided name and status.
func (c *Client) SimpleServiceCheck(name string, status ServiceCheckStatus) error ***REMOVED***
	sc := NewServiceCheck(name, status)
	return c.ServiceCheck(sc)
***REMOVED***

// Close the client connection.
func (c *Client) Close() error ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	select ***REMOVED***
	case c.stop <- struct***REMOVED******REMOVED******REMOVED******REMOVED***:
	default:
	***REMOVED***

	// if this client is buffered, flush before closing the writer
	if c.bufferLength > 0 ***REMOVED***
		if err := c.Flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return c.writer.Close()
***REMOVED***

// Events support
// EventAlertType and EventAlertPriority became exported types after this issue was submitted: https://github.com/DataDog/datadog-go/issues/41
// The reason why they got exported is so that client code can directly use the types.

// EventAlertType is the alert type for events
type EventAlertType string

const (
	// Info is the "info" AlertType for events
	Info EventAlertType = "info"
	// Error is the "error" AlertType for events
	Error EventAlertType = "error"
	// Warning is the "warning" AlertType for events
	Warning EventAlertType = "warning"
	// Success is the "success" AlertType for events
	Success EventAlertType = "success"
)

// EventPriority is the event priority for events
type EventPriority string

const (
	// Normal is the "normal" Priority for events
	Normal EventPriority = "normal"
	// Low is the "low" Priority for events
	Low EventPriority = "low"
)

// An Event is an object that can be posted to your DataDog event stream.
type Event struct ***REMOVED***
	// Title of the event.  Required.
	Title string
	// Text is the description of the event.  Required.
	Text string
	// Timestamp is a timestamp for the event.  If not provided, the dogstatsd
	// server will set this to the current time.
	Timestamp time.Time
	// Hostname for the event.
	Hostname string
	// AggregationKey groups this event with others of the same key.
	AggregationKey string
	// Priority of the event.  Can be statsd.Low or statsd.Normal.
	Priority EventPriority
	// SourceTypeName is a source type for the event.
	SourceTypeName string
	// AlertType can be statsd.Info, statsd.Error, statsd.Warning, or statsd.Success.
	// If absent, the default value applied by the dogstatsd server is Info.
	AlertType EventAlertType
	// Tags for the event.
	Tags []string
***REMOVED***

// NewEvent creates a new event with the given title and text.  Error checking
// against these values is done at send-time, or upon running e.Check.
func NewEvent(title, text string) *Event ***REMOVED***
	return &Event***REMOVED***
		Title: title,
		Text:  text,
	***REMOVED***
***REMOVED***

// Check verifies that an event is valid.
func (e Event) Check() error ***REMOVED***
	if len(e.Title) == 0 ***REMOVED***
		return fmt.Errorf("statsd.Event title is required")
	***REMOVED***
	if len(e.Text) == 0 ***REMOVED***
		return fmt.Errorf("statsd.Event text is required")
	***REMOVED***
	return nil
***REMOVED***

// Encode returns the dogstatsd wire protocol representation for an event.
// Tags may be passed which will be added to the encoded output but not to
// the Event's list of tags, eg. for default tags.
func (e Event) Encode(tags ...string) (string, error) ***REMOVED***
	err := e.Check()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	text := e.escapedText()

	var buffer bytes.Buffer
	buffer.WriteString("_e***REMOVED***")
	buffer.WriteString(strconv.FormatInt(int64(len(e.Title)), 10))
	buffer.WriteRune(',')
	buffer.WriteString(strconv.FormatInt(int64(len(text)), 10))
	buffer.WriteString("***REMOVED***:")
	buffer.WriteString(e.Title)
	buffer.WriteRune('|')
	buffer.WriteString(text)

	if !e.Timestamp.IsZero() ***REMOVED***
		buffer.WriteString("|d:")
		buffer.WriteString(strconv.FormatInt(int64(e.Timestamp.Unix()), 10))
	***REMOVED***

	if len(e.Hostname) != 0 ***REMOVED***
		buffer.WriteString("|h:")
		buffer.WriteString(e.Hostname)
	***REMOVED***

	if len(e.AggregationKey) != 0 ***REMOVED***
		buffer.WriteString("|k:")
		buffer.WriteString(e.AggregationKey)

	***REMOVED***

	if len(e.Priority) != 0 ***REMOVED***
		buffer.WriteString("|p:")
		buffer.WriteString(string(e.Priority))
	***REMOVED***

	if len(e.SourceTypeName) != 0 ***REMOVED***
		buffer.WriteString("|s:")
		buffer.WriteString(e.SourceTypeName)
	***REMOVED***

	if len(e.AlertType) != 0 ***REMOVED***
		buffer.WriteString("|t:")
		buffer.WriteString(string(e.AlertType))
	***REMOVED***

	writeTagString(&buffer, tags, e.Tags)

	return buffer.String(), nil
***REMOVED***

// ServiceCheckStatus support
type ServiceCheckStatus byte

const (
	// Ok is the "ok" ServiceCheck status
	Ok ServiceCheckStatus = 0
	// Warn is the "warning" ServiceCheck status
	Warn ServiceCheckStatus = 1
	// Critical is the "critical" ServiceCheck status
	Critical ServiceCheckStatus = 2
	// Unknown is the "unknown" ServiceCheck status
	Unknown ServiceCheckStatus = 3
)

// An ServiceCheck is an object that contains status of DataDog service check.
type ServiceCheck struct ***REMOVED***
	// Name of the service check.  Required.
	Name string
	// Status of service check.  Required.
	Status ServiceCheckStatus
	// Timestamp is a timestamp for the serviceCheck.  If not provided, the dogstatsd
	// server will set this to the current time.
	Timestamp time.Time
	// Hostname for the serviceCheck.
	Hostname string
	// A message describing the current state of the serviceCheck.
	Message string
	// Tags for the serviceCheck.
	Tags []string
***REMOVED***

// NewServiceCheck creates a new serviceCheck with the given name and status.  Error checking
// against these values is done at send-time, or upon running sc.Check.
func NewServiceCheck(name string, status ServiceCheckStatus) *ServiceCheck ***REMOVED***
	return &ServiceCheck***REMOVED***
		Name:   name,
		Status: status,
	***REMOVED***
***REMOVED***

// Check verifies that an event is valid.
func (sc ServiceCheck) Check() error ***REMOVED***
	if len(sc.Name) == 0 ***REMOVED***
		return fmt.Errorf("statsd.ServiceCheck name is required")
	***REMOVED***
	if byte(sc.Status) < 0 || byte(sc.Status) > 3 ***REMOVED***
		return fmt.Errorf("statsd.ServiceCheck status has invalid value")
	***REMOVED***
	return nil
***REMOVED***

// Encode returns the dogstatsd wire protocol representation for an serviceCheck.
// Tags may be passed which will be added to the encoded output but not to
// the Event's list of tags, eg. for default tags.
func (sc ServiceCheck) Encode(tags ...string) (string, error) ***REMOVED***
	err := sc.Check()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	message := sc.escapedMessage()

	var buffer bytes.Buffer
	buffer.WriteString("_sc|")
	buffer.WriteString(sc.Name)
	buffer.WriteRune('|')
	buffer.WriteString(strconv.FormatInt(int64(sc.Status), 10))

	if !sc.Timestamp.IsZero() ***REMOVED***
		buffer.WriteString("|d:")
		buffer.WriteString(strconv.FormatInt(int64(sc.Timestamp.Unix()), 10))
	***REMOVED***

	if len(sc.Hostname) != 0 ***REMOVED***
		buffer.WriteString("|h:")
		buffer.WriteString(sc.Hostname)
	***REMOVED***

	writeTagString(&buffer, tags, sc.Tags)

	if len(message) != 0 ***REMOVED***
		buffer.WriteString("|m:")
		buffer.WriteString(message)
	***REMOVED***

	return buffer.String(), nil
***REMOVED***

func (e Event) escapedText() string ***REMOVED***
	return strings.Replace(e.Text, "\n", "\\n", -1)
***REMOVED***

func (sc ServiceCheck) escapedMessage() string ***REMOVED***
	msg := strings.Replace(sc.Message, "\n", "\\n", -1)
	return strings.Replace(msg, "m:", `m\:`, -1)
***REMOVED***

func removeNewlines(str string) string ***REMOVED***
	return strings.Replace(str, "\n", "", -1)
***REMOVED***

func writeTagString(w io.Writer, tagList1, tagList2 []string) ***REMOVED***
	// the tag lists may be shared with other callers, so we cannot modify
	// them in any way (which means we cannot append to them either)
	// therefore we must make an entirely separate copy just for this call
	totalLen := len(tagList1) + len(tagList2)
	if totalLen == 0 ***REMOVED***
		return
	***REMOVED***
	tags := make([]string, 0, totalLen)
	tags = append(tags, tagList1...)
	tags = append(tags, tagList2...)

	io.WriteString(w, "|#")
	io.WriteString(w, removeNewlines(tags[0]))
	for _, tag := range tags[1:] ***REMOVED***
		io.WriteString(w, ",")
		io.WriteString(w, removeNewlines(tag))
	***REMOVED***
***REMOVED***
