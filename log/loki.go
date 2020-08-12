/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type lokiHook struct ***REMOVED***
	addr           string
	labels         [][2]string
	ch             chan *logrus.Entry
	limit          int
	msgMaxSize     int
	levels         []logrus.Level
	pushPeriod     time.Duration
	client         *http.Client
	ctx            context.Context
	fallbackLogger logrus.FieldLogger
	profile        bool
***REMOVED***

// LokiFromConfigLine returns a new logrus.Hook that pushes logrus.Entrys to loki and is configured
// through the provided line
//nolint:funlen
func LokiFromConfigLine(ctx context.Context, fallbackLogger logrus.FieldLogger, line string) (logrus.Hook, error) ***REMOVED***
	h := &lokiHook***REMOVED***
		addr:           "http://127.0.0.1:3100/loki/api/v1/push",
		limit:          100,
		levels:         logrus.AllLevels,
		pushPeriod:     time.Second * 1,
		ctx:            ctx,
		msgMaxSize:     1024 * 1024, // 1mb
		ch:             make(chan *logrus.Entry, 1000),
		fallbackLogger: fallbackLogger,
	***REMOVED***
	if line == "loki" ***REMOVED***
		return h, nil
	***REMOVED***

	parts := strings.SplitN(line, "=", 2)
	if parts[0] != "loki" ***REMOVED***
		return nil, fmt.Errorf("loki configuration should be in the form `loki=url-to-push` but is `%s`", line)
	***REMOVED***
	args := strings.Split(parts[1], ",")
	h.addr = args[0]
	// TODO use something better ... maybe
	// https://godoc.org/github.com/kubernetes/helm/pkg/strvals
	// atleast until https://github.com/loadimpact/k6/issues/926?
	if len(args) == 1 ***REMOVED***
		return h, nil
	***REMOVED***

	for _, arg := range args[1:] ***REMOVED***
		paramParts := strings.SplitN(arg, "=", 2)

		if len(paramParts) != 2 ***REMOVED***
			return nil, fmt.Errorf("loki arguments should be in the form `address,key1=value1,key2=value2`, got %s", arg)
		***REMOVED***

		key, value := paramParts[0], paramParts[1]

		var err error
		switch key ***REMOVED***
		case "pushPeriod":
			h.pushPeriod, err = time.ParseDuration(value)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("couldn't parse the loki pushPeriod %w", err)
			***REMOVED***
		case "profile":
			h.profile = true
		case "limit":
			h.limit, err = strconv.Atoi(value)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("couldn't parse the loki limit as a number %w", err)
			***REMOVED***
			if !(h.limit > 0) ***REMOVED***
				return nil, fmt.Errorf("loki limit needs to be a positive number, is %d", h.limit)
			***REMOVED***
		case "msgMaxSize":
			h.msgMaxSize, err = strconv.Atoi(value)
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("couldn't parse the loki msgMaxSize as a number %w", err)
			***REMOVED***
			if !(h.msgMaxSize > 0) ***REMOVED***
				return nil, fmt.Errorf("loki msgMaxSize needs to be a positive number, is %d", h.msgMaxSize)
			***REMOVED***
		case "level":
			h.levels, err = getLevels(value)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		default:
			if strings.HasPrefix(key, "label.") ***REMOVED***
				labelKey := strings.TrimPrefix(key, "label.")
				h.labels = append(h.labels, [2]string***REMOVED***labelKey, value***REMOVED***)

				continue
			***REMOVED***

			return nil, fmt.Errorf("unknown loki config key %s", key)
		***REMOVED***
	***REMOVED***

	h.client = &http.Client***REMOVED***Timeout: h.pushPeriod***REMOVED***

	go h.loop()

	return h, nil
***REMOVED***

func getLevels(level string) ([]logrus.Level, error) ***REMOVED***
	lvl, err := logrus.ParseLevel(level)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("unknown log level %s", level) // specifically use a custom error
	***REMOVED***
	index := sort.Search(len(logrus.AllLevels), func(i int) bool ***REMOVED***
		return logrus.AllLevels[i] > lvl
	***REMOVED***)

	return logrus.AllLevels[:index], nil
