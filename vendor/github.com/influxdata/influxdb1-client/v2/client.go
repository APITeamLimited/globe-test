// Package client (v2) is the current official Go client for InfluxDB.
package client // import "github.com/influxdata/influxdb1-client/v2"

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb1-client/models"
)

// HTTPConfig is the config data needed to create an HTTP Client.
type HTTPConfig struct ***REMOVED***
	// Addr should be of the form "http://host:port"
	// or "http://[ipv6-host%zone]:port".
	Addr string

	// Username is the influxdb username, optional.
	Username string

	// Password is the influxdb password, optional.
	Password string

	// UserAgent is the http User Agent, defaults to "InfluxDBClient".
	UserAgent string

	// Timeout for influxdb writes, defaults to no timeout.
	Timeout time.Duration

	// InsecureSkipVerify gets passed to the http client, if true, it will
	// skip https certificate verification. Defaults to false.
	InsecureSkipVerify bool

	// TLSConfig allows the user to set their own TLS config for the HTTP
	// Client. If set, this option overrides InsecureSkipVerify.
	TLSConfig *tls.Config

	// Proxy configures the Proxy function on the HTTP client.
	Proxy func(req *http.Request) (*url.URL, error)
***REMOVED***

// BatchPointsConfig is the config data needed to create an instance of the BatchPoints struct.
type BatchPointsConfig struct ***REMOVED***
	// Precision is the write precision of the points, defaults to "ns".
	Precision string

	// Database is the database to write points to.
	Database string

	// RetentionPolicy is the retention policy of the points.
	RetentionPolicy string

	// Write consistency is the number of servers required to confirm write.
	WriteConsistency string
***REMOVED***

// Client is a client interface for writing & querying the database.
type Client interface ***REMOVED***
	// Ping checks that status of cluster, and will always return 0 time and no
	// error for UDP clients.
	Ping(timeout time.Duration) (time.Duration, string, error)

	// Write takes a BatchPoints object and writes all Points to InfluxDB.
	Write(bp BatchPoints) error

	// Query makes an InfluxDB Query on the database. This will fail if using
	// the UDP client.
	Query(q Query) (*Response, error)

	// QueryAsChunk makes an InfluxDB Query on the database. This will fail if using
	// the UDP client.
	QueryAsChunk(q Query) (*ChunkedResponse, error)

	// Close releases any resources a Client may be using.
	Close() error
***REMOVED***

// NewHTTPClient returns a new Client from the provided config.
// Client is safe for concurrent use by multiple goroutines.
func NewHTTPClient(conf HTTPConfig) (Client, error) ***REMOVED***
	if conf.UserAgent == "" ***REMOVED***
		conf.UserAgent = "InfluxDBClient"
	***REMOVED***

	u, err := url.Parse(conf.Addr)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if u.Scheme != "http" && u.Scheme != "https" ***REMOVED***
		m := fmt.Sprintf("Unsupported protocol scheme: %s, your address"+
			" must start with http:// or https://", u.Scheme)
		return nil, errors.New(m)
	***REMOVED***

	tr := &http.Transport***REMOVED***
		TLSClientConfig: &tls.Config***REMOVED***
			InsecureSkipVerify: conf.InsecureSkipVerify,
		***REMOVED***,
		Proxy: conf.Proxy,
	***REMOVED***
	if conf.TLSConfig != nil ***REMOVED***
		tr.TLSClientConfig = conf.TLSConfig
	***REMOVED***
	return &client***REMOVED***
		url:       *u,
		username:  conf.Username,
		password:  conf.Password,
		useragent: conf.UserAgent,
		httpClient: &http.Client***REMOVED***
			Timeout:   conf.Timeout,
			Transport: tr,
		***REMOVED***,
		transport: tr,
	***REMOVED***, nil
***REMOVED***

