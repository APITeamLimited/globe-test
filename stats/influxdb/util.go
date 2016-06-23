package influxdb

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
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