***REMOVED***

// fill one of two equally sized slices with entries and then push it while filling the other one
// TODO benchmark this
//nolint:funlen
func (h *lokiHook) loop() ***REMOVED***
	var (
		msgs       = make([]tmpMsg, h.limit)
		msgsToPush = make([]tmpMsg, h.limit)
		dropped    int
		count      int
		ticker     = time.NewTicker(h.pushPeriod)
		pushCh     = make(chan chan int64)
	)

	defer ticker.Stop()
	defer close(pushCh)

	go func() ***REMOVED***
		oldLogs := make([]tmpMsg, 0, h.limit*2)
		for ch := range pushCh ***REMOVED***
			msgsToPush, msgs = msgs, msgsToPush
			oldCount, oldDropped := count, dropped
			count, dropped = 0, 0
			cutOff := <-ch
			close(ch) // signal that more buffering can continue

			copy(oldLogs[len(oldLogs):len(oldLogs)+oldCount], msgsToPush[:oldCount])
			oldLogs = oldLogs[:len(oldLogs)+oldCount]

			t := time.Now()
			cutOffIndex := sortAndSplitMsgs(oldLogs, cutOff)
			if cutOffIndex == 0 ***REMOVED***
				continue
			***REMOVED***
			t1 := time.Since(t)

			pushMsg := h.createPushMessage(oldLogs, cutOffIndex, oldDropped)
			if cutOffIndex > len(oldLogs) ***REMOVED***
				oldLogs = oldLogs[:0]

				continue
			***REMOVED***
			oldLogs = oldLogs[:copy(oldLogs, oldLogs[cutOffIndex:])]
			t2 := time.Since(t) - t1

			var b bytes.Buffer
			_, err := pushMsg.WriteTo(&b)
			if err != nil ***REMOVED***
				h.fallbackLogger.WithError(err).Error("Error while marshaling logs for loki")

				continue
			***REMOVED***
			size := b.Len()
			t3 := time.Since(t) - t2 - t1

			err = h.push(b)
			if err != nil ***REMOVED***
				h.fallbackLogger.WithError(err).Error("Error while sending logs to loki")

				continue
			***REMOVED***
			t4 := time.Since(t) - t3 - t2 - t1

			if h.profile ***REMOVED***
				h.fallbackLogger.Infof(
					"sorting=%s, adding=%s marshalling=%s sending=%s count=%d final_size=%d\n",
					t1, t2, t3, t4, cutOffIndex, size)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	for ***REMOVED***
		select ***REMOVED***
		case entry := <-h.ch:
			if count == h.limit ***REMOVED***
				dropped++

				continue
			***REMOVED***

			// Arguably we can directly generate the final marshalled version of the labels right here
			// through sorting the entry.Data, removing additionalparams from it and then dumping it
			// as the final marshal and appending level and h.labels after it.
			// If we reuse some kind of big enough `[]byte` buffer we can also possibly skip on some
			// of allocation. Combined with the cutoff part and directly pushing in the final data
			// type this can be really a lot faster and to use a lot less memory
			labels := make(map[string]string, len(entry.Data)+1)
			for k, v := range entry.Data ***REMOVED***
				labels[k] = fmt.Sprint(v) // TODO optimize ?
			***REMOVED***
			for _, params := range h.labels ***REMOVED***
				labels[params[0]] = params[1]
			***REMOVED***
			labels["level"] = entry.Level.String()
			// have the cutoff here ?
			// if we cutoff here we can cut somewhat on the backbuffers and optimize the inserting
			// in/creating of the final Streams that we push
			msgs[count] = tmpMsg***REMOVED***
				labels: labels,
				msg:    entry.Message,
				t:      entry.Time.UnixNano(),
			***REMOVED***
			count++
		case t := <-ticker.C:
			ch := make(chan int64)
			pushCh <- ch
			ch <- t.Add(-(h.pushPeriod / 2)).UnixNano()
			<-ch
		case <-h.ctx.Done():
			ch := make(chan int64)
			pushCh <- ch
			ch <- 0
			<-ch

			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func sortAndSplitMsgs(msgs []tmpMsg, cutOff int64) int ***REMOVED***
	if len(msgs) == 0 ***REMOVED***
		return 0
	***REMOVED***

	// TODO using time.Before was giving a lot of out of order, but even now, there are some, if the
	// limit is big enough ...
	sort.Slice(msgs, func(i, j int) bool ***REMOVED***
		return msgs[i].t < msgs[j].t
	***REMOVED***)

	cutOffIndex := sort.Search(len(msgs), func(i int) bool ***REMOVED***
		return !(msgs[i].t < cutOff)
	***REMOVED***)

	return cutOffIndex
