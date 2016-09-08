package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/loadimpact/speedboat/api"
	"github.com/loadimpact/speedboat/lib"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strconv"
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

var commandScale = cli.Command***REMOVED***
	Name:      "scale",
	Usage:     "Scales a running test",
	ArgsUsage: "vus",
	Action:    actionScale,
	Description: `Scale will change the number of active VUs of a running test.

   It is an error to scale a test beyond vus-max; this is because instantiating
   new VUs is a very expensive operation, which may skew test results if done
   during a running test, and should thus be done deliberately.

   Endpoint: /v1/status`,
***REMOVED***

var commandCap = cli.Command***REMOVED***
	Name:      "cap",
	Usage:     "Changes the VU cap for a running test",
	ArgsUsage: "max",
	Action:    actionCap,
	Description: `Cap will change the maximum number of VUs for a test.

   Because instantiating new VUs is a potentially very expensive operation,
   both in terms of CPU and RAM, you should be aware that you may see a bump in
   response times and skewed averages if you increase the cap during a running
   test.
   
   It's recommended to pause the test before creating a large number of VUs.
   
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
	Description: `Resume resumes a previously paused test.

   This is the opposite of the pause command, and will do nothing to an already
   running test.

   Endpoint: /v1/status`,
***REMOVED***

func dumpYAML(v interface***REMOVED******REMOVED***) error ***REMOVED***
	bytes, err := yaml.Marshal(v)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Serialization Error")
		return err
	***REMOVED***
	_, _ = os.Stdout.Write(bytes)
	return nil
***REMOVED***

func actionStatus(cc *cli.Context) error ***REMOVED***
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	status, err := client.Status()
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionScale(cc *cli.Context) error ***REMOVED***
	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***
	vus, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***

	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	status, err := client.UpdateStatus(lib.Status***REMOVED***VUs: null.IntFrom(vus)***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionCap(cc *cli.Context) error ***REMOVED***
	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***
	max, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***

	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	status, err := client.UpdateStatus(lib.Status***REMOVED***VUsMax: null.IntFrom(max)***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionPause(cc *cli.Context) error ***REMOVED***
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	status, err := client.UpdateStatus(lib.Status***REMOVED***Running: null.BoolFrom(false)***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***

func actionResume(cc *cli.Context) error ***REMOVED***
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	status, err := client.UpdateStatus(lib.Status***REMOVED***Running: null.BoolFrom(true)***REMOVED***)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Error")
		return err
	***REMOVED***
	return dumpYAML(status)
***REMOVED***
