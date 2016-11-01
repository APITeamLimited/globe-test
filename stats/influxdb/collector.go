package influxdb

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/speedboat/stats"
	"net/url"
	"sync"
	"time"
)

const pushInterval = 1 * time.Second

type Collector struct ***REMOVED***
	client    client.Client
	batchConf client.BatchPointsConfig

	buffers      []*Buffer
	buffersMutex sync.Mutex
***REMOVED***

func New(u *url.URL) (*Collector, error) ***REMOVED***
	cl, batchConf, err := parseURL(u)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Collector***REMOVED***
		client:    cl,
		batchConf: batchConf,
	***REMOVED***, nil
***REMOVED***

func (c *Collector) Run(ctx context.Context) ***REMOVED***
	log.Debug("InfluxDB: Running!")
	ticker := time.NewTicker(pushInterval)
	for ***REMOVED***
		select ***REMOVED***
		case <-ticker.C:
			c.commit()
		case <-ctx.Done():
			c.commit()
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *Collector) Buffer() stats.Buffer ***REMOVED***
	buf := &(Buffer***REMOVED******REMOVED***)
	c.buffersMutex.Lock()
	c.buffers = append(c.buffers, buf)
	c.buffersMutex.Unlock()
	return buf
***REMOVED***

func (c *Collector) commit() ***REMOVED***
	log.Debug("InfluxDB: Committing...")
	batch, err := client.NewBatchPoints(c.batchConf)
	if err != nil ***REMOVED***
		log.WithError(err).Error("InfluxDB: Couldn't make a batch")
		return
	***REMOVED***

	buffers := c.buffers
	samples := []stats.Sample***REMOVED******REMOVED***
	for _, buf := range buffers ***REMOVED***
		samples = append(samples, buf.Drain()...)
	***REMOVED***

	for _, sample := range samples ***REMOVED***
		p, err := client.NewPoint(
			sample.Metric.Name,
			sample.Tags,
			map[string]interface***REMOVED******REMOVED******REMOVED***"value": sample.Value***REMOVED***,
			sample.Time,
		)
		if err != nil ***REMOVED***
			log.WithError(err).Error("InfluxDB: Couldn't make point from sample!")
			return
		***REMOVED***
		batch.AddPoint(p)
	***REMOVED***

	log.WithField("points", len(batch.Points())).Debug("InfluxDB: Writing points...")
	if err := c.client.Write(batch); err != nil ***REMOVED***
		log.WithError(err).Error("InfluxDB: Couldn't write stats")
	***REMOVED***
***REMOVED***

type Buffer []stats.Sample

func (b *Buffer) Add(samples ...stats.Sample) ***REMOVED***
	*b = append(*b, samples...)
***REMOVED***

func (b *Buffer) Drain() []stats.Sample ***REMOVED***
	old := *b
	*b = (*b)[:0]
	return old
***REMOVED***