// Ping will check to see if the server is up with an optional timeout on waiting for leader.
// Ping returns how long the request took, the version of the server it connected to, and an error if one occurred.
func (c *client) Ping(timeout time.Duration) (time.Duration, string, error) ***REMOVED***
	now := time.Now()

	u := c.url
	u.Path = path.Join(u.Path, "ping")

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil ***REMOVED***
		return 0, "", err
	***REMOVED***

	req.Header.Set("User-Agent", c.useragent)

	if c.username != "" ***REMOVED***
		req.SetBasicAuth(c.username, c.password)
	***REMOVED***

	if timeout > 0 ***REMOVED***
		params := req.URL.Query()
		params.Set("wait_for_leader", fmt.Sprintf("%.0fs", timeout.Seconds()))
		req.URL.RawQuery = params.Encode()
	***REMOVED***

	resp, err := c.httpClient.Do(req)
	if err != nil ***REMOVED***
		return 0, "", err
	***REMOVED***
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil ***REMOVED***
		return 0, "", err
	***REMOVED***

	if resp.StatusCode != http.StatusNoContent ***REMOVED***
		var err = errors.New(string(body))
		return 0, "", err
	***REMOVED***

	version := resp.Header.Get("X-Influxdb-Version")
	return time.Since(now), version, nil
***REMOVED***

// Close releases the client's resources.
func (c *client) Close() error ***REMOVED***
	c.transport.CloseIdleConnections()
	return nil
***REMOVED***

// client is safe for concurrent use as the fields are all read-only
// once the client is instantiated.
type client struct ***REMOVED***
	// N.B - if url.UserInfo is accessed in future modifications to the
	// methods on client, you will need to synchronize access to url.
	url        url.URL
	username   string
	password   string
	useragent  string
	httpClient *http.Client
	transport  *http.Transport
***REMOVED***

// BatchPoints is an interface into a batched grouping of points to write into
// InfluxDB together. BatchPoints is NOT thread-safe, you must create a separate
// batch for each goroutine.
type BatchPoints interface ***REMOVED***
	// AddPoint adds the given point to the Batch of points.
	AddPoint(p *Point)
	// AddPoints adds the given points to the Batch of points.
	AddPoints(ps []*Point)
	// Points lists the points in the Batch.
	Points() []*Point

	// Precision returns the currently set precision of this Batch.
	Precision() string
	// SetPrecision sets the precision of this batch.
	SetPrecision(s string) error

	// Database returns the currently set database of this Batch.
	Database() string
	// SetDatabase sets the database of this Batch.
	SetDatabase(s string)

	// WriteConsistency returns the currently set write consistency of this Batch.
	WriteConsistency() string
	// SetWriteConsistency sets the write consistency of this Batch.
	SetWriteConsistency(s string)

	// RetentionPolicy returns the currently set retention policy of this Batch.
	RetentionPolicy() string
	// SetRetentionPolicy sets the retention policy of this Batch.
	SetRetentionPolicy(s string)
***REMOVED***

// NewBatchPoints returns a BatchPoints interface based on the given config.
func NewBatchPoints(conf BatchPointsConfig) (BatchPoints, error) ***REMOVED***
	if conf.Precision == "" ***REMOVED***
		conf.Precision = "ns"
	***REMOVED***
	if _, err := time.ParseDuration("1" + conf.Precision); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	bp := &batchpoints***REMOVED***
		database:         conf.Database,
		precision:        conf.Precision,
		retentionPolicy:  conf.RetentionPolicy,
		writeConsistency: conf.WriteConsistency,
	***REMOVED***
	return bp, nil
***REMOVED***

type batchpoints struct ***REMOVED***
	points           []*Point
	database         string
	precision        string
	retentionPolicy  string
	writeConsistency string
***REMOVED***

func (bp *batchpoints) AddPoint(p *Point) ***REMOVED***
	bp.points = append(bp.points, p)
***REMOVED***

func (bp *batchpoints) AddPoints(ps []*Point) ***REMOVED***
	bp.points = append(bp.points, ps...)
***REMOVED***

func (bp *batchpoints) Points() []*Point ***REMOVED***
	return bp.points
***REMOVED***

func (bp *batchpoints) Precision() string ***REMOVED***
	return bp.precision
***REMOVED***

func (bp *batchpoints) Database() string ***REMOVED***
	return bp.database
***REMOVED***

func (bp *batchpoints) WriteConsistency() string ***REMOVED***
	return bp.writeConsistency
***REMOVED***

func (bp *batchpoints) RetentionPolicy() string ***REMOVED***
	return bp.retentionPolicy
***REMOVED***

func (bp *batchpoints) SetPrecision(p string) error ***REMOVED***
	if _, err := time.ParseDuration("1" + p); err != nil ***REMOVED***
		return err
	***REMOVED***
	bp.precision = p
	return nil
