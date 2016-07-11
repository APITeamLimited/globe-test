package influxdb

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/speedboat/stats"
)

type Backend struct ***REMOVED***
	Client   client.Client
	Database string
***REMOVED***

func New(conf client.HTTPConfig, db string) (*Backend, error) ***REMOVED***
	c, err := client.NewHTTPClient(conf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Backend***REMOVED***
		Client:   c,
		Database: db,
	***REMOVED***, nil
***REMOVED***

func NewFromURL(url string) (*Backend, error) ***REMOVED***
	conf, db, err := parseURL(url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return New(conf, db)
***REMOVED***

func (b *Backend) Submit(batches [][]stats.Sample) error ***REMOVED***
	pb, err := client.NewBatchPoints(client.BatchPointsConfig***REMOVED***
		Database: b.Database,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, batch := range batches ***REMOVED***
		for _, s := range batch ***REMOVED***
			pt, err := makeInfluxPoint(s)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			pb.AddPoint(pt)
		***REMOVED***
	***REMOVED***

	if err := b.Client.Write(pb); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
