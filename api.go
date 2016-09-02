package main

import (
	"context"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/lib"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"strconv"
	"time"
)

type APIServer struct ***REMOVED***
	Engine *lib.Engine
	Cancel context.CancelFunc

	Info lib.Info
***REMOVED***

// Run runs the API server.
// I'm not sure how idiomatic this is, probably not particularly...
func (s *APIServer) Run(ctx context.Context, addr string) ***REMOVED***
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(s.logRequestsMiddleware)
	router.Use(s.jsonErrorsMiddleware)

	router.Use(static.Serve("/", static.LocalFile("web/dist", true)))
	router.GET("/ping", func(c *gin.Context) ***REMOVED***
		c.Data(http.StatusNoContent, "", nil)
	***REMOVED***)
	v1 := router.Group("/v1")
	***REMOVED***
		v1.GET("/info", func(c *gin.Context) ***REMOVED***
			c.JSON(200, s.Info)
		***REMOVED***)
		v1.GET("/status", func(c *gin.Context) ***REMOVED***
			c.JSON(200, s.Engine.Status)
		***REMOVED***)
		v1.GET("/metrics", func(c *gin.Context) ***REMOVED***
			metrics := make(map[string]client.Metric)
			for m, samples := range s.Engine.Metrics ***REMOVED***
				metrics[m.Name] = client.Metric***REMOVED***
					Name:     m.Name,
					Type:     client.MetricType(m.Type),
					Contains: client.ValueType(m.Contains),
					Data:     m.Format(samples),
				***REMOVED***
			***REMOVED***
			c.JSON(200, metrics)
		***REMOVED***)
		v1.GET("/metrics/:name", func(c *gin.Context) ***REMOVED***
			name := c.Param("name")
			for m, samples := range s.Engine.Metrics ***REMOVED***
				if m.Name != name ***REMOVED***
					continue
				***REMOVED***

				c.JSON(200, client.Metric***REMOVED***
					Name:     m.Name,
					Type:     client.MetricType(m.Type),
					Contains: client.ValueType(m.Contains),
					Data:     m.Format(samples),
				***REMOVED***)
				return
			***REMOVED***
			c.AbortWithError(404, errors.New("No such metric"))
		***REMOVED***)
		v1.POST("/abort", func(c *gin.Context) ***REMOVED***
			s.Cancel()
			c.JSON(202, gin.H***REMOVED***"success": true***REMOVED***)
		***REMOVED***)
		v1.POST("/scale", func(c *gin.Context) ***REMOVED***
			vus, err := strconv.ParseInt(c.Query("vus"), 10, 64)
			if err != nil ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			if err := s.Engine.Scale(vus); err != nil ***REMOVED***
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			***REMOVED***

			c.JSON(202, gin.H***REMOVED***"success": true***REMOVED***)
		***REMOVED***)
	***REMOVED***
	router.NoRoute(func(c *gin.Context) ***REMOVED***
		c.JSON(404, gin.H***REMOVED***"error": "Not Found"***REMOVED***)
	***REMOVED***)

	srv := graceful.Server***REMOVED***NoSignalHandling: true, Server: &http.Server***REMOVED***Addr: addr, Handler: router***REMOVED******REMOVED***
	go srv.ListenAndServe()

	<-ctx.Done()
	srv.Stop(10 * time.Second)
	<-srv.StopChan()
***REMOVED***

func (s *APIServer) logRequestsMiddleware(c *gin.Context) ***REMOVED***
	path := c.Request.URL.Path
	c.Next()
	log.WithField("status", c.Writer.Status()).Debugf("%s %s", c.Request.Method, path)
***REMOVED***

func (s *APIServer) jsonErrorsMiddleware(c *gin.Context) ***REMOVED***
	c.Next()
	if c.Writer.Size() == 0 && len(c.Errors) > 0 ***REMOVED***
		c.JSON(c.Writer.Status(), c.Errors)
	***REMOVED***
***REMOVED***