***REMOVED***

func (bp *batchpoints) SetDatabase(db string) ***REMOVED***
	bp.database = db
***REMOVED***

func (bp *batchpoints) SetWriteConsistency(wc string) ***REMOVED***
	bp.writeConsistency = wc
***REMOVED***

func (bp *batchpoints) SetRetentionPolicy(rp string) ***REMOVED***
	bp.retentionPolicy = rp
***REMOVED***

// Point represents a single data point.
type Point struct ***REMOVED***
	pt models.Point
***REMOVED***

// NewPoint returns a point with the given timestamp. If a timestamp is not
// given, then data is sent to the database without a timestamp, in which case
// the server will assign local time upon reception. NOTE: it is recommended to
// send data with a timestamp.
func NewPoint(
	name string,
	tags map[string]string,
	fields map[string]interface***REMOVED******REMOVED***,
	t ...time.Time,
) (*Point, error) ***REMOVED***
	var T time.Time
	if len(t) > 0 ***REMOVED***
		T = t[0]
	***REMOVED***

	pt, err := models.NewPoint(name, models.NewTags(tags), fields, T)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &Point***REMOVED***
		pt: pt,
	***REMOVED***, nil
***REMOVED***

// String returns a line-protocol string of the Point.
func (p *Point) String() string ***REMOVED***
	return p.pt.String()
***REMOVED***

// PrecisionString returns a line-protocol string of the Point,
// with the timestamp formatted for the given precision.
func (p *Point) PrecisionString(precision string) string ***REMOVED***
	return p.pt.PrecisionString(precision)
***REMOVED***

// Name returns the measurement name of the point.
func (p *Point) Name() string ***REMOVED***
	return string(p.pt.Name())
***REMOVED***

// Tags returns the tags associated with the point.
func (p *Point) Tags() map[string]string ***REMOVED***
	return p.pt.Tags().Map()
***REMOVED***

// Time return the timestamp for the point.
func (p *Point) Time() time.Time ***REMOVED***
	return p.pt.Time()
***REMOVED***

// UnixNano returns timestamp of the point in nanoseconds since Unix epoch.
func (p *Point) UnixNano() int64 ***REMOVED***
	return p.pt.UnixNano()
***REMOVED***

