package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/loadimpact/speedboat/lib"
	"gopkg.in/urfave/cli.v1"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var commandRun = cli.Command***REMOVED***
	Name:      "run",
	Usage:     "Starts running a load test",
	ArgsUsage: "url|filename",
	Flags: []cli.Flag***REMOVED***
		cli.IntFlag***REMOVED***
			Name:  "vus, u",
			Usage: "virtual users to simulate",
			Value: 10,
		***REMOVED***,
		cli.DurationFlag***REMOVED***
			Name:  "duration, d",
			Usage: "test duration, 0 to run until cancelled",
			Value: 10 * time.Second,
		***REMOVED***,
		cli.StringFlag***REMOVED***
			Name:  "type, t",
			Usage: "input type, one of: auto, url, js",
			Value: "auto",
		***REMOVED***,
	***REMOVED***,
	Action: actionRun,
***REMOVED***

func guessType(filename string) string ***REMOVED***
	switch ***REMOVED***
	case strings.Contains(filename, "://"):
		return "url"
	case strings.HasSuffix(filename, ".js"):
		return "js"
	default:
		return ""
	***REMOVED***
***REMOVED***

func makeRunner(filename, t string) (lib.Runner, error) ***REMOVED***
	if t == "auto" ***REMOVED***
		t = guessType(filename)
	***REMOVED***
	return nil, nil
***REMOVED***

func actionRun(cc *cli.Context) error ***REMOVED***
	args := cc.Args()
	if len(args) != 1 ***REMOVED***
		return cli.NewExitError("Wrong number of arguments!", 1)
	***REMOVED***

	filename := args[0]
	runnerType := cc.String("type")
	runner, err := makeRunner(filename, runnerType)
	if err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't create a runner")
	***REMOVED***

	engine := lib.Engine***REMOVED***
		Runner: runner,
	***REMOVED***

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		defer func() ***REMOVED***
			log.Debug("Engine terminated")
			wg.Done()
		***REMOVED***()
		log.Debug("Starting engine...")
		if err := engine.Run(ctx); err != nil ***REMOVED***
			log.WithError(err).Error("Runtime Error")
		***REMOVED***
	***REMOVED***()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) ***REMOVED***
		path := c.Request.URL.Path
		c.Next()
		log.WithField("status", c.Writer.Status()).Debugf("%s %s", c.Request.Method, path)
	***REMOVED***)
	router.Use(func(c *gin.Context) ***REMOVED***
		c.Next()
		if c.Writer.Size() == 0 && len(c.Errors) > 0 ***REMOVED***
			c.JSON(c.Writer.Status(), c.Errors)
		***REMOVED***
	***REMOVED***)
	v1 := router.Group("/v1")
	***REMOVED***
		v1.GET("/info", func(c *gin.Context) ***REMOVED***
			c.JSON(200, gin.H***REMOVED***"version": cc.App.Version***REMOVED***)
		***REMOVED***)
		v1.GET("/status", func(c *gin.Context) ***REMOVED***
			c.JSON(200, engine.Status)
		***REMOVED***)
		v1.POST("/abort", func(c *gin.Context) ***REMOVED***
			cancel()
			wg.Wait()
			c.JSON(202, gin.H***REMOVED***"success": true***REMOVED***)
		***REMOVED***)
		v1.POST("/scale", func(c *gin.Context) ***REMOVED***
			vus, err := strconv.ParseInt(c.Query("vus"), 10, 64)
			if err != nil ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			if err := engine.Scale(vus); err != nil ***REMOVED***
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			***REMOVED***

			c.JSON(202, gin.H***REMOVED***"success": true***REMOVED***)
		***REMOVED***)
	***REMOVED***
	router.NoRoute(func(c *gin.Context) ***REMOVED***
		c.JSON(404, gin.H***REMOVED***"error": "Not Found"***REMOVED***)
	***REMOVED***)
	router.Run(cc.GlobalString("address"))

	return nil
***REMOVED***
