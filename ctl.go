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
***REMOVED***

var commandScale = cli.Command***REMOVED***
	Name:      "scale",
	Usage:     "Scales a running test",
	ArgsUsage: "vus",
	Action:    actionScale,
***REMOVED***

var commandPause = cli.Command***REMOVED***
	Name:      "pause",
	Usage:     "Pauses a running test",
	ArgsUsage: " ",
	Action:    actionPause,
***REMOVED***

var commandResume = cli.Command***REMOVED***
	Name:      "resume",
	Usage:     "Resumes a paused test",
	ArgsUsage: " ",
	Action:    actionResume,
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