// Fields returns the fields for the point.
func (p *Point) Fields() (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return p.pt.Fields()
***REMOVED***

// NewPointFrom returns a point from the provided models.Point.
func NewPointFrom(pt models.Point) *Point ***REMOVED***
	return &Point***REMOVED***pt: pt***REMOVED***
***REMOVED***

func (c *client) Write(bp BatchPoints) error ***REMOVED***
	var b bytes.Buffer

	for _, p := range bp.Points() ***REMOVED***
		if p == nil ***REMOVED***
			continue
		***REMOVED***
		if _, err := b.WriteString(p.pt.PrecisionString(bp.Precision())); err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := b.WriteByte('\n'); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	u := c.url
	u.Path = path.Join(u.Path, "write")

	req, err := http.NewRequest("POST", u.String(), &b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.Header.Set("Content-Type", "")
	req.Header.Set("User-Agent", c.useragent)
	if c.username != "" ***REMOVED***
		req.SetBasicAuth(c.username, c.password)
	***REMOVED***

	params := req.URL.Query()
	params.Set("db", bp.Database())
	params.Set("rp", bp.RetentionPolicy())
	params.Set("precision", bp.Precision())
	params.Set("consistency", bp.WriteConsistency())
	req.URL.RawQuery = params.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK ***REMOVED***
		var err = errors.New(string(body))
		return err
	***REMOVED***

	return nil
***REMOVED***

// Query defines a query to send to the server.
type Query struct ***REMOVED***
	Command         string
	Database        string
	RetentionPolicy string
	Precision       string
	Chunked         bool
	ChunkSize       int
	Parameters      map[string]interface***REMOVED******REMOVED***
***REMOVED***

// NewQuery returns a query object.
// The database and precision arguments can be empty strings if they are not needed for the query.
func NewQuery(command, database, precision string) Query ***REMOVED***
	return Query***REMOVED***
		Command:    command,
		Database:   database,
		Precision:  precision,
		Parameters: make(map[string]interface***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// NewQueryWithRP returns a query object.
// The database, retention policy, and precision arguments can be empty strings if they are not needed
// for the query. Setting the retention policy only works on InfluxDB versions 1.6 or greater.
func NewQueryWithRP(command, database, retentionPolicy, precision string) Query ***REMOVED***
	return Query***REMOVED***
		Command:         command,
		Database:        database,
		RetentionPolicy: retentionPolicy,
		Precision:       precision,
		Parameters:      make(map[string]interface***REMOVED******REMOVED***),
	***REMOVED***
***REMOVED***

// NewQueryWithParameters returns a query object.
// The database and precision arguments can be empty strings if they are not needed for the query.
// parameters is a map of the parameter names used in the command to their values.
func NewQueryWithParameters(command, database, precision string, parameters map[string]interface***REMOVED******REMOVED***) Query ***REMOVED***
	return Query***REMOVED***
		Command:    command,
		Database:   database,
		Precision:  precision,
		Parameters: parameters,
	***REMOVED***
***REMOVED***

// Response represents a list of statement results.
type Response struct ***REMOVED***
	Results []Result
	Err     string `json:"error,omitempty"`
***REMOVED***

// Error returns the first error from any statement.
// It returns nil if no errors occurred on any statements.
func (r *Response) Error() error ***REMOVED***
	if r.Err != "" ***REMOVED***
		return errors.New(r.Err)
	***REMOVED***
	for _, result := range r.Results ***REMOVED***
		if result.Err != "" ***REMOVED***
			return errors.New(result.Err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Message represents a user message.
type Message struct ***REMOVED***
	Level string
	Text  string
***REMOVED***

// Result represents a resultset returned from a single statement.
type Result struct ***REMOVED***
	Series   []models.Row
	Messages []*Message
	Err      string `json:"error,omitempty"`
***REMOVED***

// Query sends a command to the server and returns the Response.
func (c *client) Query(q Query) (*Response, error) ***REMOVED***
	req, err := c.createDefaultRequest(q)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	params := req.URL.Query()
	if q.Chunked ***REMOVED***
		params.Set("chunked", "true")
		if q.ChunkSize > 0 ***REMOVED***
			params.Set("chunk_size", strconv.Itoa(q.ChunkSize))
		***REMOVED***
		req.URL.RawQuery = params.Encode()
	***REMOVED***
	resp, err := c.httpClient.Do(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var response Response
	if q.Chunked ***REMOVED***
		cr := NewChunkedResponse(resp.Body)
		for ***REMOVED***
			r, err := cr.NextResponse()
			if err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					break
				***REMOVED***
				// If we got an error while decoding the response, send that back.
				return nil, err
			***REMOVED***

			if r == nil ***REMOVED***
				break
			***REMOVED***

			response.Results = append(response.Results, r.Results...)
			if r.Err != "" ***REMOVED***
				response.Err = r.Err
				break
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		dec := json.NewDecoder(resp.Body)
		dec.UseNumber()
		decErr := dec.Decode(&response)

		// ignore this error if we got an invalid status code
		if decErr != nil && decErr.Error() == "EOF" && resp.StatusCode != http.StatusOK ***REMOVED***
			decErr = nil
		***REMOVED***
		// If we got a valid decode error, send that back
		if decErr != nil ***REMOVED***
			return nil, fmt.Errorf("unable to decode json: received status code %d err: %s", resp.StatusCode, decErr)
		***REMOVED***
	***REMOVED***

	// If we don't have an error in our json response, and didn't get statusOK
	// then send back an error
	if resp.StatusCode != http.StatusOK && response.Error() == nil ***REMOVED***
		return &response, fmt.Errorf("received status code %d from server", resp.StatusCode)
	***REMOVED***
	return &response, nil
***REMOVED***

// QueryAsChunk sends a command to the server and returns the Response.
func (c *client) QueryAsChunk(q Query) (*ChunkedResponse, error) ***REMOVED***
	req, err := c.createDefaultRequest(q)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	params := req.URL.Query()
	params.Set("chunked", "true")
	if q.ChunkSize > 0 ***REMOVED***
		params.Set("chunk_size", strconv.Itoa(q.ChunkSize))
	***REMOVED***
	req.URL.RawQuery = params.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := checkResponse(resp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return NewChunkedResponse(resp.Body), nil
***REMOVED***

func checkResponse(resp *http.Response) error ***REMOVED***
	// If we lack a X-Influxdb-Version header, then we didn't get a response from influxdb
	// but instead some other service. If the error code is also a 500+ code, then some
	// downstream loadbalancer/proxy/etc had an issue and we should report that.
	if resp.Header.Get("X-Influxdb-Version") == "" && resp.StatusCode >= http.StatusInternalServerError ***REMOVED***
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil || len(body) == 0 ***REMOVED***
			return fmt.Errorf("received status code %d from downstream server", resp.StatusCode)
		***REMOVED***

		return fmt.Errorf("received status code %d from downstream server, with response body: %q", resp.StatusCode, body)
	***REMOVED***

	// If we get an unexpected content type, then it is also not from influx direct and therefore
	// we want to know what we received and what status code was returned for debugging purposes.
	if cType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type")); cType != "application/json" ***REMOVED***
		// Read up to 1kb of the body to help identify downstream errors and limit the impact of things
		// like downstream serving a large file
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024))
		if err != nil || len(body) == 0 ***REMOVED***
			return fmt.Errorf("expected json response, got empty body, with status: %v", resp.StatusCode)
		***REMOVED***

		return fmt.Errorf("expected json response, got %q, with status: %v and response body: %q", cType, resp.StatusCode, body)
	***REMOVED***
	return nil
***REMOVED***

func (c *client) createDefaultRequest(q Query) (*http.Request, error) ***REMOVED***
	u := c.url
	u.Path = path.Join(u.Path, "query")

	jsonParameters, err := json.Marshal(q.Parameters)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	req.Header.Set("Content-Type", "")
	req.Header.Set("User-Agent", c.useragent)

	if c.username != "" ***REMOVED***
		req.SetBasicAuth(c.username, c.password)
	***REMOVED***

	params := req.URL.Query()
	params.Set("q", q.Command)
	params.Set("db", q.Database)
	if q.RetentionPolicy != "" ***REMOVED***
		params.Set("rp", q.RetentionPolicy)
	***REMOVED***
	params.Set("params", string(jsonParameters))

	if q.Precision != "" ***REMOVED***
		params.Set("epoch", q.Precision)
	***REMOVED***
	req.URL.RawQuery = params.Encode()

	return req, nil

***REMOVED***

// duplexReader reads responses and writes it to another writer while
// satisfying the reader interface.
type duplexReader struct ***REMOVED***
	r io.ReadCloser
	w io.Writer
***REMOVED***

func (r *duplexReader) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = r.r.Read(p)
	if err == nil ***REMOVED***
		r.w.Write(p[:n])
	***REMOVED***
	return n, err
***REMOVED***

// Close closes the response.
func (r *duplexReader) Close() error ***REMOVED***
	return r.r.Close()
***REMOVED***

// ChunkedResponse represents a response from the server that
// uses chunking to stream the output.
type ChunkedResponse struct ***REMOVED***
	dec    *json.Decoder
	duplex *duplexReader
	buf    bytes.Buffer
***REMOVED***

// NewChunkedResponse reads a stream and produces responses from the stream.
func NewChunkedResponse(r io.Reader) *ChunkedResponse ***REMOVED***
	rc, ok := r.(io.ReadCloser)
	if !ok ***REMOVED***
		rc = ioutil.NopCloser(r)
	***REMOVED***
	resp := &ChunkedResponse***REMOVED******REMOVED***
	resp.duplex = &duplexReader***REMOVED***r: rc, w: &resp.buf***REMOVED***
	resp.dec = json.NewDecoder(resp.duplex)
	resp.dec.UseNumber()
	return resp
***REMOVED***

// NextResponse reads the next line of the stream and returns a response.
func (r *ChunkedResponse) NextResponse() (*Response, error) ***REMOVED***
	var response Response
	if err := r.dec.Decode(&response); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			return nil, err
		***REMOVED***
		// A decoding error happened. This probably means the server crashed
		// and sent a last-ditch error message to us. Ensure we have read the
		// entirety of the connection to get any remaining error text.
		io.Copy(ioutil.Discard, r.duplex)
		return nil, errors.New(strings.TrimSpace(r.buf.String()))
	***REMOVED***

	r.buf.Reset()
	return &response, nil
***REMOVED***

// Close closes the response.
func (r *ChunkedResponse) Close() error ***REMOVED***
	return r.duplex.Close()
***REMOVED***
