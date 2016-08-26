package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat/client"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

func actionStatus(cc *cli.Context) error ***REMOVED***
	client, err := client.New(cc.GlobalString("address"))
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

	client, err := client.New(cc.GlobalString("address"))
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
	client, err := client.New(cc.GlobalString("address"))
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a client")
		return err
	***REMOVED***

	if err := client.Abort(); err != nil ***REMOVED***
		log.WithError(err).Error("Error")
	***REMOVED***
	return nil
***REMOVED***
