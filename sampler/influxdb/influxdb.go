package influxdb

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/speedboat/sampler"
	neturl "net/url"
	"sync"
)

type Output struct ***REMOVED***
	Client   client.Client
	Database string

	currentBatch client.BatchPoints
	batchMutex   sync.Mutex
***REMOVED***

func New(conf client.HTTPConfig, db string) (*Output, error) ***REMOVED***
	c, err := client.NewHTTPClient(conf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Output***REMOVED***
		Client:   c,
		Database: db,
	***REMOVED***, nil
***REMOVED***

func NewFromURL(url string) (*Output, error) ***REMOVED***
	conf, db, err := parseURL(url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return New(conf, db)
***REMOVED***

func parseURL(url string) (conf client.HTTPConfig, db string, err error) ***REMOVED***
	u, err := neturl.Parse(url)
	if err != nil ***REMOVED***
		return conf, db, err
	***REMOVED***

	if u.Path == "" || u.Path == "/" ***REMOVED***
		return conf, db, errors.New("No InfluxDB database specified")
	***REMOVED***
	db = u.Path[1:]

	conf.Addr = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if u.User != nil ***REMOVED***
		conf.Username = u.User.Username()
		conf.Password, _ = u.User.Password()
	***REMOVED***
	return conf, db, nil
***REMOVED***

func (o *Output) Write(m *sampler.Metric, e *sampler.Entry) error ***REMOVED***
	o.batchMutex.Lock()
	defer o.batchMutex.Unlock()

	if o.currentBatch == nil ***REMOVED***
		batch, err := client.NewBatchPoints(client.BatchPointsConfig***REMOVED***
			Database: o.Database,
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		o.currentBatch = batch
	***REMOVED***

	tags := make(map[string]string)
	for key, value := range e.Fields ***REMOVED***
		tags[key] = fmt.Sprint(value)
	***REMOVED***
	fields := map[string]interface***REMOVED******REMOVED******REMOVED***"value": e.Value***REMOVED***

	point, err := client.NewPoint(m.Name, tags, fields, e.Time)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.currentBatch.AddPoint(point)
	return nil
***REMOVED***

func (o *Output) Commit() error ***REMOVED***
	o.batchMutex.Lock()
	defer o.batchMutex.Unlock()

	if o.currentBatch == nil ***REMOVED***
		return nil
	***REMOVED***

	err := o.Client.Write(o.currentBatch)
	o.currentBatch = nil
	return err
***REMOVED***
