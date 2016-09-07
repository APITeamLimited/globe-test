package api

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/loadimpact/speedboat/lib"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"net/http"
	// "strconv"
	"time"
)

var contentType = "application/vnd.api+json"

type Server struct ***REMOVED***
	Engine *lib.Engine
	Cancel context.CancelFunc

	Info lib.Info
***REMOVED***

// Run runs the API server.
// I'm not sure how idiomatic this is, probably not particularly...
func (s *Server) Run(ctx context.Context, addr string) ***REMOVED***
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
			data, err := jsonapi.Marshal(s.Info)
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/error", func(c *gin.Context) ***REMOVED***
			c.AbortWithError(500, errors.New("This is an error"))
		***REMOVED***)
		v1.GET("/status", func(c *gin.Context) ***REMOVED***
			data, err := jsonapi.Marshal(s.Engine.Status)
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.PATCH("/status", func(c *gin.Context) ***REMOVED***
			// TODO: Allow full control of running/active/inactive VUs; stopping a test shouldn't
			// be final, and shouldn't implicitly affect anything else.
			if !s.Engine.Status.Running ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, errors.New("Test is stopped"))
				return
			***REMOVED***

			status := s.Engine.Status
			data, _ := ioutil.ReadAll(c.Request.Body)
			if err := jsonapi.Unmarshal(data, &status); err != nil ***REMOVED***
				c.AbortWithError(http.StatusBadRequest, err)
				return
			***REMOVED***

			if status.ActiveVUs != s.Engine.Status.ActiveVUs ***REMOVED***
				s.Engine.Scale(status.ActiveVUs)
			***REMOVED***
			if !s.Engine.Status.Running ***REMOVED***
				s.Cancel()
			***REMOVED***
			s.Engine.Status = status

			data, err := jsonapi.Marshal(s.Engine.Status)
			if err != nil ***REMOVED***
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/metrics", func(c *gin.Context) ***REMOVED***
			metrics := make([]interface***REMOVED******REMOVED***, 0, len(s.Engine.Metrics))
			for metric, sink := range s.Engine.Metrics ***REMOVED***
				metric.Sample = sink.Format()
				metrics = append(metrics, metric)
			***REMOVED***
			data, err := jsonapi.Marshal(metrics)
			if err != nil ***REMOVED***
				c.AbortWithError(500, err)
				return
			***REMOVED***
			c.Data(200, contentType, data)
		***REMOVED***)
		v1.GET("/metrics/:id", func(c *gin.Context) ***REMOVED***
			id := c.Param("id")
			for metric, sink := range s.Engine.Metrics ***REMOVED***
				if metric.Name != id ***REMOVED***
					continue
				***REMOVED***
				metric.Sample = sink.Format()
				data, err := jsonapi.Marshal(metric)
				if err != nil ***REMOVED***
					c.AbortWithError(500, err)
					return
				***REMOVED***
				c.Data(200, contentType, data)
				return
			***REMOVED***
			c.AbortWithError(404, errors.New("Metric not found"))
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

func (s *Server) logRequestsMiddleware(c *gin.Context) ***REMOVED***
	path := c.Request.URL.Path
	c.Next()
	log.WithField("status", c.Writer.Status()).Debugf("%s %s", c.Request.Method, path)
***REMOVED***

func (s *Server) jsonErrorsMiddleware(c *gin.Context) ***REMOVED***
	c.Header("Content-Type", contentType)
	c.Next()
	if len(c.Errors) > 0 ***REMOVED***
		var errors api2go.HTTPError
		for _, err := range c.Errors ***REMOVED***
			errors.Errors = append(errors.Errors, api2go.Error***REMOVED***
				Title: err.Error(),
			***REMOVED***)
		***REMOVED***
		data, _ := json.Marshal(errors)
		c.Data(c.Writer.Status(), contentType, data)
	***REMOVED***
***REMOVED***
