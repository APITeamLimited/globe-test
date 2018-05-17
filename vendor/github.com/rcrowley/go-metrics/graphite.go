package metrics

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// GraphiteConfig provides a container with configuration parameters for
// the Graphite exporter
type GraphiteConfig struct ***REMOVED***
	Addr          *net.TCPAddr  // Network address to connect to
	Registry      Registry      // Registry to be exported
	FlushInterval time.Duration // Flush interval
	DurationUnit  time.Duration // Time conversion unit for durations
	Prefix        string        // Prefix to be prepended to metric names
	Percentiles   []float64     // Percentiles to export from timers and histograms
***REMOVED***

// Graphite is a blocking exporter function which reports metrics in r
// to a graphite server located at addr, flushing them every d duration
// and prepending metric names with prefix.
func Graphite(r Registry, d time.Duration, prefix string, addr *net.TCPAddr) ***REMOVED***
	GraphiteWithConfig(GraphiteConfig***REMOVED***
		Addr:          addr,
		Registry:      r,
		FlushInterval: d,
		DurationUnit:  time.Nanosecond,
		Prefix:        prefix,
		Percentiles:   []float64***REMOVED***0.5, 0.75, 0.95, 0.99, 0.999***REMOVED***,
	***REMOVED***)
***REMOVED***

// GraphiteWithConfig is a blocking exporter function just like Graphite,
// but it takes a GraphiteConfig instead.
func GraphiteWithConfig(c GraphiteConfig) ***REMOVED***
	log.Printf("WARNING: This go-metrics client has been DEPRECATED! It has been moved to https://github.com/cyberdelia/go-metrics-graphite and will be removed from rcrowley/go-metrics on August 12th 2015")
	for _ = range time.Tick(c.FlushInterval) ***REMOVED***
		if err := graphite(&c); nil != err ***REMOVED***
			log.Println(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// GraphiteOnce performs a single submission to Graphite, returning a
// non-nil error on failed connections. This can be used in a loop
// similar to GraphiteWithConfig for custom error handling.
func GraphiteOnce(c GraphiteConfig) error ***REMOVED***
	log.Printf("WARNING: This go-metrics client has been DEPRECATED! It has been moved to https://github.com/cyberdelia/go-metrics-graphite and will be removed from rcrowley/go-metrics on August 12th 2015")
	return graphite(&c)
***REMOVED***

func graphite(c *GraphiteConfig) error ***REMOVED***
	now := time.Now().Unix()
	du := float64(c.DurationUnit)
	conn, err := net.DialTCP("tcp", nil, c.Addr)
	if nil != err ***REMOVED***
		return err
	***REMOVED***
	defer conn.Close()
	w := bufio.NewWriter(conn)
	c.Registry.Each(func(name string, i interface***REMOVED******REMOVED***) ***REMOVED***
		switch metric := i.(type) ***REMOVED***
		case Counter:
			fmt.Fprintf(w, "%s.%s.count %d %d\n", c.Prefix, name, metric.Count(), now)
		case Gauge:
			fmt.Fprintf(w, "%s.%s.value %d %d\n", c.Prefix, name, metric.Value(), now)
		case GaugeFloat64:
			fmt.Fprintf(w, "%s.%s.value %f %d\n", c.Prefix, name, metric.Value(), now)
		case Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles(c.Percentiles)
			fmt.Fprintf(w, "%s.%s.count %d %d\n", c.Prefix, name, h.Count(), now)
			fmt.Fprintf(w, "%s.%s.min %d %d\n", c.Prefix, name, h.Min(), now)
			fmt.Fprintf(w, "%s.%s.max %d %d\n", c.Prefix, name, h.Max(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", c.Prefix, name, h.Mean(), now)
			fmt.Fprintf(w, "%s.%s.std-dev %.2f %d\n", c.Prefix, name, h.StdDev(), now)
			for psIdx, psKey := range c.Percentiles ***REMOVED***
				key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
				fmt.Fprintf(w, "%s.%s.%s-percentile %.2f %d\n", c.Prefix, name, key, ps[psIdx], now)
			***REMOVED***
		case Meter:
			m := metric.Snapshot()
			fmt.Fprintf(w, "%s.%s.count %d %d\n", c.Prefix, name, m.Count(), now)
			fmt.Fprintf(w, "%s.%s.one-minute %.2f %d\n", c.Prefix, name, m.Rate1(), now)
			fmt.Fprintf(w, "%s.%s.five-minute %.2f %d\n", c.Prefix, name, m.Rate5(), now)
			fmt.Fprintf(w, "%s.%s.fifteen-minute %.2f %d\n", c.Prefix, name, m.Rate15(), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", c.Prefix, name, m.RateMean(), now)
		case Timer:
			t := metric.Snapshot()
			ps := t.Percentiles(c.Percentiles)
			fmt.Fprintf(w, "%s.%s.count %d %d\n", c.Prefix, name, t.Count(), now)
			fmt.Fprintf(w, "%s.%s.min %d %d\n", c.Prefix, name, t.Min()/int64(du), now)
			fmt.Fprintf(w, "%s.%s.max %d %d\n", c.Prefix, name, t.Max()/int64(du), now)
			fmt.Fprintf(w, "%s.%s.mean %.2f %d\n", c.Prefix, name, t.Mean()/du, now)
			fmt.Fprintf(w, "%s.%s.std-dev %.2f %d\n", c.Prefix, name, t.StdDev()/du, now)
			for psIdx, psKey := range c.Percentiles ***REMOVED***
				key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
				fmt.Fprintf(w, "%s.%s.%s-percentile %.2f %d\n", c.Prefix, name, key, ps[psIdx], now)
			***REMOVED***
			fmt.Fprintf(w, "%s.%s.one-minute %.2f %d\n", c.Prefix, name, t.Rate1(), now)
			fmt.Fprintf(w, "%s.%s.five-minute %.2f %d\n", c.Prefix, name, t.Rate5(), now)
			fmt.Fprintf(w, "%s.%s.fifteen-minute %.2f %d\n", c.Prefix, name, t.Rate15(), now)
			fmt.Fprintf(w, "%s.%s.mean-rate %.2f %d\n", c.Prefix, name, t.RateMean(), now)
		***REMOVED***
		w.Flush()
	***REMOVED***)
	return nil
***REMOVED***