***REMOVED***

func (h *lokiHook) createPushMessage(msgs []tmpMsg, cutOffIndex, dropped int) *lokiPushMessage ***REMOVED***
	pushMsg := new(lokiPushMessage)
	pushMsg.maxSize = h.msgMaxSize
	for _, msg := range msgs[:cutOffIndex] ***REMOVED***
		pushMsg.add(msg)
	***REMOVED***
	if dropped != 0 ***REMOVED***
		labels := make(map[string]string, 2+len(h.labels))
		labels["level"] = logrus.WarnLevel.String()
		labels["dropped"] = strconv.Itoa(dropped)
		for _, params := range h.labels ***REMOVED***
			labels[params[0]] = params[1]
		***REMOVED***

		msg := tmpMsg***REMOVED***
			labels: labels,
			msg: fmt.Sprintf("k6 dropped some log messages because they were above the limit of %d/%s",
				h.limit, h.pushPeriod),
			t: msgs[cutOffIndex-1].t,
		***REMOVED***
		pushMsg.add(msg)
	***REMOVED***

	return pushMsg
***REMOVED***

func (h *lokiHook) push(b bytes.Buffer) error ***REMOVED***
	body := b.Bytes()

	req, err := http.NewRequestWithContext(context.Background(), "GET", h.addr, &b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	req.GetBody = func() (io.ReadCloser, error) ***REMOVED***
		return ioutil.NopCloser(bytes.NewBuffer(body)), nil
	***REMOVED***

	req.Header.Set("Content-Type", "application/json")

	res, err := h.client.Do(req)

	if res != nil ***REMOVED***
		if res.StatusCode >= 400 ***REMOVED***
			r, _ := ioutil.ReadAll(res.Body) // maybe limit it to something like the first 1000 characters?

			return fmt.Errorf("got %d from loki: %s", res.StatusCode, string(r))
		***REMOVED***
		_, _ = io.Copy(ioutil.Discard, res.Body)
		_ = res.Body.Close()
	***REMOVED***

	return err
***REMOVED***

func mapEqual(a, b map[string]string) bool ***REMOVED***
	if len(a) != len(b) ***REMOVED***
		return false
	***REMOVED***
	for k, v := range a ***REMOVED***
		if v2, ok := b[k]; !ok || v2 != v ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func (pushMsg *lokiPushMessage) add(entry tmpMsg) ***REMOVED***
	var foundStrm *stream
	for _, strm := range pushMsg.Streams ***REMOVED***
		if mapEqual(strm.Stream, entry.labels) ***REMOVED***
			foundStrm = strm

			break
		***REMOVED***
	***REMOVED***

	if foundStrm == nil ***REMOVED***
		foundStrm = &stream***REMOVED***Stream: entry.labels***REMOVED***
		pushMsg.Streams = append(pushMsg.Streams, foundStrm)
	***REMOVED***

	foundStrm.Values = append(foundStrm.Values, logEntry***REMOVED***t: entry.t, msg: entry.msg***REMOVED***)
***REMOVED***

// this is temporary message format used to not keep the logrus.Entry around too long and to make
// sorting easier
type tmpMsg struct ***REMOVED***
	labels map[string]string
	t      int64
	msg    string
***REMOVED***

func (h *lokiHook) Fire(entry *logrus.Entry) error ***REMOVED***
	h.ch <- entry

	return nil
***REMOVED***

func (h *lokiHook) Levels() []logrus.Level ***REMOVED***
	return h.levels
***REMOVED***

/*
***REMOVED***
  "streams": [
    ***REMOVED***
      "stream": ***REMOVED***
        "label1": "value1"
        "label2": "value2"
      ***REMOVED***,
      "values": [ // the nanoseconds need to be in order
          [ "<unix epoch in nanoseconds>", "<log line>" ],
          [ "<unix epoch in nanoseconds>", "<log line>" ]
      ]
    ***REMOVED***
  ]
***REMOVED***
*/
type lokiPushMessage struct ***REMOVED***
	Streams []*stream `json:"streams"`
	maxSize int
***REMOVED***

func (pushMsg *lokiPushMessage) WriteTo(w io.Writer) (n int64, err error) ***REMOVED***
	var k int
	write := func(b []byte) ***REMOVED***
		if err != nil ***REMOVED***
			return
		***REMOVED***
		k, err = w.Write(b)
		n += int64(k)
	***REMOVED***
	// 10+ 9 for the amount of nanoseconds between 2001 and 2286 also it overflows in the year 2262 ;)
	var nanoseconds [19]byte
	write([]byte(`***REMOVED***"streams":[`))
	var b []byte
	for i, str := range pushMsg.Streams ***REMOVED***
		if i != 0 ***REMOVED***
			write([]byte(`,`))
		***REMOVED***
		write([]byte(`***REMOVED***"stream":***REMOVED***`))
		var f bool
		for k, v := range str.Stream ***REMOVED***
			if f ***REMOVED***
				write([]byte(`,`))
			***REMOVED***
			f = true
			write([]byte(`"`))
			write([]byte(k))
			write([]byte(`":`))
			b, err = json.Marshal(v)
			if err != nil ***REMOVED***
				return n, err
			***REMOVED***
			write(b)
		***REMOVED***
		write([]byte(`***REMOVED***,"values":[`))
		for j, v := range str.Values ***REMOVED***
			if j != 0 ***REMOVED***
				write([]byte(`,`))
			***REMOVED***
			write([]byte(`["`))
			strconv.AppendInt(nanoseconds[:0], v.t, 10)
			write(nanoseconds[:])
			write([]byte(`",`))
			if len([]rune(v.msg)) > pushMsg.maxSize ***REMOVED***
				difference := int64(len(v.msg) - pushMsg.maxSize)
				omitMsg := append(strconv.AppendInt([]byte("... omitting "), difference, 10), " characters ..."...)
				v.msg = strings.Join([]string***REMOVED***
					string([]rune(v.msg)[:pushMsg.maxSize/2]),
					string([]rune(v.msg)[len([]rune(v.msg))-pushMsg.maxSize/2:]),
				***REMOVED***, string(omitMsg))
			***REMOVED***

			b, err = json.Marshal(v.msg)
			if err != nil ***REMOVED***
				return n, err
			***REMOVED***
			write(b)
			write([]byte(`]`))
		***REMOVED***
		write([]byte(`]***REMOVED***`))
	***REMOVED***

	write([]byte(`]***REMOVED***`))

	return n, err
***REMOVED***

type stream struct ***REMOVED***
	Stream map[string]string `json:"stream"`
	Values []logEntry        `json:"values"`
***REMOVED***

type logEntry struct ***REMOVED***
	t   int64  // nanoseconds
	msg string // maybe intern those as they are likely to be the same for an interval
***REMOVED***

// rewrite this either with easyjson or with a custom marshalling
func (l logEntry) MarshalJSON() ([]byte, error) ***REMOVED***
	// 2 for '[]', 1 for ',', 4 for '"' and 10 + 9 for the amount of nanoseconds between 2001 and
	// 2286 also it overflows in the year 2262 ;)
	b := make([]byte, 2, len(l.msg)+26)
	b[0] = '['
	b[1] = '"'
	b = strconv.AppendInt(b, l.t, 10)
	b = append(b, '"', ',', '"')
	b = append(b, l.msg...)
	b = append(b, '"', ']')

	return b, nil
***REMOVED***
