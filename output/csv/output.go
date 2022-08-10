package csv

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

// Output implements the lib.Output interface for saving to CSV files.
type Output struct ***REMOVED***
	output.SampleBuffer

	params          output.Params
	periodicFlusher *output.PeriodicFlusher

	logger    logrus.FieldLogger
	fname     string
	csvWriter *csv.Writer
	csvLock   sync.Mutex
	closeFn   func() error

	resTags      []string
	ignoredTags  []string
	row          []string
	saveInterval time.Duration
	timeFormat   TimeFormat
***REMOVED***

// New Creates new instance of CSV output
func New(params output.Params) (output.Output, error) ***REMOVED***
	return newOutput(params)
***REMOVED***

func newOutput(params output.Params) (*Output, error) ***REMOVED***
	resTags := []string***REMOVED******REMOVED***
	ignoredTags := []string***REMOVED******REMOVED***
	tags := params.ScriptOptions.SystemTags.Map()
	for tag, flag := range tags ***REMOVED***
		if flag ***REMOVED***
			resTags = append(resTags, tag)
		***REMOVED*** else ***REMOVED***
			ignoredTags = append(ignoredTags, tag)
		***REMOVED***
	***REMOVED***

	sort.Strings(resTags)
	sort.Strings(ignoredTags)

	logger := params.Logger.WithFields(logrus.Fields***REMOVED***
		"output":   "csv",
		"filename": params.ConfigArgument,
	***REMOVED***)
	config, err := GetConsolidatedConfig(params.JSONConfig, params.Environment, params.ConfigArgument, logger.Logger)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	timeFormat, err := TimeFormatString(config.TimeFormat.String)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	saveInterval := config.SaveInterval.TimeDuration()
	fname := config.FileName.String

	if fname == "" || fname == "-" ***REMOVED***
		stdoutWriter := csv.NewWriter(os.Stdout)
		return &Output***REMOVED***
			fname:        "-",
			resTags:      resTags,
			ignoredTags:  ignoredTags,
			csvWriter:    stdoutWriter,
			row:          make([]string, 3+len(resTags)+1),
			saveInterval: saveInterval,
			timeFormat:   timeFormat,
			closeFn:      func() error ***REMOVED*** return nil ***REMOVED***,
			logger:       logger,
			params:       params,
		***REMOVED***, nil
	***REMOVED***

	logFile, err := params.FS.Create(fname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c := Output***REMOVED***
		fname:        fname,
		resTags:      resTags,
		ignoredTags:  ignoredTags,
		row:          make([]string, 3+len(resTags)+1),
		saveInterval: saveInterval,
		timeFormat:   timeFormat,
		logger:       logger,
		params:       params,
	***REMOVED***

	if strings.HasSuffix(fname, ".gz") ***REMOVED***
		outfile := gzip.NewWriter(logFile)
		csvWriter := csv.NewWriter(outfile)
		c.csvWriter = csvWriter
		c.closeFn = func() error ***REMOVED***
			_ = outfile.Close()
			return logFile.Close()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		csvWriter := csv.NewWriter(logFile)
		c.csvWriter = csvWriter
		c.closeFn = logFile.Close
	***REMOVED***

	return &c, nil
***REMOVED***

// Description returns a human-readable description of the output.
func (o *Output) Description() string ***REMOVED***
	if o.fname == "" || o.fname == "-" ***REMOVED*** // TODO rename
		return "csv (stdout)"
	***REMOVED***
	return fmt.Sprintf("csv (%s)", o.fname)
***REMOVED***

// Start writes the csv header and starts a new output.PeriodicFlusher
func (o *Output) Start() error ***REMOVED***
	o.logger.Debug("Starting...")

	header := MakeHeader(o.resTags)
	err := o.csvWriter.Write(header)
	if err != nil ***REMOVED***
		o.logger.WithField("filename", o.fname).Error("CSV: Error writing column names to file")
	***REMOVED***
	o.csvWriter.Flush()

	pf, err := output.NewPeriodicFlusher(o.saveInterval, o.flushMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	o.logger.Debug("Started!")
	o.periodicFlusher = pf

	return nil
***REMOVED***

// Stop flushes any remaining metrics and stops the goroutine.
func (o *Output) Stop() error ***REMOVED***
	o.logger.Debug("Stopping...")
	defer o.logger.Debug("Stopped!")
	o.periodicFlusher.Stop()
	return o.closeFn()
***REMOVED***

// flushMetrics Writes samples to the csv file
func (o *Output) flushMetrics() ***REMOVED***
	samples := o.GetBufferedSamples()

	if len(samples) > 0 ***REMOVED***
		o.csvLock.Lock()
		defer o.csvLock.Unlock()
		for _, sc := range samples ***REMOVED***
			for _, sample := range sc.GetSamples() ***REMOVED***
				sample := sample
				row := SampleToRow(&sample, o.resTags, o.ignoredTags, o.row, o.timeFormat)
				err := o.csvWriter.Write(row)
				if err != nil ***REMOVED***
					o.logger.WithField("filename", o.fname).Error("CSV: Error writing to file")
				***REMOVED***
			***REMOVED***
		***REMOVED***
		o.csvWriter.Flush()
	***REMOVED***
***REMOVED***

// MakeHeader creates list of column names for csv file
func MakeHeader(tags []string) []string ***REMOVED***
	tags = append(tags, "extra_tags")
	return append([]string***REMOVED***"metric_name", "timestamp", "metric_value"***REMOVED***, tags...)
***REMOVED***

// SampleToRow converts sample into array of strings
func SampleToRow(sample *metrics.Sample, resTags []string, ignoredTags []string, row []string,
	timeFormat TimeFormat,
) []string ***REMOVED***
	row[0] = sample.Metric.Name

	switch timeFormat ***REMOVED***
	case TimeFormatRFC3339:
		row[1] = sample.Time.Format(time.RFC3339)
	case TimeFormatUnix:
		row[1] = strconv.FormatInt(sample.Time.Unix(), 10)
	***REMOVED***

	row[2] = fmt.Sprintf("%f", sample.Value)
	sampleTags := sample.Tags.CloneTags()

	for ind, tag := range resTags ***REMOVED***
		row[ind+3] = sampleTags[tag]
	***REMOVED***

	extraTags := bytes.Buffer***REMOVED******REMOVED***
	prev := false
	for tag, val := range sampleTags ***REMOVED***
		if !IsStringInSlice(resTags, tag) && !IsStringInSlice(ignoredTags, tag) ***REMOVED***
			if prev ***REMOVED***
				if _, err := extraTags.WriteString("&"); err != nil ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			if _, err := extraTags.WriteString(tag); err != nil ***REMOVED***
				break
			***REMOVED***

			if _, err := extraTags.WriteString("="); err != nil ***REMOVED***
				break
			***REMOVED***

			if _, err := extraTags.WriteString(val); err != nil ***REMOVED***
				break
			***REMOVED***
			prev = true
		***REMOVED***
	***REMOVED***
	row[len(row)-1] = extraTags.String()

	return row
***REMOVED***

// IsStringInSlice returns whether the string is contained within a string slice
func IsStringInSlice(slice []string, str string) bool ***REMOVED***
	if index := sort.SearchStrings(slice, str); index == len(slice) || slice[index] != str ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***
