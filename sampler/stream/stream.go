package stream

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/loadimpact/speedboat/sampler"
	"io"
	"strconv"
	"sync"
)

var (
	ErrUnknownFormat = errors.New("unknown output format")
)

func New(format string, output io.WriteCloser) (sampler.Output, error) ***REMOVED***
	switch format ***REMOVED***
	case "json":
		return &JSONOutput***REMOVED***Output: output***REMOVED***, nil
	case "csv":
		return &CSVOutput***REMOVED***Output: output***REMOVED***, nil
	default:
		return nil, ErrUnknownFormat
	***REMOVED***
***REMOVED***

type JSONOutput struct ***REMOVED***
	Output io.WriteCloser

	encoder *json.Encoder
	mutex   sync.Mutex
***REMOVED***

func (o *JSONOutput) Write(m *sampler.Metric, e *sampler.Entry) error ***REMOVED***
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.encoder == nil ***REMOVED***
		o.encoder = json.NewEncoder(o.Output)
	***REMOVED***
	if err := o.encoder.Encode(e); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (o *JSONOutput) Commit() error ***REMOVED***
	return nil
***REMOVED***

func (o *JSONOutput) Close() error ***REMOVED***
	o.mutex.Lock()
	defer o.mutex.Unlock()

	return o.Output.Close()
***REMOVED***

type CSVOutput struct ***REMOVED***
	Output io.WriteCloser

	writer *csv.Writer
	mutex  sync.Mutex
***REMOVED***

func (o *CSVOutput) Write(m *sampler.Metric, e *sampler.Entry) error ***REMOVED***
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if o.writer == nil ***REMOVED***
		o.writer = csv.NewWriter(o.Output)
	***REMOVED***

	record := []string***REMOVED***m.Name, strconv.FormatInt(e.Value, 10)***REMOVED***
	if err := o.writer.Write(record); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (o *CSVOutput) Commit() error ***REMOVED***
	o.writer.Flush()
	return nil
***REMOVED***

func (o *CSVOutput) Close() error ***REMOVED***
	o.mutex.Lock()
	defer o.mutex.Unlock()

	return o.Output.Close()
***REMOVED***
