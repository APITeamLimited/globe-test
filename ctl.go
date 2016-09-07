package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/api"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
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

var commandAbort = cli.Command***REMOVED***
	Name:      "abort",
	Usage:     "Aborts a running test",
	ArgsUsage: " ",
	Action:    actionAbort,
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

	bytes, err := yaml.Marshal(status)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Serialization Error")
		return err
	***REMOVED***
	_, _ = os.Stdout.Write(bytes)

	return nil
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

	if err := client.Scale(vus); err != nil ***REMOVED***
		log.WithError(err).Error("Error")
	***REMOVED***
	return nil
***REMOVED***

func actionAbort(cc *cli.Context) error ***REMOVED***
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	if err := client.Abort(); err != nil ***REMOVED***
		log.WithError(err).Error("Error")
	***REMOVED***
	return nil
***REMOVED***
