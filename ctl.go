/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/loadimpact/k6/api/v1"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
)

var commandStatus = cli.Command***REMOVED***
	Name:      "status",
	Usage:     "Looks up the status of a running test",
	ArgsUsage: " ",
	Action:    actionStatus,
	Description: `Status will print the status of a running test to stdout in YAML format.

   Use the global --address/-a flag to specify the host to connect to; the
   default is port 6565 on the local machine.

   Endpoint: /v1/status`,
***REMOVED***

var commandStats = cli.Command***REMOVED***
	Name:      "stats",
	Usage:     "Prints stats for a running test",
	ArgsUsage: " ",
	Action:    actionStats,
	Description: `Stats will print metrics about a running test to stdout in YAML format.

   The result is a dictionary of metrics, in no particular order.

   Endpoint: /v1/metrics`,
***REMOVED***

var commandScale = cli.Command***REMOVED***
	Name:      "scale",
	Usage:     "Scales a running test",
	ArgsUsage: "vus",
	Flags: []cli.Flag***REMOVED***
		cli.Int64Flag***REMOVED***
			Name:  "vus, u",
			Usage: "update the number of running VUs",
		***REMOVED***,
		cli.Int64Flag***REMOVED***
			Name:  "max, m",
			Usage: "update the max number of VUs allowed",
		***REMOVED***,
	***REMOVED***,
	Action: actionScale,
	Description: `Scale will change the number of active VUs of a running test.

   It is an error to scale a test beyond vus-max; this is because instantiating
   new VUs is a very expensive operation, which may skew test results if done
   during a running test. To raise vus-max, use --max/-m.

   Endpoint: /v1/status`,
***REMOVED***

var commandPause = cli.Command***REMOVED***
	Name:      "pause",
	Usage:     "Pauses a running test",
	ArgsUsage: " ",
	Action:    actionPause,
	Description: `Pause pauses a running test.

   Running VUs will finish their current iterations, then suspend themselves
   until woken by the test's resumption. A sleeping VU will consume no CPU
   cycles, but will still occupy memory.

   Endpoint: /v1/status`,
***REMOVED***

var commandResume = cli.Command***REMOVED***
	Name:      "resume",
	Usage:     "Resumes a paused test",
	ArgsUsage: " ",
	Action:    actionResume,
	Description: `Resume resumes a paused test.

   This is the opposite of the pause command, and will do nothing to an already
   running test.

   Endpoint: /v1/status`,
***REMOVED***

func endpointURL(cc *cli.Context, endpoint string) string ***REMOVED***
	return fmt.Sprintf("http://%s%s", cc.GlobalString("address"), endpoint)
***REMOVED***

func apiCall(cc *cli.Context, method, endpoint string, body []byte, dst interface***REMOVED******REMOVED***) error ***REMOVED***
	var bodyReader io.Reader
	if len(body) > 0 ***REMOVED***
		bodyReader = bytes.NewReader(body)
	***REMOVED***

	req, err := http.NewRequest(method, endpointURL(cc, endpoint), bodyReader)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	res, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED*** _ = res.Body.Close() ***REMOVED***()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if res.StatusCode >= 400 ***REMOVED***
		var envelope v1.ErrorResponse
		if err := json.Unmarshal(data, &envelope); err != nil ***REMOVED***
			return err
		***REMOVED***
		return envelope.Errors[0]
	***REMOVED***

	return jsonapi.Unmarshal(data, dst)
***REMOVED***

func actionStatus(cc *cli.Context) error ***REMOVED***
	var status v1.Status
	if err := apiCall(cc, "GET", "/v1/status", nil, &status); err != nil ***REMOVED***
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionStats(cc *cli.Context) error ***REMOVED***
	var metrics []v1.Metric
	if err := apiCall(cc, "GET", "/v1/metrics", nil, &metrics); err != nil ***REMOVED***
		return err
	***REMOVED***
	output := make(map[string]v1.Metric)
	for _, m := range metrics ***REMOVED***
		output[m.GetID()] = m
	***REMOVED***
	return dumpYAML(output)
***REMOVED***

func actionScale(cc *cli.Context) error ***REMOVED***
	patch := v1.Status***REMOVED***
		VUs:    cliInt64(cc, "vus"),
		VUsMax: cliInt64(cc, "max"),
	***REMOVED***
	if !patch.VUs.Valid && !patch.VUsMax.Valid ***REMOVED***
		log.Warn("Neither --vus/-u or --max/-m passed; doing doing nothing")
		return nil
	***REMOVED***

	body, err := jsonapi.Marshal(patch)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Serialization error")
		return err
	***REMOVED***

	var status v1.Status
	if err := apiCall(cc, "PATCH", "/v1/status", body, &status); err != nil ***REMOVED***
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionPause(cc *cli.Context) error ***REMOVED***
	body, err := jsonapi.Marshal(v1.Status***REMOVED***
		Paused: null.BoolFrom(true),
	***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Serialization error")
		return err
	***REMOVED***

	var status v1.Status
	if err := apiCall(cc, "PATCH", "/v1/status", body, &status); err != nil ***REMOVED***
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionResume(cc *cli.Context) error ***REMOVED***
	body, err := jsonapi.Marshal(v1.Status***REMOVED***
		Paused: null.BoolFrom(false),
	***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Serialization error")
		return err
	***REMOVED***

	var status v1.Status
	if err := apiCall(cc, "PATCH", "/v1/status", body, &status); err != nil ***REMOVED***
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***
