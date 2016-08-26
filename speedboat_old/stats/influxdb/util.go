package influxdb

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/loadimpact/speedboat/stats"
	neturl "net/url"
)

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

func makeInfluxPoint(s stats.Sample) (*client.Point, error) ***REMOVED***
	tags := make(map[string]string)
	for key, val := range s.Tags ***REMOVED***
		tags[key] = fmt.Sprint(val)
	***REMOVED***
	fields := make(map[string]interface***REMOVED******REMOVED***)
	for key, val := range s.Values ***REMOVED***
		fields[key] = val
	***REMOVED***
	return client.NewPoint(s.Stat.Name, tags, fields, s.Time)
***REMOVED***
