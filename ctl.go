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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/k6/api/v1"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"net/http"
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

var commandStart = cli.Command***REMOVED***
	Name:      "start",
	Usage:     "Starts a paused test",
	ArgsUsage: " ",
	Action:    actionStart,
	Description: `Start starts a paused test.

   This is the opposite of the pause command, and will do nothing to an already
   running test.

   Endpoint: /v1/status`,
***REMOVED***

func endpointURL(cc *cli.Context, endpoint string) string ***REMOVED***
	return fmt.Sprintf("http://%s%s", cc.GlobalString("address"), endpoint)
***REMOVED***

func actionStatus(cc *cli.Context) error ***REMOVED***
	res, err := http.Get(endpointURL(cc, "/v1/status"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
		return err
	***REMOVED***
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't read response")
		return err
	***REMOVED***
	var status v1.Status
	if err := jsonapi.Unmarshal(data, &status); err != nil ***REMOVED***
		log.WithError(err).Error("Invalid response")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionStats(cc *cli.Context) error ***REMOVED***
	res, err := http.Get(endpointURL(cc, "/v1/metrics"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
		return err
	***REMOVED***
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't read response")
		return err
	***REMOVED***
	var metrics []v1.Metric
	if err := jsonapi.Unmarshal(data, &metrics); err != nil ***REMOVED***
		log.WithError(err).Error("Invalid response")
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

	req, err := http.NewRequest(
		http.MethodPatch,
		endpointURL(cc, "/v1/status"),
		bytes.NewReader(body),
	)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create request")
		return err
	***REMOVED***

	res, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Request error")
		return err
	***REMOVED***
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't read response")
		return err
	***REMOVED***
	var status v1.Status
	if err := jsonapi.Unmarshal(data, &status); err != nil ***REMOVED***
		log.WithError(err).Error("Invalid response")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionPause(cc *cli.Context) error ***REMOVED***
	return nil
***REMOVED***

func actionStart(cc *cli.Context) error ***REMOVED***
	return nil
***REMOVED***
